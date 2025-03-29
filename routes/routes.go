package routes

import (
	"compress/gzip"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/dustin/go-humanize"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pocka/legit/config"
	"github.com/pocka/legit/git"
	"github.com/pocka/legit/routes/preview"
	"github.com/russross/blackfriday/v2"
)

type deps struct {
	c *config.Config
}

func (d *deps) Index(w http.ResponseWriter, r *http.Request) {
	dirs, err := os.ReadDir(d.c.Repo.ScanPath)
	if err != nil {
		d.Write500(w)
		log.Printf("reading scan path: %s", err)
		return
	}

	summaries := []repositorySummary{}

	for _, dir := range dirs {
		name := dir.Name()
		if !dir.IsDir() || d.isIgnored(name) || d.isUnlisted(name) {
			continue
		}

		path, err := securejoin.SecureJoin(d.c.Repo.ScanPath, name)
		if err != nil {
			log.Printf("securejoin error: %v", err)
			d.Write404(w)
			return
		}

		gr, err := git.Open(path, "")
		if err != nil {
			log.Println(err)
			continue
		}

		c, err := gr.LastCommit()
		if err != nil {
			d.Write500(w)
			log.Println(err)
			return
		}

		summaries = append(summaries, repositorySummary{
			DisplayName:          getDisplayName(name),
			DirName:              name,
			Description:          getDescription(path),
			LastCommitAtRelative: humanize.Time(c.Committer.When),
			LastCommit:           c,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[j].LastCommit.Committer.When.Before(summaries[i].LastCommit.Committer.When)
	})

	tpath := filepath.Join(d.c.Dirs.Templates, "*")
	t := template.Must(template.ParseGlob(tpath))

	data := repoListData{
		Config:       d.c,
		Repositories: summaries,
	}

	if err := t.ExecuteTemplate(w, "repo-list", data); err != nil {
		log.Println(err)
		return
	}
}

func (d *deps) RepoIndex(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if d.isIgnored(name) {
		d.Write404(w)
		return
	}
	name = filepath.Clean(name)
	path, err := securejoin.SecureJoin(d.c.Repo.ScanPath, name)
	if err != nil {
		log.Printf("securejoin error: %v", err)
		d.Write404(w)
		return
	}

	gr, err := git.Open(path, "")
	if err != nil {
		d.Write404(w)
		return
	}

	commits, err := gr.Commits()
	if err != nil {
		d.Write500(w)
		log.Println(err)
		return
	}

	var readmeContent template.HTML
	for _, readme := range d.c.Repo.Readme {
		ext := filepath.Ext(readme)
		content, _ := gr.FileContent(readme)
		if len(content) > 0 {
			switch ext {
			case ".md", ".mkd", ".markdown":
				unsafe := blackfriday.Run(
					[]byte(content),
					blackfriday.WithExtensions(blackfriday.CommonExtensions),
				)
				html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
				readmeContent = template.HTML(html)
			default:
				safe := bluemonday.UGCPolicy().SanitizeBytes([]byte(content))
				readmeContent = template.HTML(
					fmt.Sprintf(`<pre>%s</pre>`, safe),
				)
			}
			break
		}
	}

	if readmeContent == "" {
		log.Printf("no readme found for %s", name)
	}

	mainBranch, err := gr.FindMainBranch(d.c.Repo.MainBranch)
	if err != nil {
		d.Write500(w)
		log.Println(err)
		return
	}

	tpath := filepath.Join(d.c.Dirs.Templates, "*")
	t := template.Must(template.ParseGlob(tpath))

	if len(commits) >= 3 {
		commits = commits[:3]
	}

	data := repoTopData{
		Config: d.c,
		Meta: repositoryMeta{
			DisplayName: getDisplayName(name),
			DirName:     name,
			Description: getDescription(path),
			Ref:         mainBranch,
		},
		Readme:        readmeContent,
		DefaultBranch: mainBranch,
		RecentCommits: commits,
		IsGoModule:    isGoModule(gr),
	}

	if err := t.ExecuteTemplate(w, "repo-top", data); err != nil {
		log.Println(err)
		return
	}

	return
}

func (d *deps) RepoTree(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if d.isIgnored(name) {
		d.Write404(w)
		return
	}
	treePath := r.PathValue("rest")
	ref := r.PathValue("ref")

	name = filepath.Clean(name)
	path, err := securejoin.SecureJoin(d.c.Repo.ScanPath, name)
	if err != nil {
		log.Printf("securejoin error: %v", err)
		d.Write404(w)
		return
	}
	gr, err := git.Open(path, ref)
	if err != nil {
		d.Write404(w)
		return
	}

	files, err := gr.FileTree(treePath)
	if err != nil {
		d.Write500(w)
		log.Println(err)
		return
	}

	relpath := []string{}
	if len(treePath) > 0 {
		relpath = strings.Split(treePath, "/")
	}

	data := repoTreeRefData{
		Config: d.c,
		Meta: repositoryMeta{
			DisplayName: getDisplayName(name),
			DirName:     name,
			Description: getDescription(path),
			Ref:         ref,
		},
		Path:  relpath,
		Files: files,
	}

	tpath := filepath.Join(d.c.Dirs.Templates, "*")
	t := template.Must(template.ParseGlob(tpath))

	if err := t.ExecuteTemplate(w, "repo-tree-ref", data); err != nil {
		log.Println(err)
		return
	}

	return
}

func (d *deps) FileContent(w http.ResponseWriter, r *http.Request) {
	var raw bool
	if rawParam, err := strconv.ParseBool(r.URL.Query().Get("raw")); err == nil {
		raw = rawParam
	}

	name := r.PathValue("name")
	if d.isIgnored(name) {
		d.Write404(w)
		return
	}
	treePath := r.PathValue("rest")
	ref := r.PathValue("ref")

	name = filepath.Clean(name)
	path, err := securejoin.SecureJoin(d.c.Repo.ScanPath, name)
	if err != nil {
		log.Printf("securejoin error: %v", err)
		d.Write404(w)
		return
	}

	gr, err := git.Open(path, ref)
	if err != nil {
		d.Write404(w)
		return
	}

	contents, err := gr.FileContent(treePath)
	if err != nil {
		d.Write500(w)
		return
	}

	if raw {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(contents))
		return
	}

	meta := repositoryMeta{
		DisplayName: getDisplayName(name),
		DirName:     name,
		Description: getDescription(path),
		Ref:         ref,
	}

	relpath := []string{}
	if len(treePath) > 0 {
		relpath = strings.Split(treePath, "/")
	}

	tpath := filepath.Join(d.c.Dirs.Templates, "*")
	t := template.Must(template.ParseGlob(tpath))

	if r.URL.Query().Has("preview") {
		previewType := r.URL.Query().Get("preview")

		for _, renderer := range preview.GetPreviewRenderers(treePath) {
			resolvedPreviewType := renderer.GetPreviewType()

			if previewType != "" && resolvedPreviewType != previewType {
				continue
			}

			switch resolvedPreviewType {
			case "html":
				html, err := renderer.Render([]byte(contents))
				if err != nil {
					log.Printf("Failed to render HTML preview: %s", err)
					d.Write500(w)
					return
				}

				data := repoBlobRefHTMLPreviewData{
					Config:  d.c,
					Meta:    meta,
					Path:    relpath,
					Content: template.HTML(html),
				}

				if err := t.ExecuteTemplate(w, "repo-blob-ref-html-preview", data); err != nil {
					log.Println(err)
					return
				}

				return
			}
		}

		log.Printf("Got ?preview=%s, but not preview renderer is available for the type", previewType)
		d.Write404(w)
		return
	}

	lc, err := countLines(strings.NewReader(contents))
	if err != nil {
		log.Printf("Failed to count lines for %s: %s", r.URL.Path, err)
		d.Write500(w)
		return
	}

	lines := make([]uint, lc)
	for i := range lines {
		if i < 0 {
			continue
		}

		lines[i] = uint(i + 1)
	}

	renderers := preview.GetPreviewRenderers(treePath)
	previewTypes := make([]string, len(renderers))
	for i, renderer := range renderers {
		previewTypes[i] = renderer.GetPreviewType()
	}

	data := repoBlobRefData{
		Config:       d.c,
		Meta:         meta,
		Path:         relpath,
		Content:      contents,
		LineNumbers:  lines,
		PreviewTypes: previewTypes,
	}

	if d.c.Meta.SyntaxHighlight {
		highlighted, err := highlightCode(treePath, contents)
		if err != nil {
			log.Println(err)
		} else {
			data.SyntaxHighlightedContent = highlighted
		}
	}

	if err := t.ExecuteTemplate(w, "repo-blob-ref", data); err != nil {
		log.Println(err)
		return
	}

	return
}

func (d *deps) Archive(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if d.isIgnored(name) {
		d.Write404(w)
		return
	}

	file := r.PathValue("file")

	// TODO: extend this to add more files compression (e.g.: xz)
	if !strings.HasSuffix(file, ".tar.gz") {
		d.Write404(w)
		return
	}

	ref := strings.TrimSuffix(file, ".tar.gz")

	// This allows the browser to use a proper name for the file when
	// downloading
	filename := fmt.Sprintf("%s-%s.tar.gz", name, ref)
	setContentDisposition(w, filename)
	setGZipMIME(w)

	path, err := securejoin.SecureJoin(d.c.Repo.ScanPath, name)
	if err != nil {
		log.Printf("securejoin error: %v", err)
		d.Write404(w)
		return
	}

	gr, err := git.Open(path, ref)
	if err != nil {
		d.Write404(w)
		return
	}

	gw := gzip.NewWriter(w)
	defer gw.Close()

	prefix := fmt.Sprintf("%s-%s", name, ref)
	err = gr.WriteTar(gw, prefix)
	if err != nil {
		// once we start writing to the body we can't report error anymore
		// so we are only left with printing the error.
		log.Println(err)
		return
	}

	err = gw.Flush()
	if err != nil {
		// once we start writing to the body we can't report error anymore
		// so we are only left with printing the error.
		log.Println(err)
		return
	}
}

func (d *deps) Log(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if d.isIgnored(name) {
		d.Write404(w)
		return
	}
	ref := r.PathValue("ref")

	path, err := securejoin.SecureJoin(d.c.Repo.ScanPath, name)
	if err != nil {
		log.Printf("securejoin error: %v", err)
		d.Write404(w)
		return
	}

	gr, err := git.Open(path, ref)
	if err != nil {
		d.Write404(w)
		return
	}

	commits, err := gr.Commits()
	if err != nil {
		d.Write500(w)
		log.Println(err)
		return
	}

	tpath := filepath.Join(d.c.Dirs.Templates, "*")
	t := template.Must(template.ParseGlob(tpath))

	data := repoLogRefData{
		Config: d.c,
		Meta: repositoryMeta{
			DisplayName: getDisplayName(name),
			DirName:     name,
			Description: getDescription(path),
			Ref:         ref,
		},
		Commits: commits,
	}

	if err := t.ExecuteTemplate(w, "repo-log-ref", data); err != nil {
		log.Println(err)
		return
	}
}

func (d *deps) Diff(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if d.isIgnored(name) {
		d.Write404(w)
		return
	}
	ref := r.PathValue("ref")

	path, err := securejoin.SecureJoin(d.c.Repo.ScanPath, name)
	if err != nil {
		log.Printf("securejoin error: %v", err)
		d.Write404(w)
		return
	}
	gr, err := git.Open(path, ref)
	if err != nil {
		d.Write404(w)
		return
	}

	diff, err := gr.Diff()
	if err != nil {
		d.Write500(w)
		log.Println(err)
		return
	}

	tpath := filepath.Join(d.c.Dirs.Templates, "*")
	t := template.Must(template.ParseGlob(tpath))

	data := repoCommitData{
		Config: d.c,
		Meta: repositoryMeta{
			DisplayName: getDisplayName(name),
			DirName:     name,
			Description: getDescription(path),
			Ref:         diff.Commit.Hash.String(),
		},
		Commit: diff.Commit,
		Parent: diff.Parent,
		Diff:   diff,
	}

	if err := t.ExecuteTemplate(w, "repo-commit", data); err != nil {
		log.Println(err)
		return
	}
}

func (d *deps) Refs(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if d.isIgnored(name) {
		d.Write404(w)
		return
	}

	path, err := securejoin.SecureJoin(d.c.Repo.ScanPath, name)
	if err != nil {
		log.Printf("securejoin error: %v", err)
		d.Write404(w)
		return
	}

	gr, err := git.Open(path, "")
	if err != nil {
		d.Write404(w)
		return
	}

	tags, err := gr.Tags()
	if err != nil {
		// Non-fatal, we *should* have at least one branch to show.
		log.Println(err)
	}

	branches, err := gr.Branches()
	if err != nil {
		log.Println(err)
		d.Write500(w)
		return
	}

	mainBranch, err := gr.FindMainBranch(d.c.Repo.MainBranch)
	if err != nil {
		d.Write500(w)
		log.Println(err)
		return
	}

	tpath := filepath.Join(d.c.Dirs.Templates, "*")
	t := template.Must(template.ParseGlob(tpath))

	data := repoRefsData{
		Config: d.c,
		Meta: repositoryMeta{
			DisplayName: getDisplayName(name),
			DirName:     name,
			Description: getDescription(path),
			Ref:         mainBranch,
		},
		Tags:     tags,
		Branches: branches,
	}

	if err := t.ExecuteTemplate(w, "repo-refs", data); err != nil {
		log.Println(err)
		return
	}
}

func (d *deps) ServeStatic(w http.ResponseWriter, r *http.Request) {
	f := r.PathValue("file")
	f = filepath.Clean(f)
	f, err := securejoin.SecureJoin(d.c.Dirs.Static, f)
	if err != nil {
		d.Write404(w)
		return
	}

	http.ServeFile(w, r, f)
}

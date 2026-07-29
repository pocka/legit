package routes

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/pocka/legit/git"
)

func (d *deps) serveRepoIndex(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if d.isIgnored(name) {
		d.write404(w)
		return
	}
	name = filepath.Clean(name)
	path, err := securejoin.SecureJoin(d.c.Repo.ScanPath, name)
	if err != nil {
		log.Printf("securejoin error: %v", err)
		d.write404(w)
		return
	}

	gr, err := git.Open(path, "")
	if err != nil {
		d.write404(w)
		return
	}

	commits, _, err := gr.Commits(git.CommitsOptions{Limit: 3})
	if err != nil {
		d.write500(w)
		log.Println(err)
		return
	}

	mainBranch, err := gr.FindMainBranch(d.c.Repo.MainBranch)
	if err != nil {
		d.write500(w)
		log.Println(err)
		return
	}

	transformer := newRepoLinkTransformer(name, mainBranch)

	var readmeContent template.HTML
	for _, readme := range d.c.Repo.Readme {
		file, _ := gr.File(readme)
		if file == nil {
			continue
		}

		content, _ := file.Contents()

		// Skip empty files.
		if len(content) > 0 {
			renderer := d.htmlRenderer(file)
			if renderer == nil {
				renderer = &d.plaintext
			}

			result, err := renderer.Render([]byte(content), transformer)
			if err != nil {
				log.Printf("Unable to render %s/%s, skipping.", name, readme)
				continue
			}

			readmeContent = template.HTML(result)
			break
		}
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

	if err := d.template().ExecuteTemplate(w, "repo-top", data); err != nil {
		log.Println(err)
		return
	}
}

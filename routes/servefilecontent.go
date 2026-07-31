package routes

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

func (d *deps) serveFileContent(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Has("raw")

	ref := r.PathValue("ref")
	gr, name, err := d.openRepository(r.PathValue("name"), ref)
	if err != nil {
		d.write404(w)
		return
	}
	treePath := r.PathValue("rest")

	file, err := gr.File(treePath)
	if err != nil {
		d.write500(w)
		return
	}

	contents, err := file.Contents()
	if err != nil {
		d.write500(w)
		return
	}

	if raw {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(contents))
		return
	}

	if isBin, _ := file.IsBinary(); isBin {
		contents = "Not displaying binary file"
	}

	meta := repositoryMeta{
		DisplayName: getDisplayName(name),
		DirName:     name,
		Description: gr.GitwebDescription(),
		Ref:         ref,
	}

	relpath := []string{}
	if len(treePath) > 0 {
		relpath = strings.Split(treePath, "/")
	}

	if r.URL.Query().Has("preview") {
		previewType := r.URL.Query().Get("preview")

		switch previewType {
		case "html":
			renderer := d.htmlRenderer(file)
			if renderer == nil {
				log.Printf("Requested HTML preview for %s/%s, but the filetype has no HTML renderer", name, treePath)
				d.write404(w)
				return
			}

			html, err := renderer.Render([]byte(contents), newRepoLinkTransformer(name, ref))
			if err != nil {
				log.Printf("Failed to render HTML preview: %s", err)
				d.write500(w)
				return
			}

			data := repoBlobRefHTMLPreviewData{
				Config:  d.c,
				Meta:    meta,
				Path:    relpath,
				Content: template.HTML(html),
			}

			if err := d.template().ExecuteTemplate(w, "repo-blob-ref-html-preview", data); err != nil {
				log.Println(err)
				return
			}

			return
		default:
			log.Printf("Got ?preview=%s, but not preview renderer is available for the type", previewType)
			d.write404(w)
			return
		}
	}

	lc, err := countLines(strings.NewReader(contents))
	if err != nil {
		log.Printf("Failed to count lines for %s: %s", r.URL.Path, err)
		d.write500(w)
		return
	}

	lines := make([]uint, lc)
	for i := range lines {
		if i < 0 {
			continue
		}

		lines[i] = uint(i + 1)
	}

	previewTypes := make([]string, 0, 1)
	if d.htmlRenderer(file) != nil {
		previewTypes = append(previewTypes, "html")
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

	if err := d.template().ExecuteTemplate(w, "repo-blob-ref", data); err != nil {
		log.Println(err)
		return
	}
}

package routes

import (
	"log"
	"net/http"
	"os"
	"sort"
)

func (d *deps) serveIndex(w http.ResponseWriter, r *http.Request) {
	dirs, err := os.ReadDir(d.c.Repo.ScanPath)
	if err != nil {
		d.write500(w)
		log.Printf("reading scan path: %s", err)
		return
	}

	summaries := []repositorySummary{}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		gr, name, err := d.openRepository(dir.Name(), "")
		if err != nil {
			d.write404(w)
			return
		}

		if d.isIgnored(name) || d.isUnlisted(name) {
			continue
		}

		c, err := gr.LastCommit()
		if err != nil {
			d.write500(w)
			log.Println(err)
			return
		}

		var category string
		if d.c.UI.Category.Grouping {
			category = gr.GitwebCategory()
		}

		summaries = append(summaries, repositorySummary{
			DisplayName: getDisplayName(name),
			DirName:     name,
			Description: gr.GitwebDescription(),
			Category:    category,
			LastCommit:  c,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[j].LastCommit.Committer.When.Before(summaries[i].LastCommit.Committer.When)
	})

	data := repoListData{
		Config:       d.c,
		Repositories: summaries,
	}

	if err := d.template().ExecuteTemplate(w, "repo-list", data); err != nil {
		log.Println(err)
		return
	}
}

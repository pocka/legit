package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pocka/legit/git"
)

func (d *deps) serveLog(w http.ResponseWriter, r *http.Request) {
	ref := r.PathValue("ref")
	gr, name, err := d.openRepository(r.PathValue("name"), ref)
	if err != nil {
		d.write404(w)
		return
	}

	limit := d.c.UI.CommitsPageSize

	opts := git.CommitsOptions{Limit: limit}
	query := r.URL.Query()

	if after := query.Get("before"); plumbing.IsHash(after) {
		opts.Before = plumbing.NewHash(after)
	}
	if before := query.Get("after"); plumbing.IsHash(before) {
		opts.After = plumbing.NewHash(before)
	}

	commits, page, err := gr.Commits(opts)
	if err != nil {
		d.write500(w)
		log.Println(err)
		return
	}

	prevPageHref := ""
	nextPageHref := ""

	if len(commits) > 0 {
		if page.HasPrevPage {
			prevPageHref = fmt.Sprintf("/%s/log/%s?after=%s", name, ref, commits[0].Hash.String())
		}

		if page.HasNextPage {
			nextPageHref = fmt.Sprintf("/%s/log/%s?before=%s", name, ref, commits[len(commits)-1].Hash.String())
		}
	}

	data := repoLogRefData{
		Config: d.c,
		Meta: repositoryMeta{
			DisplayName: getDisplayName(name),
			DirName:     name,
			Description: gr.GitwebDescription(),
			Ref:         ref,
		},
		Commits:      commits,
		PrevPageHref: prevPageHref,
		NextPageHref: nextPageHref,
	}

	if err := d.template().ExecuteTemplate(w, "repo-log-ref", data); err != nil {
		log.Println(err)
		return
	}
}

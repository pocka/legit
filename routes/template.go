package routes

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func (d *deps) Write404(w http.ResponseWriter) {
	tpath := filepath.Join(d.c.Dirs.Templates, "*")
	t := template.Must(template.ParseGlob(tpath))

	data := error404Data{
		Config: d.c,
	}

	w.WriteHeader(404)
	if err := t.ExecuteTemplate(w, "404", data); err != nil {
		log.Printf("404 template: %s", err)
	}
}

func (d *deps) Write500(w http.ResponseWriter) {
	tpath := filepath.Join(d.c.Dirs.Templates, "*")
	t := template.Must(template.ParseGlob(tpath))

	data := error500Data{
		Config: d.c,
	}

	w.WriteHeader(500)
	if err := t.ExecuteTemplate(w, "500", data); err != nil {
		log.Printf("500 template: %s", err)
	}
}

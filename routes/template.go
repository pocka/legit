package routes

import (
	"html/template"
	"log"
	"net/http"
)

func (d *deps) Template() *template.Template {
	if d.t != nil {
		return d.t
	}

	t := template.Must(template.ParseFS(d.templatesDir, "*"))
	return t
}

func (d *deps) Write404(w http.ResponseWriter) {
	data := error404Data{
		Config: d.c,
	}

	w.WriteHeader(404)
	if err := d.Template().ExecuteTemplate(w, "404", data); err != nil {
		log.Printf("404 template: %s", err)
	}
}

func (d *deps) Write500(w http.ResponseWriter) {
	data := error500Data{
		Config: d.c,
	}

	w.WriteHeader(500)
	if err := d.Template().ExecuteTemplate(w, "500", data); err != nil {
		log.Printf("500 template: %s", err)
	}
}

package routes

import (
	"html/template"
	"log"
	"net/http"
)

func (d *deps) Write404(w http.ResponseWriter) {
	t := template.Must(template.ParseFS(d.templatesDir, "*"))

	data := error404Data{
		Config: d.c,
	}

	w.WriteHeader(404)
	if err := t.ExecuteTemplate(w, "404", data); err != nil {
		log.Printf("404 template: %s", err)
	}
}

func (d *deps) Write500(w http.ResponseWriter) {
	t := template.Must(template.ParseFS(d.templatesDir, "*"))

	data := error500Data{
		Config: d.c,
	}

	w.WriteHeader(500)
	if err := t.ExecuteTemplate(w, "500", data); err != nil {
		log.Printf("500 template: %s", err)
	}
}

package routes

import (
	"html/template"
	"io/fs"
	"net/http"
	"os"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pocka/legit/config"
	"github.com/pocka/legit/renderer/html"
	"github.com/pocka/legit/routes/debug"
)

// Checks for gitprotocol-http(5) specific smells; if found, passes
// the request on to the git http service, else render the web frontend.
func (d *deps) multiplex(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if d.isIgnored(name) {
		d.write404(w)
		return
	}

	path := r.PathValue("rest")

	if r.URL.RawQuery == "service=git-receive-pack" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no pushing allowed!"))
		return
	}

	if path == "info/refs" &&
		r.URL.RawQuery == "service=git-upload-pack" &&
		r.Method == "GET" {
		d.serveInfoRefs(w, r)
	} else if path == "git-upload-pack" && r.Method == "POST" {
		d.serveUploadPack(w, r)
	} else if r.Method == "GET" {
		d.serveRepoIndex(w, r)
	}
}

func Handler(c *config.Config, scanRoot *os.Root, staticDir fs.FS, templatesDir fs.FS) *http.ServeMux {
	mux := http.NewServeMux()
	ugcPolicy := bluemonday.UGCPolicy()

	d := deps{
		c:            c,
		scanRoot:     scanRoot,
		staticDir:    staticDir,
		templatesDir: templatesDir,
		markdown:     html.NewMarkdownRenderer(ugcPolicy),
		plaintext:    html.NewPlaintextRenderer(ugcPolicy),
	}

	if !c.CompileTemplatesOnRequest {
		d.t = template.Must(template.ParseFS(d.templatesDir, "*"))
	}

	debug.Register(mux)

	mux.HandleFunc("GET /", d.serveIndex)
	mux.HandleFunc("GET /static/{file}", d.serveStatic)
	mux.HandleFunc("GET /{name}", d.multiplex)
	mux.HandleFunc("POST /{name}", d.multiplex)
	mux.HandleFunc("GET /{name}/tree/{ref}/{rest...}", d.serveRepoTree)
	mux.HandleFunc("GET /{name}/blob/{ref}/{rest...}", d.serveFileContent)
	mux.HandleFunc("GET /{name}/log/{ref}", d.serveLog)
	mux.HandleFunc("GET /{name}/archive/{file}", d.serveArchive)
	mux.HandleFunc("GET /{name}/commit/{ref}", d.serveDiff)
	mux.HandleFunc("GET /{name}/refs/{$}", d.serveRefs)
	mux.HandleFunc("GET /{name}/{rest...}", d.multiplex)
	mux.HandleFunc("POST /{name}/{rest...}", d.multiplex)

	return mux
}

package routes

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"

	"github.com/pocka/legit/git/service"
)

func (d *deps) serveInfoRefs(w http.ResponseWriter, r *http.Request) {
	repo, err := d.resolveRepository(r.PathValue("name"))
	if err != nil {
		d.write404(w)
		return
	}

	w.Header().Set("content-type", "application/x-git-upload-pack-advertisement")
	w.WriteHeader(http.StatusOK)

	cmd := service.ServiceCommand{
		Dir:    repo,
		Stdout: w,
	}

	if err := cmd.InfoRefs(); err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf("git: failed to execute git-upload-pack (info/refs) %s", err)
		return
	}
}

func (d *deps) serveUploadPack(w http.ResponseWriter, r *http.Request) {
	repo, err := d.resolveRepository(r.PathValue("name"))
	if err != nil {
		d.write404(w)
		return
	}

	w.Header().Set("content-type", "application/x-git-upload-pack-result")
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	cmd := service.ServiceCommand{
		Dir:    repo,
		Stdout: w,
	}

	var reader io.ReadCloser
	reader = r.Body

	if r.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Printf("git: failed to create gzip reader: %s", err)
			return
		}
		defer reader.Close()
	}

	cmd.Stdin = reader
	if err := cmd.UploadPack(); err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf("git: failed to execute git-upload-pack %s", err)
		return
	}
}

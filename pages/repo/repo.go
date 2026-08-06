// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	e "errors"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/pocka/legit/core"
	"github.com/pocka/legit/pages/errors"
)

const (
	treePrefix    = "/tree/"
	blobPrefix    = "/blob/"
	logPrefix     = "/log/"
	commitPrefix  = "/commit/"
	archivePrefix = "/archive/"
)

var (
	errRefNotFound = e.New("ref not found")
)

type Repo struct {
	core *core.Core
	r    *git.Repository

	// innerPath is r.URL.Path without repository name.
	innerPath string

	// path represents repository path.
	path string

	// dirname is a result of "filepath.Base(path)", for convenience.
	dirname string
}

func New(repositoryPath string, innerPath string, core *core.Core) (*Repo, error) {
	var err error
	repo := Repo{core: core, innerPath: innerPath, path: repositoryPath, dirname: filepath.Base(repositoryPath)}

	repo.r, err = git.PlainOpen(repositoryPath)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", repositoryPath, err)
	}

	return &repo, nil
}

func (repo *Repo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	service := r.URL.Query().Get("service")

	if service == "git-receive-pack" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This is readonly repository"))
		return
	}

	if repo.innerPath == "/info/refs" && service == "git-upload-pack" && r.Method == http.MethodGet {
		repo.inforefs(w, r)
		return
	}

	if repo.innerPath == "/git-upload-pack" && r.Method == http.MethodPost {
		repo.uploadpack(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	if repo.innerPath == "/" {
		repo.index(w, r)
		return
	}

	// "/refs"
	if repo.innerPath == "/refs" {
		repo.refs(w, r)
		return
	}

	// "/tree/<ref>*/<path>*"
	if strings.Index(repo.innerPath, treePrefix) == 0 {
		repo.tree(strings.TrimPrefix(repo.innerPath, treePrefix), w, r)
		return
	}

	// "/blob/<ref>*/<path>*"
	if strings.Index(repo.innerPath, blobPrefix) == 0 {
		repo.blob(strings.TrimPrefix(repo.innerPath, blobPrefix), w, r)
		return
	}

	// "/log/<ref>*"
	if strings.Index(repo.innerPath, logPrefix) == 0 {
		repo.log(strings.TrimPrefix(repo.innerPath, logPrefix), w, r)
		return
	}

	// "/commit/<hash>"
	if strings.Index(repo.innerPath, commitPrefix) == 0 {
		repo.commit(strings.TrimPrefix(repo.innerPath, commitPrefix), w, r)
		return
	}

	// "/archive/<ref>*"
	if strings.Index(repo.innerPath, archivePrefix) == 0 {
		repo.archive(strings.TrimPrefix(repo.innerPath, archivePrefix), w, r)
		return
	}

	errors.WriteNotFound(repo.core, w, r)
}

func takeRef(repo *git.Repository, relpath string) (resolved *plumbing.Hash, ref string, childPath string, err error) {
	segments := strings.Split(relpath, "/")

	for i := range segments {
		ref = path.Join(segments[:i+1]...)
		resolved, err = repo.ResolveRevision(plumbing.Revision(ref))
		if err != nil {
			continue
		}

		if i < len(segments)-1 {
			childPath = path.Join(segments[i+1:]...)
		}

		return
	}

	return nil, "", "", errRefNotFound
}

// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

// save-commit-pages is a helper program for detecting markup changes
// for commit pages. This program saves all commit pages of a repository
// into a single directory. The intended usecase is to run this program
// against a repository A on legit at revision X and the same repository
// A on legit at revision Y, then compare the two output directories using
// diff command line tool.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	baseUrl := flag.String("base-url", "http://localhost:5555/", "legit server URL, excluding project directory")
	outdir := flag.String("outdir", "", "Output directory to store downloaded HTML files")
	workers := flag.Uint("workers", 1, "Number of workers")
	flag.Parse()

	if *outdir == "" {
		log.Fatal("--outdir is required")
	}

	if *workers == 0 {
		log.Fatal("--workers must be greater than 0")
	}

	repoPath := flag.Arg(0)
	if repoPath == "" {
		log.Fatal("Argument is required for repository path")
	}

	repoName := filepath.Base(repoPath)
	commitUrlPrefix, err := url.JoinPath(*baseUrl, repoName, "commit")
	if err != nil {
		log.Fatal(err)
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir(*outdir, 0o755); err != nil {
		log.Fatal(err)
	}

	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	commits := make(chan plumbing.Hash)

	var wg sync.WaitGroup
	for range *workers {
		wg.Go(func() {
			startWorker(ctx, commitUrlPrefix, *outdir, commits)
		})
	}

	for {
		if next, err := iter.Next(); err != nil {
			if err == io.EOF {
				break
			}

			log.Fatal(err)
		} else {
			commits <- next.Hash
		}
	}

	close(commits)
	wg.Wait()
}

func startWorker(ctx context.Context, baseUrl string, outdir string, commits <-chan plumbing.Hash) {
	for {
		select {
		case hash, ok := <-commits:
			if ok {
				if err := fetchAndSave(ctx, hash, baseUrl, outdir); err != nil {
					log.Printf("Download failed (%s): %s", hash, err)
				}
			} else {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func fetchAndSave(_ context.Context, hash plumbing.Hash, baseUrl string, outdir string) error {
	target, err := url.JoinPath(baseUrl, hash.String())
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Get(target)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected HTTP %d, got HTTP %d", http.StatusOK, resp.StatusCode)
	}

	file, err := os.Create(filepath.Join(outdir, fmt.Sprintf("%s.html", hash)))
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}

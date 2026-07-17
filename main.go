package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/pocka/legit/config"
	"github.com/pocka/legit/embed"
	"github.com/pocka/legit/routes"
)

func main() {
	var cfg string
	var host string
	var port uint
	flag.StringVar(&cfg, "config", "./config.yaml", "path to config file")
	flag.StringVar(&host, "server.host", "", "override server.host config")
	flag.UintVar(&port, "server.port", 0, "override server.port config")
	flag.Parse()

	c, err := config.Read(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if port > 0 {
		if port > math.MaxUint16 {
			log.Fatalf("server.port should be in 0 < x <= %d range", math.MaxUint16)
		}

		c.Server.Port = int(port)
	}

	if host != "" {
		c.Server.Host = host
	}

	if c.UI.CommitsPageSize == 0 {
		c.UI.CommitsPageSize = 30
	}

	allowedDirs := make([]string, 1, 3)
	allowedDirs[0] = c.Repo.ScanPath

	if c.Dirs.Static != "" {
		allowedDirs = append(allowedDirs, c.Dirs.Static)
	}

	if c.Dirs.Templates != "" {
		allowedDirs = append(allowedDirs, c.Dirs.Templates)
	}

	if err := restrictFileAccessTo(allowedDirs...); err != nil {
		log.Fatalf("Unable to restrict filesystem access: %s", err)
	}

	var staticDir fs.FS
	if c.Dirs.Static != "" {
		root, err := os.OpenRoot(c.Dirs.Static)
		if err != nil {
			log.Fatalf("Unable to open static dir: %s", err)
		}
		defer root.Close()

		staticDir = root.FS()
	} else {
		staticDir = embed.StaticDir()
	}

	var templatesDir fs.FS
	if c.Dirs.Templates != "" {
		templatesDir = os.DirFS(c.Dirs.Templates)
	} else {
		templatesDir = embed.TemplatesDir()
	}

	mux := routes.Handlers(c, staticDir, templatesDir)
	addr := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
	log.Println("starting server on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

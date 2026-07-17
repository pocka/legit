package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
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
	var scanPath string
	flag.StringVar(&cfg, "config", "", "path to config file")
	flag.StringVar(&host, "server.host", "", "override server.host config")
	flag.UintVar(&port, "server.port", 0, "override server.port config")
	flag.StringVar(&scanPath, "repo.scanPath", "", "override repo.scanPath config")
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	c := config.NewWithDefaults()

	if cfg != "" {
		if err := c.ReadFromFile(cfg); err != nil {
			log.Fatalf("Unable to read config file: %s", err)
		}
	}

	if port > 0 {
		c.Server.Port = port
	}

	if host != "" {
		c.Server.Host = host
	}

	if scanPath != "" {
		c.Repo.ScanPath = scanPath
	}

	if err := c.Resolve(cwd); err != nil {
		log.Fatal(err)
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

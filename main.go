package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/pocka/legit/config"
	"github.com/pocka/legit/routes"
)

//go:embed static/*
var defaultStaticDir embed.FS

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

	if err := UnveilPaths([]string{
		c.Dirs.Static,
		c.Repo.ScanPath,
		c.Dirs.Templates,
	},
		"r"); err != nil {
		log.Fatalf("unveil: %s", err)
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
		staticDir, err = fs.Sub(defaultStaticDir, "static")
		if err != nil {
			log.Fatalf("Unable to open default static dir: %s", err)
		}
	}

	mux := routes.Handlers(c, staticDir)
	addr := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
	log.Println("starting server on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

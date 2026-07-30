package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pocka/legit/config"
	"github.com/pocka/legit/embed"
	"github.com/pocka/legit/git/exe"
	"github.com/pocka/legit/routes"
)

var additionalAccessDirs string

func main() {
	var cfg string
	var host string
	var port uint
	var scanPath string
	var compileTemplatesOnRequest bool
	var staticDirOverride string
	var templatesDirOverride string
	flag.StringVar(&cfg, "config", "", "path to config file")
	flag.StringVar(&host, "server.host", "", "override server.host config")
	flag.UintVar(&port, "server.port", 0, "override server.port config")
	flag.StringVar(&scanPath, "repo.scanPath", "", "override repo.scanPath config")
	flag.BoolVar(&compileTemplatesOnRequest, "compileTemplatesOnRequest", false, "override compileTemplatesOnRequest config")
	flag.StringVar(&staticDirOverride, "dirs.static", "", "override dirs.static config")
	flag.StringVar(&templatesDirOverride, "dirs.templates", "", "override dirs.templates config")
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

	if compileTemplatesOnRequest {
		c.CompileTemplatesOnRequest = true
	}

	if staticDirOverride != "" {
		resolved, err := filepath.Abs(staticDirOverride)
		if err != nil {
			log.Fatalf("Cannot resolve -dirs.static")
		}

		c.Dirs.Static = resolved
	}

	if templatesDirOverride != "" {
		resolved, err := filepath.Abs(templatesDirOverride)
		if err != nil {
			log.Fatalf("Cannot resolve -dirs.templates")
		}

		c.Dirs.Templates = resolved
	}

	if err := c.Resolve(cwd); err != nil {
		log.Fatal(err)
	}

	fsAllowList := make([]filesystemAccess, 2, 10)
	fsAllowList[0] = filesystemAccess{
		path:  c.Repo.ScanPath,
		isDir: true,
		read:  true,
	}

	// os.exec.Cmd use /dev/null. Without this, git operations fail with
	// "open /dev/null: permission denied".
	// https://rohitpaulk.com/articles/cmd-run-dev-null
	fsAllowList[1] = filesystemAccess{
		path:  "/dev/null",
		read:  true,
		write: true,
	}

	if path, err := exec.LookPath(exe.GitBin()); err != nil {
		log.Printf("Unable to find git binary, git operations will fail: %s", err)
	} else {
		fsAllowList = append(fsAllowList, filesystemAccess{
			path:    path,
			read:    true,
			execute: true,
		})
	}

	if c.Dirs.Static != "" {
		fsAllowList = append(fsAllowList, filesystemAccess{
			path:  c.Dirs.Static,
			isDir: true,
			read:  true,
		})
	}

	if c.Dirs.Templates != "" {
		fsAllowList = append(fsAllowList, filesystemAccess{
			path:  c.Dirs.Templates,
			isDir: true,
			read:  true,
		})
	}

	if additionalAccessDirs != "" {
		for path := range strings.SplitSeq(additionalAccessDirs, ",") {
			path := strings.Trim(path, " ")

			fsAllowList = append(fsAllowList, filesystemAccess{
				path:  path,
				isDir: true,
				read:  true,
			})
		}
	}

	if err := restrictFileAccessTo(fsAllowList...); err != nil {
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

	mux := routes.Handler(c, staticDir, templatesDir)
	addr := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
	log.Println("starting server on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

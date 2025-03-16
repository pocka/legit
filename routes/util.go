package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"git.icyphox.sh/legit/git"
)

func isGoModule(gr *git.GitRepo) bool {
	_, err := gr.FileContent("go.mod")
	return err == nil
}

func getDisplayName(name string) string {
	return strings.TrimSuffix(name, ".git")
}

func getDescription(path string) (desc string) {
	db, err := os.ReadFile(filepath.Join(path, "description"))
	if err == nil {
		desc = string(db)
	} else {
		desc = ""
	}
	return
}

func (d *deps) isUnlisted(name string) bool {
	return slices.Contains(d.c.Repo.Unlisted, name)
}

func (d *deps) isIgnored(name string) bool {
	return slices.Contains(d.c.Repo.Ignore, name)
}

type repoInfo struct {
	Git      *git.GitRepo
	Path     string
	Category string
}

func setContentDisposition(w http.ResponseWriter, name string) {
	h := "inline; filename=\"" + name + "\""
	w.Header().Add("Content-Disposition", h)
}

func setGZipMIME(w http.ResponseWriter) {
	setMIME(w, "application/gzip")
}

func setMIME(w http.ResponseWriter, mime string) {
	w.Header().Add("Content-Type", mime)
}

func countLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	bufLen := 0
	count := 0
	nl := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		if c > 0 {
			bufLen += c
		}
		count += bytes.Count(buf[:c], nl)

		switch {
		case err == io.EOF:
			/* handle last line not having a newline at the end */
			if bufLen >= 1 && buf[(bufLen-1)%(32*1024)] != '\n' {
				count++
			}
			return count, nil
		case err != nil:
			return 0, err
		}
	}
}

func highlightCode(fileName string, code string, styleQuery string) (template.HTML, error) {
	lexer := lexers.Get(fileName)

	// Do not process if no appropriate highlighter was found.
	if lexer == nil {
		return "", nil
	}

	style := styles.Get(styleQuery)
	if style == nil {
		return "", fmt.Errorf("No chroma style found for '%s'", styleQuery)
	}

	formatter := html.New()

	iter, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", fmt.Errorf("Failed to tokenize code: %s", err)
	}

	var output bytes.Buffer
	err = formatter.Format(&output, style, iter)
	if err != nil {
		return "", fmt.Errorf("Failed to highlight code: %s", err)
	}

	return template.HTML(output.String()), nil
}

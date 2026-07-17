package config

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Repo struct {
		ScanPath   string   `yaml:"scanPath"`
		Readme     []string `yaml:"readme"`
		MainBranch []string `yaml:"mainBranch"`
		Ignore     []string `yaml:"ignore,omitempty"`
		Unlisted   []string `yaml:"unlisted,omitempty"`
	} `yaml:"repo"`
	Dirs struct {
		Templates string `yaml:"templates"`
		Static    string `yaml:"static"`
	} `yaml:"dirs"`
	Meta struct {
		Title           string `yaml:"title"`
		Description     string `yaml:"description"`
		SyntaxHighlight bool   `yaml:"syntaxHighlight"`
	} `yaml:"meta"`
	UI struct {
		CommitsPageSize uint32 `yaml:"commitsPageSize"`
		Footer          struct {
			Links []struct {
				Text string `yaml:"text"`
				Href string `yaml:"href"`
			} `yaml:"links"`
			PoweredBy bool `yaml:"poweredBy"`
		} `yaml:"footer"`
	} `yaml:"ui"`
	CompileTemplatesOnRequest bool `yaml:"compileTemplatesOnRequest"`
	Server                    struct {
		Name string `yaml:"name,omitempty"`
		Host string `yaml:"host"`
		Port uint   `yaml:"port"`
	} `yaml:"server"`

	filepath string
}

func NewWithDefaults() *Config {
	c := Config{}

	c.Repo.MainBranch = []string{"trunk", "master", "main"}
	c.Repo.Readme = []string{
		"README", "README.txt", "README.md", "README.adoc",
		"readme", "readme.txt", "readme.md", "readme.adoc",
	}

	c.Meta.Title = "Repositories"

	return &c
}

func (c *Config) ReadFromFile(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Read error (%s): %w", path, err)
	}

	if err := yaml.Unmarshal(b, &c); err != nil {
		return fmt.Errorf("Parsing error: %w", err)
	}

	c.filepath = path

	return nil
}

func (c *Config) Resolve(cwd string) error {
	var err error

	basePath := cwd
	if c.filepath != "" {
		basePath = filepath.Dir(c.filepath)
	}

	if c.Repo.ScanPath == "" {
		c.Repo.ScanPath = cwd
	} else if c.Repo.ScanPath, err = resolvePath(c.Repo.ScanPath, basePath); err != nil {
		return err
	}

	// Override templates dir
	if c.Dirs.Templates != "" {
		if c.Dirs.Templates, err = resolvePath(c.Dirs.Templates, basePath); err != nil {
			return err
		}
	}

	// Override static dir
	if c.Dirs.Static != "" {
		if c.Dirs.Static, err = resolvePath(c.Dirs.Static, basePath); err != nil {
			return err
		}
	}

	if c.UI.CommitsPageSize == 0 {
		c.UI.CommitsPageSize = 30
	}

	if c.Server.Host == "" {
		c.Server.Host = "localhost"
	}

	if c.Server.Port == 0 {
		c.Server.Port = 5555
	} else if c.Server.Port > math.MaxUint16 {
		return fmt.Errorf("server.port should be in 0 < x <= %d range", math.MaxUint16)
	}

	return nil
}

func resolvePath(target string, basePath string) (string, error) {
	if filepath.IsAbs(target) {
		return target, nil
	}

	return filepath.Abs(filepath.Join(basePath, target))
}

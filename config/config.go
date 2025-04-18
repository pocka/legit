package config

import (
	"fmt"
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
	Server struct {
		Name string `yaml:"name,omitempty"`
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
}

func Read(f string) (*Config, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	c := Config{}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if c.Repo.ScanPath, err = resolvePath(c.Repo.ScanPath, f); err != nil {
		return nil, err
	}
	if c.Dirs.Templates, err = resolvePath(c.Dirs.Templates, f); err != nil {
		return nil, err
	}
	if c.Dirs.Static, err = resolvePath(c.Dirs.Static, f); err != nil {
		return nil, err
	}

	return &c, nil
}

func resolvePath(target string, configPath string) (string, error) {
	if filepath.IsAbs(target) {
		return target, nil
	}

	dir := filepath.Dir(configPath)

	return filepath.Abs(filepath.Join(dir, target))
}

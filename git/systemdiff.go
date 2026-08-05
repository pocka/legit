package git

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/pocka/legit/git/exe"
)

// SystemDiff generates diff lines using system git.
func SystemDiff(repoPath string, commit *object.Commit) (*NiceDiff, error) {
	var err error
	var parent *object.Commit
	if parent, err = commit.Parent(0); err != nil {
		if err != object.ErrParentNotFound {
			return nil, fmt.Errorf("unable to read parent of %s: %w", commit.Hash, err)
		}
	}

	nd := NiceDiff{}
	nd.Commit = commit
	nd.Parent = parent
	if parent == nil {
		nd.Parent = &object.Commit{}
	}

	diffs, err := getSystemGitDiff(repoPath, commit)
	if err == nil {
		nd.Files = diffs

		for _, d := range diffs {
			for _, tf := range d.TextFragments {
				for _, l := range tf.Lines {
					switch l.Op {
					case gitdiff.OpAdd:
						nd.Stat.Insertions += 1
					case gitdiff.OpDelete:
						nd.Stat.Deletions += 1
					}
				}
			}
		}

		nd.Stat.FilesChanged = len(diffs)
	} else {
		log.Printf("diff generation failed: %s", err)
	}

	return &nd, nil
}

func getSystemGitDiff(repoPath string, commit *object.Commit) ([]*gitdiff.File, error) {
	cmd := exec.Command(exe.GitBin(), []string{
		// repo.scanPath is trusted directory, thus g.path (inside repo.scanPath) is
		// trusted too.
		// https://git-scm.com/docs/git-config#Documentation/git-config.txt-safedirectory
		"-c", fmt.Sprintf("safe.directory=%s", repoPath),
		"-C", repoPath,
		"show", commit.Hash.String(),
		"--no-ext-diff",
	}...)

	cmd.Env = []string{
		// By default, system git tries to read from various places.
		// We want deterministic, configuration independent behavior here.
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_SYSTEM=/dev/null",
		"GIT_CONFIG_NOSYSTEM=1",
		fmt.Sprintf("GIT_CEILING_DIRECTORIES=%s", repoPath),
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe init failure: %w", err)
	}
	defer stdout.Close()
	stderr := strings.Builder{}
	cmd.Stderr = &stderr
	cmd.Stdin = &bytes.Reader{}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("unable to start system git-diff: %w", err)
	}

	diffs, _, err := gitdiff.Parse(stdout)
	if err != nil {
		return nil, fmt.Errorf("unable to parse git diff: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Print(stderr.String())
		return nil, fmt.Errorf("unable to run system git-diff: %w", err)
	}

	return diffs, nil
}

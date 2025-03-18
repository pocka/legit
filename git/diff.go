package git

import (
	"fmt"
	"log"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// A nicer git diff representation.
type NiceDiff struct {
	Commit *object.Commit
	Parent *object.Commit
	Stat   struct {
		FilesChanged int
		Insertions   int
		Deletions    int
	}
	Files []*gitdiff.File
}

func (g *GitRepo) Diff() (*NiceDiff, error) {
	c, err := g.r.CommitObject(g.h)
	if err != nil {
		return nil, fmt.Errorf("commit object: %w", err)
	}

	patch := &object.Patch{}
	commitTree, err := c.Tree()
	parent := &object.Commit{}
	if err == nil {
		parentTree := &object.Tree{}
		if c.NumParents() != 0 {
			parent, err = c.Parents().Next()
			if err == nil {
				parentTree, err = parent.Tree()
				if err == nil {
					patch, err = parentTree.Patch(commitTree)
					if err != nil {
						return nil, fmt.Errorf("patch: %w", err)
					}
				}
			}
		} else {
			patch, err = parentTree.Patch(commitTree)
			if err != nil {
				return nil, fmt.Errorf("patch: %w", err)
			}
		}
	}

	diffs, _, err := gitdiff.Parse(strings.NewReader(patch.String()))
	if err != nil {
		log.Println(err)
	}

	nd := NiceDiff{}
	nd.Commit = c
	nd.Parent = parent
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

	return &nd, nil
}

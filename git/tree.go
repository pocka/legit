package git

import (
	"github.com/go-git/go-git/v5/plumbing/object"
)

func NewNiceTree(tree *object.Tree, path string) ([]NiceTree, error) {
	files := []NiceTree{}

	if path == "" {
		files = makeNiceTree(tree)
	} else {
		o, err := tree.FindEntry(path)
		if err != nil {
			return nil, err
		}

		if !o.Mode.IsFile() {
			subtree, err := tree.Tree(path)
			if err != nil {
				return nil, err
			}

			files = makeNiceTree(subtree)
		}
	}

	return files, nil
}

// A nicer git tree representation.
type NiceTree struct {
	Name      string
	Mode      string
	Size      int64
	IsFile    bool
	IsSubtree bool
}

func makeNiceTree(t *object.Tree) []NiceTree {
	nts := make([]NiceTree, 0, len(t.Entries))

	for _, e := range t.Entries {
		mode, _ := e.Mode.ToOSFileMode()
		sz, _ := t.Size(e.Name)
		nts = append(nts, NiceTree{
			Name:   e.Name,
			Mode:   mode.String(),
			IsFile: e.Mode.IsFile(),
			Size:   sz,
		})
	}

	return nts
}

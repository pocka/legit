package git

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

type TagList struct {
	refs []*TagReference
	r    *git.Repository
}

// TagReference is used to list both tag and non-annotated tags.
// Non-annotated tags should only contains a reference.
// Annotated tags should contain its reference and its tag information.
type TagReference struct {
	ref *plumbing.Reference
	tag *object.Tag
}

// infoWrapper wraps the property of a TreeEntry so it can export fs.FileInfo
// to tar WriteHeader
type infoWrapper struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
}

func (self *TagList) Len() int {
	return len(self.refs)
}

func (self *TagList) Swap(i, j int) {
	self.refs[i], self.refs[j] = self.refs[j], self.refs[i]
}

// sorting tags in reverse chronological order
func (self *TagList) Less(i, j int) bool {
	var dateI time.Time
	var dateJ time.Time

	if self.refs[i].tag != nil {
		dateI = self.refs[i].tag.Tagger.When
	} else {
		c, err := self.r.CommitObject(self.refs[i].ref.Hash())
		if err != nil {
			dateI = time.Now()
		} else {
			dateI = c.Committer.When
		}
	}

	if self.refs[j].tag != nil {
		dateJ = self.refs[j].tag.Tagger.When
	} else {
		c, err := self.r.CommitObject(self.refs[j].ref.Hash())
		if err != nil {
			dateJ = time.Now()
		} else {
			dateJ = c.Committer.When
		}
	}

	return dateI.After(dateJ)
}

// GitwebDescription returns description text.
// See https://git-scm.com/docs/gitweb#Documentation/gitweb.txt-descriptionorgitwebdescription
func GitwebDescription(repo *git.Repository) string {
	if storage, ok := repo.Storer.(*filesystem.Storage); ok {
		file, err := storage.Filesystem().Open("description")
		if err == nil {
			defer file.Close()
			contents, err := io.ReadAll(file)
			if err == nil {
				text := string(contents)
				// git by default copies "description" file from template directory,
				// and there is no way to detect it other than this "heuristic".
				if strings.Index(text, "Unnamed repository;") != 0 {
					return text
				}
			}
		}
	}

	config, err := repo.Config()
	if err != nil {
		return ""
	}

	gitweb := config.Raw.Section("gitweb")
	if gitweb == nil {
		return ""
	}

	return gitweb.Option("description")
}

// GitwebCategory returns category text.
// See https://git-scm.com/docs/gitweb#Documentation/gitweb.txt-categoryorgitwebcategory
func GitwebCategory(repo *git.Repository) string {
	if storage, ok := repo.Storer.(*filesystem.Storage); ok {
		file, err := storage.Filesystem().Open("category")
		if err == nil {
			defer file.Close()
			contents, err := io.ReadAll(file)
			if err == nil {
				return string(contents)
			}
		}
	}

	config, err := repo.Config()
	if err != nil {
		return ""
	}

	gitweb := config.Raw.Section("gitweb")
	if gitweb == nil {
		return ""
	}

	return gitweb.Option("category")
}

func Tags(repo *git.Repository) ([]*TagReference, error) {
	iter, err := repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("tag objects: %w", err)
	}

	tags := make([]*TagReference, 0)

	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		obj, err := repo.TagObject(ref.Hash())
		switch err {
		case nil:
			tags = append(tags, &TagReference{
				ref: ref,
				tag: obj,
			})
		case plumbing.ErrObjectNotFound:
			tags = append(tags, &TagReference{
				ref: ref,
			})
		default:
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	tagList := &TagList{r: repo, refs: tags}
	sort.Sort(tagList)
	return tags, nil
}

func Branches(repo *git.Repository) ([]*plumbing.Reference, error) {
	bi, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("branchs: %w", err)
	}

	branches := []*plumbing.Reference{}

	_ = bi.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref)
		return nil
	})

	return branches, nil
}

func FindMainBranch(repo *git.Repository, branches []string) (string, error) {
	for _, b := range branches {
		_, err := repo.ResolveRevision(plumbing.Revision(b))
		if err == nil {
			return b, nil
		}
	}
	return "", fmt.Errorf("unable to find main branch")
}

// WriteTar writes itself from a tree into a binary tar file format.
// prefix is root folder to be appended.
func WriteTar(w io.Writer, prefix string, tree *object.Tree) error {
	tw := tar.NewWriter(w)
	defer tw.Close()

	walker := object.NewTreeWalker(tree, true, nil)
	defer walker.Close()

	name, entry, err := walker.Next()
	for ; err == nil; name, entry, err = walker.Next() {
		info, err := newInfoWrapper(name, prefix, &entry, tree)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := tree.File(name)
			if err != nil {
				return err
			}

			reader, err := file.Blob.Reader()
			if err != nil {
				return err
			}

			_, err = io.Copy(tw, reader)
			if err != nil {
				reader.Close()
				return err
			}
			reader.Close()
		}
	}

	return nil
}

func newInfoWrapper(
	name string,
	prefix string,
	entry *object.TreeEntry,
	tree *object.Tree,
) (*infoWrapper, error) {
	var (
		size  int64
		mode  fs.FileMode
		isDir bool
	)

	if entry.Mode.IsFile() {
		file, err := tree.TreeEntryFile(entry)
		if err != nil {
			return nil, err
		}
		mode = fs.FileMode(file.Mode)

		size, err = tree.Size(name)
		if err != nil {
			return nil, err
		}
	} else {
		isDir = true
		mode = fs.ModeDir | fs.ModePerm
	}

	fullname := path.Join(prefix, name)
	return &infoWrapper{
		name:    fullname,
		size:    size,
		mode:    mode,
		modTime: time.Unix(0, 0),
		isDir:   isDir,
	}, nil
}

func (i *infoWrapper) Name() string {
	return i.name
}

func (i *infoWrapper) Size() int64 {
	return i.size
}

func (i *infoWrapper) Mode() fs.FileMode {
	return i.mode
}

func (i *infoWrapper) ModTime() time.Time {
	return i.modTime
}

func (i *infoWrapper) IsDir() bool {
	return i.isDir
}

func (i *infoWrapper) Sys() any {
	return nil
}

func (t *TagReference) Name() string {
	return t.ref.Name().Short()
}

func (t *TagReference) Message() string {
	if t.tag != nil {
		return t.tag.Message
	}
	return ""
}

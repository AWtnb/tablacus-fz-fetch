package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

func MakeDir(path string) error {
	if f, err := os.Stat(path); err == nil && f.IsDir() {
		return nil
	}
	return os.Mkdir(path, os.ModePerm)
}

type Dir struct {
	Path      string
	Exception string
}

func (d Dir) getChildren() (paths []string) {
	fs, err := os.ReadDir(d.Path)
	if err != nil {
		return
	}
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".ini") || strings.HasPrefix(f.Name(), "~$") {
			continue
		}
		p := filepath.Join(d.Path, f.Name())
		if p == d.Exception {
			continue
		}
		paths = append(paths, p)
	}
	return
}

func (d Dir) SelectItems(query string) (ps []string, err error) {
	paths := d.getChildren()
	if len(paths) < 1 {
		return
	}
	idxs, err := fuzzyfinder.FindMulti(paths, func(i int) string {
		return filepath.Base(paths[i])
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop), fuzzyfinder.WithQuery(query))
	if err != nil {
		return
	}
	for _, i := range idxs {
		ps = append(ps, paths[i])
	}
	return
}

func (d Dir) ShowResult() {
	left := d.getChildren()
	if len(left) < 1 {
		fmt.Printf("No items left on '%s'.\n", d.Path)
		return
	}
	if len(left) == 1 {
		fmt.Printf("Left item on '%s':\n- '%s'\n", d.Path, filepath.Base(left[0]))
		return
	}
	fmt.Printf("Left items on '%s':\n", d.Path)
	for i, l := range left {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(left), filepath.Base(l))
	}
}

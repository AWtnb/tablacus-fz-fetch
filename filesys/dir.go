package filesys

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

type Dir struct {
	Path string
}

func (d Dir) getChildren() (ps []string) {
	fs, err := os.ReadDir(d.Path)
	if err != nil {
		return
	}
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".ini") || strings.HasPrefix(f.Name(), "~$") {
			continue
		}
		p := filepath.Join(d.Path, f.Name())
		ps = append(ps, p)
	}
	return
}

func (d Dir) getFiles() (ps []string) {
	items := d.getChildren()
	for _, p := range items {
		if f, err := os.Stat(p); err == nil && !f.IsDir() {
			ps = append(ps, p)
		}
	}
	return
}

func (d Dir) getDirs() (ps []string) {
	items := d.getChildren()
	for _, p := range items {
		if f, err := os.Stat(p); err == nil && f.IsDir() {
			ps = append(ps, p)
		}
	}
	return
}

func (d Dir) SelectItems(file bool, dir bool) (ps []string, err error) {
	var paths []string
	if file {
		paths = append(paths, d.getFiles()...)
	}
	if dir {
		paths = append(paths, d.getDirs()...)
	}
	if len(paths) < 1 {
		return
	}
	idxs, err := fuzzyfinder.FindMulti(paths, func(i int) string {
		return filepath.Base(paths[i])
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
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
		fmt.Printf("Left item on '%s':\n- '%s'", d.Path, left[0])
		return
	}
	fmt.Printf("Left items on '%s':\n", d.Path)
	for i, p := range left {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(left), filepath.Base(p))
	}
}

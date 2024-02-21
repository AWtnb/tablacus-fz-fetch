package filesys

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

type Entry struct {
	path string
}

func (et Entry) name() string {
	return filepath.Base(et.path)
}

type Dir struct {
	Path      string
	Exception string
}

func (d Dir) getChildren() (ets []Entry) {
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
		et := Entry{path: p}
		ets = append(ets, et)
	}
	return
}

func (d Dir) SelectItems() (ps []string, err error) {
	ets := d.getChildren()
	if len(ets) < 1 {
		return
	}
	idxs, err := fuzzyfinder.FindMulti(ets, func(i int) string {
		return ets[i].name()
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
	if err != nil {
		return
	}
	for _, i := range idxs {
		ps = append(ps, ets[i].path)
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
		fmt.Printf("Left item on '%s':\n- '%s'\n", d.Path, left[0].name())
		return
	}
	fmt.Printf("Left items on '%s':\n", d.Path)
	for i, l := range left {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(left), l.name())
	}
}

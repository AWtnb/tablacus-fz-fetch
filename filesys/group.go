package filesys

import (
	"fmt"
	"os"

	"github.com/AWtnb/go-dircopy"
)

type Group struct {
	entries []Entry
}

func (g *Group) Register(paths []string) {
	for _, p := range paths {
		ent := Entry{path: p}
		g.entries = append(g.entries, ent)
	}
}

func (g Group) PreExists(dest string) (paths []string) {
	for _, ent := range g.entries {
		if ent.existsOn(dest) {
			paths = append(paths, ent.path)
		}
	}
	return
}

func (g Group) Size() int {
	return len(g.entries)
}

func (g *Group) Drop(path string) {
	var ents []Entry
	for _, ent := range g.entries {
		if ent.path != path {
			ents = append(ents, ent)
		}
	}
	g.entries = ents
}

func (g Group) CopyTo(dest string) error {
	for _, ent := range g.entries {
		if ent.isDir() {
			np := ent.reborn(dest)
			if err := dircopy.Copy(ent.path, np); err != nil {
				return err
			}
			continue
		}
		if err := ent.copyTo(dest); err != nil {
			return err
		}
	}
	return nil
}

func (g Group) Show() {
	for i, ent := range g.entries {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(g.entries), ent.name())
	}
}

func (g Group) Remove() error {
	for _, ent := range g.entries {
		if ent.isDir() {
			if err := os.RemoveAll(ent.path); err != nil {
				return err
			}
			continue
		}
		if err := os.Remove(ent.path); err != nil {
			return err
		}
	}
	return nil
}

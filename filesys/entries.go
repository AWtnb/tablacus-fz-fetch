package filesys

import (
	"fmt"
	"os"

	"github.com/AWtnb/tablacus-fz-fetch/dir"
)

type Entries struct {
	entries []Entry
}

func (ets *Entries) Register(paths []string) {
	for _, p := range paths {
		c := Entry{path: p}
		ets.entries = append(ets.entries, c)
	}
}

func (ets Entries) Dupls(dest string) (dupls []string) {
	for _, ent := range ets.entries {
		if ent.existsOn(dest) {
			dupls = append(dupls, ent.path)
		}
	}
	return
}

func (ets Entries) Count() int {
	return len(ets.entries)
}

func (ets *Entries) Drop(path string) {
	var ents []Entry
	for _, ent := range ets.entries {
		if ent.path != path {
			ents = append(ents, ent)
		}
	}
	ets.entries = ents
}

func (ets Entries) CopyTo(dest string) error {
	for _, ent := range ets.entries {
		if ent.isDir() {
			np := ent.reborn(dest)
			return dir.Copy(ent.path, np)
		}
		if err := ent.copyTo(dest); err != nil {
			return err
		}
	}
	return nil
}

func (ets Entries) Show() {
	for i, ent := range ets.entries {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(ets.entries), ent.name())
	}
}

func (ets Entries) Remove() error {
	for _, ent := range ets.entries {
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

package filesys

import (
	"fmt"
)

type Entries struct {
	entries []Entry
}

func (es *Entries) Register(path string) {
	ent := Entry{Path: path}
	es.entries = append(es.entries, ent)
}

func (es *Entries) RegisterMulti(paths []string) {
	for _, p := range paths {
		es.Register(p)
	}
}

func (es Entries) UnMovable(dest string) (paths []string) {
	for _, ent := range es.entries {
		if ent.ExistsOn(dest) {
			paths = append(paths, ent.Path)
		}
	}
	return
}

func (es Entries) Size() int {
	return len(es.entries)
}

func (es *Entries) Exclude(path string) {
	var ents []Entry
	for _, ent := range es.entries {
		if ent.Path != path {
			ents = append(ents, ent)
		}
	}
	es.entries = ents
}

func (es Entries) CopyTo(dest string) error {
	for _, ent := range es.entries {
		if err := ent.CopyTo(dest); err != nil {
			return err
		}
	}
	return nil
}

func (es Entries) Show() {
	for i, ent := range es.entries {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(es.entries), ent.Name())
	}
}

func (es Entries) Remove() error {
	for _, ent := range es.entries {
		if err := ent.Remove(); err != nil {
			return err
		}
	}
	return nil
}

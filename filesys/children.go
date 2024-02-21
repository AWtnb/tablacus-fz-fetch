package filesys

import (
	"fmt"
	"os"
	"path/filepath"
)

type Children struct {
	Paths []string
}

func (cs Children) Dupls(dest string) (dupls []string) {
	for _, p := range cs.Paths {
		f := Child{Path: p}
		if f.ExistsOn(dest) {
			dupls = append(dupls, p)
		}
	}
	return
}

func (cs *Children) Drop(path string) {
	var ps []string
	for _, p := range cs.Paths {
		if p != path {
			ps = append(ps, p)
		}
	}
	cs.Paths = ps
}

func (cs Children) Scan(dest string) (uniques []string, dupls []string) {
	for _, p := range cs.Paths {
		f := Child{Path: p}
		if f.ExistsOn(dest) {
			dupls = append(dupls, p)
		} else {
			uniques = append(uniques, p)
		}
	}
	return
}

func (cs Children) CopyTo(dest string) error {
	for _, p := range cs.Paths {
		c := Child{Path: p}
		if err := c.CopyTo(dest); err != nil {
			return err
		}
	}
	return nil
}

func (cs Children) Show() {
	for i, p := range cs.Paths {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(cs.Paths), filepath.Base(p))
	}
}

func (cs Children) Remove() error {
	for _, p := range cs.Paths {
		fi, err := os.Stat(p)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			if err := os.RemoveAll(p); err != nil {
				return err
			}
			continue
		}
		if err := os.Remove(p); err != nil {
			return err
		}
	}
	return nil
}

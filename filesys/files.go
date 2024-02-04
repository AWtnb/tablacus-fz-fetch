package filesys

import (
	"fmt"
	"os"
	"path/filepath"
)

type Files struct {
	Paths []string
}

func (fs Files) GetDuplicates(dest string) (paths []string) {
	for _, path := range fs.Paths {
		f := File{Path: path}
		if f.existsOn(dest) {
			paths = append(paths, path)
		}
	}
	return
}

func (fs Files) GetNonDuplicates(dest string) (paths []string) {
	for _, path := range fs.Paths {
		f := File{Path: path}
		if !f.existsOn(dest) {
			paths = append(paths, path)
		}
	}
	return
}

func (fs Files) CopyFiles(dest string) error {
	for _, path := range fs.Paths {
		sf := File{Path: path}
		if err := sf.copyTo(dest); err != nil {
			return err
		}
	}
	return nil
}

func (fs Files) Show() {
	for i, path := range fs.Paths {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(fs.Paths), filepath.Base(path))
	}
}

func (fs Files) RemoveFiles() error {
	for _, path := range fs.Paths {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

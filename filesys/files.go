package filesys

import (
	"fmt"
	"os"
	"path/filepath"
)

type Files struct {
	Paths []string
}

func (fs Files) CopyFiles(dest string) (result []string, err error) {
	for _, path := range fs.Paths {
		sf := File{Path: path}
		if sf.existsOn(dest) {
			p := fmt.Sprintf("Name duplicated: '%s'\noverwrite?", sf.name())
			a := Asker{Prompt: p, Accept: "y", Reject: "n"}
			if !a.Accepted() {
				fmt.Println("==> skipped")
				continue
			}
		}
		if err = sf.copyTo(dest); err != nil {
			return
		}
		result = append(result, path)
	}
	return
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

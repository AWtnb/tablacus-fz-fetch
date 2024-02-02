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
		sf := File{path: path}
		if sf.existsOn(dest) {
			p := fmt.Sprintf("Name duplicated: '%s'\noverwrite?", sf.name())
			a := Asker{prompt: p, accept: "y", reject: "n"}
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

func (fs Files) show() {
	for i, path := range fs.Paths {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(fs.Paths), filepath.Base(path))
	}
}

func (fs Files) RemoveFiles() error {
	fs.show()
	p := "\nsuccessfully copied everything.\nDELETE original?"
	a := Asker{prompt: p, accept: "y", reject: "n"}
	if !a.Accepted() {
		return nil
	}
	for _, path := range fs.Paths {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

package filesys

import (
	"io"
	"os"
	"path/filepath"
)

type File struct {
	Path string
}

func (f File) name() string {
	return filepath.Base(f.Path)
}

func (f File) existsOn(dirPath string) bool {
	p := filepath.Join(dirPath, filepath.Base(f.Path))
	_, err := os.Stat(p)
	return err == nil
}

func (f File) copyTo(dest string) error {
	srcFile, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	newPath := filepath.Join(dest, filepath.Base(f.Path))
	newFile, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer newFile.Close()
	if _, err = io.Copy(newFile, srcFile); err != nil {
		return err
	}
	return nil
}

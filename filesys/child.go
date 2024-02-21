package filesys

import (
	"io"
	"os"
	"path/filepath"
)

type Child struct {
	Path string
}

func (c Child) newPath(dest string) string {
	return filepath.Join(dest, filepath.Base(c.Path))
}

func (c Child) ExistsOn(dirPath string) bool {
	p := c.newPath(dirPath)
	_, err := os.Stat(p)
	return err == nil
}

func (c Child) copyFile(dest string) error {
	srcFile, err := os.Open(c.Path)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	newPath := c.newPath(dest)
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

func (c Child) CopyTo(dest string) error {
	fi, err := os.Stat(c.Path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		np := c.newPath(dest)
		return CopyDir(c.Path, np)
	}
	return c.copyFile(dest)
}

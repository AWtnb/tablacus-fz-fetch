package filesys

import (
	"io"
	"os"
	"path/filepath"
)

type Entry struct {
	path string
}

func (et Entry) name() string {
	return filepath.Base(et.path)
}

func (et Entry) isDir() bool {
	fi, err := os.Stat(et.path)
	return err == nil && fi.IsDir()
}

func (et Entry) reborn(dest string) string {
	return filepath.Join(dest, filepath.Base(et.path))
}

func (et Entry) existsOn(dirPath string) bool {
	p := et.reborn(dirPath)
	_, err := os.Stat(p)
	return err == nil
}

func (et Entry) copyTo(dest string) error {
	srcFile, err := os.Open(et.path)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	newPath := et.reborn(dest)
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

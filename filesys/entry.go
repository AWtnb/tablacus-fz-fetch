package filesys

import (
	"io"
	"os"
	"path/filepath"

	"github.com/AWtnb/go-dircopy"
)

func copyFile(src string, newPath string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	nf, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer nf.Close()
	if _, err = io.Copy(nf, sf); err != nil {
		return err
	}
	return nil
}

type Entry struct {
	Path string
}

func (e Entry) Name() string {
	return filepath.Base(e.Path)
}

func (e Entry) isDir() bool {
	fi, err := os.Stat(e.Path)
	return err == nil && fi.IsDir()
}

func (e Entry) reborn(dest string) string {
	return filepath.Join(dest, filepath.Base(e.Path))
}

func (e Entry) ExistsOn(dirPath string) bool {
	p := e.reborn(dirPath)
	_, err := os.Stat(p)
	return err == nil
}

func (e Entry) CopyTo(dest string) error {
	fs, err := os.Stat(e.Path)
	if err != nil {
		return err
	}

	newPath := e.reborn(dest)
	if fs.IsDir() {
		return dircopy.Copy(e.Path, newPath)
	}

	return copyFile(e.Path, newPath)
}

func (e Entry) Remove() error {
	if e.isDir() {
		return os.RemoveAll(e.Path)
	}
	return os.Remove(e.Path)
}

func (e Entry) Member() (entries []Entry) {
	if !e.isDir() {
		return
	}
	fs, err := os.ReadDir(e.Path)
	if err != nil {
		return
	}
	for _, f := range fs {
		p := filepath.Join(e.Path, f.Name())
		ent := Entry{Path: p}
		entries = append(entries, ent)
	}
	return
}

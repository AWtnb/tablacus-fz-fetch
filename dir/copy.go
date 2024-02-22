package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Copy(src string, newPath string) error {
	if _, err := os.Stat(newPath); err == nil {
		err := os.RemoveAll(newPath)
		if err != nil {
			return fmt.Errorf("failed to remove already-existing dest path")
		}
	}
	if strings.HasPrefix(newPath, src) {
		return fmt.Errorf("danger to invoke infinit-loop")
	}
	return copy(src, newPath)
}

func isLink(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink != 0 || fi.Mode()&os.ModeDevice != 0
}

func copy(src string, newPath string) error {
	if isLink(src) {
		return fmt.Errorf("'%s' is a link to atnother location", src)
	}
	fs, err := os.Stat(src)
	if err != nil {
		return err
	}

	if fs.IsDir() {
		if err := addDir(src, newPath); err != nil {
			return err
		}
	} else {
		if err := addFile(src, newPath); err != nil {
			return err
		}
	}

	return nil
}

func addDir(src string, newPath string) error {
	if err := os.Mkdir(newPath, 0700); err != nil {
		return err
	}

	fi, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, f := range fi {
		sp := filepath.Join(src, f.Name())
		np := filepath.Join(newPath, f.Name())
		err := copy(sp, np)
		if err != nil {
			return err
		}
	}

	return nil
}

func addFile(src string, newPath string) error {
	d, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	df, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer df.Close()

	if _, err = df.Write(d); err != nil {
		return err
	}

	return nil
}

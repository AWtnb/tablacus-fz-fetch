package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		src  string
		dest string
	)
	flag.StringVar(&src, "src", "", "location of files to copy")
	flag.StringVar(&dest, "dest", "", "destination to copy files")
	flag.Parse()
	if len(src) < 1 {
		src = os.ExpandEnv(`C:\Users\${USERNAME}\Desktop`)
	}
	os.Exit(run(src, dest))
}

func report(s string) {
	fmt.Printf("ERROR: %s\n", s)
	fmt.Scanln()
}

func askUser(prompt string) string {
	fmt.Printf("%s (y/N) ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

type Dir struct {
	path string
}

func (d Dir) getChildren() (ps []string) {
	fs, err := os.ReadDir(d.path)
	if err != nil {
		return
	}
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".ini") || strings.HasPrefix(f.Name(), "~$") || f.IsDir() {
			continue
		}
		p := filepath.Join(d.path, f.Name())
		ps = append(ps, p)
	}
	return
}

func (d Dir) selectFiles() (ps []string, err error) {
	paths := d.getChildren()
	if len(paths) < 1 {
		return
	}
	idxs, err := fuzzyfinder.FindMulti(paths, func(i int) string {
		return filepath.Base(paths[i])
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
	if err != nil {
		return
	}
	for _, i := range idxs {
		ps = append(ps, paths[i])
	}
	return
}

func (d Dir) showLeftFiles() {
	left := d.getChildren()
	fmt.Printf("\n[FINISHED] ")
	if len(left) < 1 {
		fmt.Printf("No files left on '%s'.\n", d.path)
	} else {
		fmt.Printf("Left file(s) on '%s':\n", d.path)
		for _, p := range left {
			fmt.Printf("- '%s'\n", filepath.Base(p))
		}
	}
	fmt.Scanln()
}

type File struct {
	path string
}

func (f File) name() string {
	return filepath.Base(f.path)
}

func (f File) existsOn(dirPath string) bool {
	p := filepath.Join(dirPath, filepath.Base(f.path))
	_, err := os.Stat(p)
	return err == nil
}

func (f File) copyTo(dest string) error {
	srcFile, err := os.Open(f.path)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	newPath := filepath.Join(dest, filepath.Base(f.path))
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

type Files struct {
	paths []string
}

func (fs Files) copyFiles(dest string) (result []string, err error) {
	for _, path := range fs.paths {
		sf := File{path: path}
		if sf.existsOn(dest) {
			msg := fmt.Sprintf("Name duplicated: '%s'\noverwrite?", sf.name())
			ans := askUser(msg)
			if strings.ToLower(ans) != "y" {
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
	for i, path := range fs.paths {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(fs.paths), filepath.Base(path))
	}
}

func (fs Files) removeFiles() error {
	fs.show()
	ans := askUser("\nsuccessfully copied everything.\nDELETE original?")
	if strings.ToLower(ans) != "y" {
		return nil
	}
	for _, path := range fs.paths {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

func run(src string, dest string) int {
	if src == dest {
		return 1
	}
	if src == ".." {
		src = filepath.Dir(dest)
	}
	d := Dir{path: src}
	selected, err := d.selectFiles()
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			report(err.Error())
		}
		return 1
	}

	sfs := Files{paths: selected}
	copied, err := sfs.copyFiles(dest)
	if err != nil {
		report(err.Error())
		return 1
	}
	if len(copied) < 1 {
		return 0
	}
	dupls := Files{paths: copied}
	if err := dupls.removeFiles(); err != nil {
		report(err.Error())
	}
	d.showLeftFiles()
	return 0
}

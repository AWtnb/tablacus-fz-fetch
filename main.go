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

func isValidDirPath(path string) bool {
	s, err := os.Stat(path)
	return err == nil && s.IsDir()
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

func removeFile(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
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
		fmt.Printf("No files left on '%s'.\n", d)
	} else {
		fmt.Printf("Left file(s) on '%s':\n", d)
		for _, p := range left {
			fmt.Printf("- '%s'\n", filepath.Base(p))
		}
	}
	fmt.Scanln()
}

type SelectedFile struct {
	path string
}

func (sf SelectedFile) name() string {
	return filepath.Base(sf.path)
}

func (sf SelectedFile) existsOn(dirPath string) bool {
	p := filepath.Join(dirPath, filepath.Base(sf.path))
	_, err := os.Stat(p)
	return err == nil
}

func (sf SelectedFile) copyTo(dest string) error {
	srcFile, err := os.Open(sf.path)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	newPath := filepath.Join(dest, filepath.Base(sf.path))
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

type SelectedFiles struct {
	paths []string
}

func (sfs SelectedFiles) copyFiles(dest string) (result []string, err error) {
	for _, path := range sfs.paths {
		sf := SelectedFile{path: path}
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

func run(src string, dest string) int {
	if src == dest {
		return 1
	}
	if src == ".." {
		src = filepath.Dir(dest)
	}
	if !isValidDirPath(src) {
		report(fmt.Sprintf("invalid source-path: '%s'", src))
		return 1
	}
	if !isValidDirPath(dest) {
		report(fmt.Sprintf("invalid destination-path: '%s'", dest))
		return 1
	}
	d := Dir{path: src}
	selected, err := d.selectFiles()
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			report(err.Error())
		}
		return 1
	}

	sfs := SelectedFiles{paths: selected}
	copied, err := sfs.copyFiles(dest)
	if err != nil {
		report(err.Error())
		return 1
	}
	if len(copied) < 1 {
		return 0
	}
	if err := askRemoveFiles(copied); err != nil {
		report(err.Error())
	}
	d.showLeftFiles()
	return 0
}

func askRemoveFiles(paths []string) error {
	for i, path := range paths {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(paths), filepath.Base(path))
	}
	fmt.Printf("\nsuccessfully copied everything.\n")
	ans := askUser("DELETE original?")
	if strings.ToLower(ans) == "y" {
		for _, path := range paths {
			err := removeFile(path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

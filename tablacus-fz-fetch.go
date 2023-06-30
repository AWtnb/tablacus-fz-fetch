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

func run(src string, dest string) int {
	if !isValidDirPath(src) {
		report(fmt.Sprintf("invalid source-path: '%s'", src))
		return 1
	}
	if !isValidDirPath(dest) {
		report(fmt.Sprintf("invalid destination-path: '%s'", dest))
		return 1
	}
	ps := getChildren(src)
	if len(ps) < 1 {
		return 1
	}
	selected, err := selectFilePaths(ps)
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			report(err.Error())
		}
		return 1
	}
	var copied []string
	for _, path := range selected {
		n := filepath.Base(path)
		newPath := filepath.Join(dest, n)
		if isValidPath(newPath) {
			msg := fmt.Sprintf("Name duplicated: '%s'\noverwrite?", n)
			ans := askUser(msg)
			if strings.ToLower(ans) != "y" {
				fmt.Println("==> skipped")
				continue
			}
		}
		err := copyFile(path, newPath)
		if err != nil {
			report(err.Error())
			continue
		}
		copied = append(copied, path)
	}
	if len(copied) < 1 {
		return 0
	}
	fmt.Println()
	if err := removeFiles(copied); err != nil {
		report(err.Error())
	}
	showLeftFiles(src)
	return 0
}

func report(s string) {
	fmt.Printf("ERROR: %s\n", s)
	fmt.Scanln()
}

func askUser(prompt string) string {
	fmt.Printf("%s (y/n) ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func isValidPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isValidDirPath(path string) bool {
	s, err := os.Stat(path)
	return err == nil && s.IsDir()
}

func getChildren(d string) []string {
	var ps []string
	fs, err := os.ReadDir(d)
	if err != nil {
		return ps
	}
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".ini") || strings.HasPrefix(f.Name(), "~$") || f.IsDir() {
			continue
		}
		p := filepath.Join(d, f.Name())
		ps = append(ps, p)
	}
	return ps
}

func showLeftFiles(d string) {
	left := getChildren(d)
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

func selectFilePaths(paths []string) ([]string, error) {
	var ps []string
	idxs, err := fuzzyfinder.FindMulti(paths, func(i int) string {
		return filepath.Base(paths[i])
	})
	if err != nil {
		return ps, err
	}
	for _, i := range idxs {
		ps = append(ps, paths[i])
	}
	return ps, err
}

func copyFile(src string, newPath string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
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

func removeFile(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}

func removeFiles(paths []string) error {
	for _, path := range paths {
		fmt.Printf("- '%s' successfully copied!\n", filepath.Base(path))
	}
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

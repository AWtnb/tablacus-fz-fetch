package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AWtnb/tablacus-fz-fetch/filesys"
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

func run(src string, dest string) int {
	if src == dest {
		return 1
	}
	if src == ".." {
		src = filepath.Dir(dest)
	}
	d := filesys.Dir{Path: src}
	selected, err := d.SelectItems(true, false)
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			report(err.Error())
		}
		return 1
	}

	sfs := filesys.Files{Paths: selected}
	copied, err := sfs.CopyFiles(dest)
	if err != nil {
		report(err.Error())
		return 1
	}
	if len(copied) < 1 {
		return 0
	}
	disposals := filesys.Files{Paths: copied}
	disposals.Show()
	p := "\nsuccessfully copied everything.\nDELETE original?"
	a := filesys.Asker{Prompt: p, Accept: "y", Reject: "n"}
	if !a.Accepted() {
		return 0
	}
	if err := disposals.RemoveFiles(); err != nil {
		report(err.Error())
		return 1
	}
	fmt.Printf("\n[FINISHED] ")
	d.ShowResult()
	fmt.Scanln()
	return 0
}

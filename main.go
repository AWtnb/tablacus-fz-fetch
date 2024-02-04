package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func border(s string) {
	fmt.Printf("\n======================================\n %s\n======================================\n", strings.ToUpper(s))
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
	targets := sfs.GetNonDuplicates(dest)
	dupls := sfs.GetDuplicates(dest)
	if 0 < len(dupls) {
		for _, dp := range dupls {
			pr := fmt.Sprintf("Name duplicated: '%s'\noverwrite?", filepath.Base(dp))
			a := Asker{Prompt: pr, Accept: "y", Reject: "n"}
			if !a.Accepted() {
				fmt.Printf("==> skipped\n")
			} else {
				targets = append(targets, dp)
			}
		}
	}
	if 0 < len(targets) {
		t := filesys.Files{Paths: targets}
		if err := t.CopyFiles(dest); err != nil {
			report(err.Error())
			return 1
		}
		border("successfully copied everything")
		t.Show()
		p := "==> Delete original?"
		a := Asker{Prompt: p, Accept: "y", Reject: "n"}
		if a.Accepted() {
			if err := t.RemoveFiles(); err != nil {
				report(err.Error())
				return 1
			}
		}
	}
	border("finished")
	d.ShowResult(true, false)
	fmt.Scanln()
	return 0
}

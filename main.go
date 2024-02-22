package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/tablacus-fz-fetch/dir"
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

func report(err error) {
	fmt.Printf("ERROR: %s\n", err.Error())
	fmt.Scanln()
}

func showLabel(heading string, s string) {
	fmt.Printf("\n\n[%s] %s:\n\n", strings.ToUpper(heading), s)
}

func run(src string, dest string) int {
	if src == dest {
		return 1
	}
	if src == ".." {
		src = filepath.Dir(dest)
	}
	d := dir.Dir{Path: src, Exception: dest}
	selected, err := d.SelectItems()
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			report(err)
		}
		return 1
	}

	var ets filesys.Entries
	ets.Register(selected)
	dupls := ets.Dupls(dest)
	if 0 < len(dupls) {
		for _, dp := range dupls {
			pr := fmt.Sprintf("Name duplicated: '%s'\noverwrite?", filepath.Base(dp))
			a := Asker{Prompt: pr, Accept: "y", Reject: "n"}
			if !a.Accepted() {
				fmt.Printf("==> skipped\n")
				ets.Drop(dp)
			}
		}
	}
	if ets.Count() < 1 {
		return 0
	}
	if err := ets.CopyTo(dest); err != nil {
		report(err)
		return 1
	}
	showLabel("done", "successfully copied everything")
	ets.Show()
	p := "\n==> Delete original?"
	a := Asker{Prompt: p, Accept: "y", Reject: "n"}
	if a.Accepted() {
		if err := ets.Remove(); err != nil {
			report(err)
			return 1
		}
	}
	showLabel("finished", "")
	d.ShowResult()
	fmt.Scanln()
	return 0
}

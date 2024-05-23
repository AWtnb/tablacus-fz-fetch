package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/go-asker"
	"github.com/AWtnb/go-filesys"
	"github.com/AWtnb/tablacus-fz-fetch/dir"
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
	var d dir.Dir
	d.Init(src)
	selected, err := d.SelectItems("")
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			report(err)
		}
		return 1
	}

	var fes filesys.Entries
	fes.RegisterMulti(selected)
	dupls := fes.UnMovable(dest)
	if 0 < len(dupls) {
		for _, dp := range dupls {
			a := asker.Asker{Accept: "y", Reject: "n"}
			a.Ask(fmt.Sprintf("Name duplicated: '%s'\noverwrite?", filepath.Base(dp)))
			if !a.Accepted() {
				fmt.Printf("==> skipped\n")
				fes.Exclude(dp)
			}
		}
	}
	if fes.Size() < 1 {
		return 0
	}
	if err := fes.CopyTo(dest); err != nil {
		report(err)
		return 1
	}
	showLabel("done", "successfully copied everything")
	fes.Show()
	a := asker.Asker{Accept: "y", Reject: "n"}
	a.Ask("\n==> Delete original?")
	if a.Accepted() {
		if err := fes.Remove(); err != nil {
			report(err)
			return 1
		}
	}
	showLabel("finished", "")
	dir.ShowDir(src)
	fmt.Scanln()
	return 0
}

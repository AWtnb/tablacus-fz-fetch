package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Asker struct {
	Prompt string
	Accept string
	Reject string
}

func (a Asker) option() string {
	return fmt.Sprintf(" (%s/%s)", a.Accept, strings.ToUpper(a.Reject))
}

func (a Asker) Accepted() bool {
	fmt.Printf(a.Prompt + a.option())
	scn := bufio.NewScanner(os.Stdin)
	scn.Scan()
	s := scn.Text()
	return strings.EqualFold(s, a.Accept)
}

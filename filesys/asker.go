package filesys

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Asker struct {
	prompt string
	accept string
	reject string
}

func (a Asker) option() string {
	return fmt.Sprintf(" (%s/%s)", a.accept, strings.ToUpper(a.reject))
}

func (a Asker) Accepted() bool {
	fmt.Printf(a.prompt + a.option())
	scn := bufio.NewScanner(os.Stdin)
	scn.Scan()
	s := scn.Text()
	return strings.ToLower(s) == a.accept
}

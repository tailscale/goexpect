// process is a simple example of spawning a process from the expect package.
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/google/goterm/term"
	expect "github.com/tailscale/goexpect"
)

const (
	command = `bc -l`
	timeout = 10 * time.Minute
)

var piRE = regexp.MustCompile(`3.14[0-9]*`)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("Usage: process <nr of digits>")
		os.Exit(1)
	}

	if err := os.Setenv("BC_LINE_LENGTH", "0"); err != nil {
		panic(err)
	}

	scale, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	if scale < 3 {
		panic("scale must be at least 3 for this sample to work")
	}

	e, _, err := expect.Spawn(command, -1)
	if err != nil {
		panic(err)
	}

	if err := e.Send("scale=" + strconv.Itoa(scale) + "\n"); err != nil {
		panic(err)
	}
	if err := e.Send("4*a(1)\n"); err != nil {
		panic(err)
	}
	out, match, err := e.Expect(piRE, timeout)
	if err != nil {
		fmt.Printf("e.Expect(%q,%v) failed: %v, out: %q", piRE.String(), timeout, err, out)
		os.Exit(1)
	}

	fmt.Println(term.Bluef("Pi with %d digits: %s", scale, match[0]))
}

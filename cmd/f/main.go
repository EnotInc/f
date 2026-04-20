package main

import (
	"flag"
	"fmt"

	"github.com/EnotInc/f/internal"
)

func main() {
	exitIf := func(err error) { // this is so cool. I just learned I can to this
		if err != nil {
			fmt.Printf("\033[31m%s", err)
			return
		}
	}

	var all bool = false
	var path bool = false

	flag.BoolVar(&all, "a", false, "show all files")
	flag.BoolVar(&path, "p", false, "show path")

	flag.Parse()

	s, err := internal.NewScanner(all, path)
	exitIf(err)

	err = s.Scan()
	exitIf(err)

	if s.Deny {
		fmt.Println("\n\n \033[33mWARRNING:\033[0m Unable to display some files due to lack of permission")
	}
}

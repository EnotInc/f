package main

import (
	"flag"
	"fmt"

	"github.com/EnotInc/f/internal"
)

func main() {
	var all bool = false
	flag.BoolVar(&all, "a", false, "show all files")
	flag.Parse()

	s, err := internal.NewScanner(all)
	if err != nil {
		fmt.Print(err)
		return
	}

	err = s.Scan()
	if err != nil {
		fmt.Print(err)
	}
}

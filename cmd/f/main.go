package main

import (
	"fmt"

	"github.com/EnotInc/f/internal"
)

func main() {
	err := internal.Scan()
	if err != nil {
		fmt.Print(err)
	}
}

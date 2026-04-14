package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

const defaultWidth = 18

const (
	reset = "\033[0m"
	cyan  = "\033[36m"
)

type scanner struct {
	tabs int
	tabW int
}

func newScanner(width int) *scanner {
	tabs := width / defaultWidth
	tabw := width / tabs
	tabs -= 1

	s := scanner{tabs: tabs, tabW: tabw}
	return &s
}

func Scan() error {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("Unable to run f: %s", err)
	}
	s := newScanner(width)

	var i int = 0

	return filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if i == 0 { // skipping '.' directory
			i += 1
			return nil
		}
		// TODO: ignore hidden files
		// TODO: add flag -a to show all files
		if strings.Count(path, string(os.PathSeparator)) > 0 {
			return fs.SkipDir // reading only current directory (working with depth = 0)
		}

		rString := s.toRuneColumn(path)

		color := reset
		var icon string
		if info.IsDir() {
			color = cyan
			if i, ok := folderIcons[info.Name()]; ok {
				icon = i
			} else {
				icon = " " // default directory icon
			}
		} else {
			if i, ok := fileIcons[filepath.Ext(info.Name())]; ok {
				icon = i
			} else {
				icon = " " // defautl file icon
			}
		}

		fmt.Printf(" %s%s%s", color, icon, string(rString))
		if i%s.tabs == 0 {
			fmt.Println()
		}
		i += 1
		return nil
	})
}

func (s *scanner) toRuneColumn(str string) []rune {
	rString := []rune(str)
	if len(rString) > s.tabW-2 { // tab width - space between columns
		rString = rString[:s.tabW-5]                  // removing both the last 3 runes and 2 space symbols
		rString = append(rString, []rune("...  ")...) // adding '...' and returning the spaces
	} else {
		amout := s.tabW - len(rString)
		space := strings.Repeat(" ", amout)
		rString = append(rString, []rune(space)...)
	}
	return rString
}

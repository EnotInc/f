package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

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
	show bool
}

func NewScanner(showHidden bool) (*scanner, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, fmt.Errorf("Unable to run f: %s", err)
	}

	tabs := width / defaultWidth
	tabw := width / tabs
	tabs -= 1

	s := scanner{tabs: tabs, tabW: tabw, show: showHidden}
	return &s, nil
}

func (s *scanner) Scan() error {
	var i int = 0

	return filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if i == 0 { // skipping '.' directory
			i += 1
			return nil
		}
		if strings.Count(path, string(os.PathSeparator)) > 0 {
			return fs.SkipDir // reading only current directory (working with depth = 0)
		}
		if s.isHidden(path) {
			return nil
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
			if i, ok := fileIcons[info.Name()]; ok {
				icon = i // files without extations
			} else if i, ok := fileIcons[filepath.Ext(info.Name())]; ok {
				icon = i // files with extations
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

func (s *scanner) isHidden(path string) bool {
	if s.show {
		return false
	}

	// is hidden on linux
	linux := path[0] == '.'

	// is hidden on windows
	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return linux
	}
	attr, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return linux
	}
	windows := attr&syscall.FILE_ATTRIBUTE_HIDDEN != 0

	return linux || windows
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

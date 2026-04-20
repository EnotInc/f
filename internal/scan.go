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
	dir   = "\033[1;36m"
)

type scanner struct {
	tabs int
	tabW int
	show bool
	path bool
	Deny bool
}

func NewScanner(showHidden bool, path bool) (*scanner, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, fmt.Errorf("Unable to run 'f': %s", err)
	}

	tabs := width / defaultWidth
	tabw := width / tabs

	// Fixin borders
	//
	//  # foobar     # qwerty
	// ^^^          ^^^
	//
	// ' ', icon, ' ' - 3 runes
	tabw -= 3

	s := scanner{
		tabs: tabs,
		tabW: tabw,
		show: showHidden,
		path: path,
		Deny: false,
	}

	return &s, nil
}

func (s *scanner) Scan() error {
	var i int = 0

	if s.path {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Printf(" \033[33mWARRNING:\033[0m Unable to get path: %s\n", err)
		} else {
			fmt.Printf("[\033[34m%s\033[0m]\n", dir)
		}
	}

	return filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				s.Deny = true
				return nil
			}
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
			color = dir
			if _i, ok := folderIcons[info.Name()]; ok {
				icon = _i
			} else {
				icon = " " // default directory icon
			}
		} else {
			name := strings.ToLower(info.Name())
			if _i, ok := fileIcons[name]; ok {
				icon = _i
			} else if _i, ok := fileIcons[filepath.Ext(name)]; ok {
				icon = _i
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

	hidden := path[0] == '.'

	return hidden
}

func (s *scanner) toRuneColumn(str string) []rune {
	if strings.Contains(str, " ") {
		str = fmt.Sprintf("'%s'", str)
	}
	rString := []rune(str)

	end := func(rs []rune) []rune {
		rs = rs[:s.tabW-5]
		rs = append(rs, []rune("...  ")...)
		return rs
	}

	if len(rString) > s.tabW-2 { // tab width - space between columns
		rString = end(rString)
	} else {
		amout := s.tabW - len(rString)
		space := strings.Repeat(" ", amout)
		rString = append(rString, []rune(space)...)
	}
	return rString
}

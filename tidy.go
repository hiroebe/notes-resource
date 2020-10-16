package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mattn/go-zglob"
	"github.com/rhysd/notes-cli"
)

type target struct {
	id, path string
}

func tidy(config *notes.Config) error {
	cats, err := notes.CollectCategories(config, 0)
	if err != nil {
		return err
	}
	pattern, err := regexp.Compile(`\((\S*[0-9a-v]{20}\.\S+)\)`)
	if err != nil {
		return err
	}
	targets := []target{}
	for _, cat := range cats {
		for _, path := range cat.NotePaths {
			matches, err := findMatchInFile(path, pattern)
			if err != nil {
				fmt.Println(err)
			}
			for _, m := range matches {
				fname := m[1]
				if strings.HasPrefix(fname, "http") {
					continue
				}
				targets = append(targets, target{
					id:   extractID(fname),
					path: filepath.Join(path, "..", filepath.FromSlash(fname)),
				})
			}
		}
	}
	for _, t := range targets {
		pattern := fmt.Sprintf("%s/**/*%s%s", config.HomePath, t.id, filepath.Ext(t.path))
		matches, err := zglob.Glob(pattern)
		if err != nil {
			return err
		}
		if len(matches) != 1 {
			return fmt.Errorf("could not find exactly one file by pattern: %s", pattern)
		}
		match := matches[0]
		if match == t.path {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(t.path), 0755); err != nil {
			return err
		}
		if err := os.Rename(match, t.path); err != nil {
			return err
		}
		removeDirRec(filepath.Dir(match))
		fmt.Printf("Successfully moved %s to %s\n", match, t.path)
	}
	return nil
}

func findMatchInFile(path string, pattern *regexp.Regexp) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	matches := [][]string{}
	for scanner.Scan() {
		for _, m := range pattern.FindAllStringSubmatch(scanner.Text(), -1) {
			matches = append(matches, m)
		}
	}
	return matches, nil
}

func removeDirRec(dir string) {
	for {
		// Remove directory if empty
		if err := os.Remove(dir); err != nil {
			break
		}
		dir = filepath.Dir(dir)
	}
}

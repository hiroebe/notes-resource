package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	targets := []target{}
	for _, cat := range cats {
		for _, path := range cat.NotePaths {
			resources, err := findResourcesInFile(path)
			if err != nil {
				fmt.Println(err)
			}
			for _, r := range resources {
				targets = append(targets, target{
					id:   extractID(r),
					path: filepath.Join(path, "..", filepath.FromSlash(r)),
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

func removeDirRec(dir string) {
	for {
		// Remove directory if empty
		if err := os.Remove(dir); err != nil {
			break
		}
		dir = filepath.Dir(dir)
	}
}

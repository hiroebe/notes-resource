package main

import (
	"fmt"

	"github.com/mattn/go-zglob"
	"github.com/rhysd/notes-cli"
)

func listUnused(config *notes.Config, args []string) error {
	cats, err := notes.CollectCategories(config, 0)
	if err != nil {
		return err
	}
	ids := map[string]struct{}{}
	for _, cat := range cats {
		for _, path := range cat.NotePaths {
			resources, err := findResourcesInMarkdown(path)
			if err != nil {
				fmt.Println(err)
			}
			for _, r := range resources {
				ids[r.id] = struct{}{}
			}
		}
	}

	matches, err := zglob.Glob(fmt.Sprintf("%s/**/*", config.HomePath))
	if err != nil {
		return err
	}
	for _, m := range matches {
		r, err := newResource(m)
		if err != nil {
			continue
		}
		if _, ok := ids[r.id]; !ok {
			fmt.Println(m)
		}
	}
	return nil
}

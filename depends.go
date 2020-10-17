package main

import (
	"fmt"
	"os"

	"github.com/rhysd/notes-cli"
)

func listDepends(config *notes.Config, note string) error {
	f, err := os.Open(note)
	if err != nil {
		return err
	}
	defer f.Close()

	resources, err := findResourcesInFile(note)
	if err != nil {
		return err
	}
	for _, r := range resources {
		fmt.Println(r)
	}
	return nil
}

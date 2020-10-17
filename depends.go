package main

import (
	"fmt"
	"os"

	"github.com/rhysd/notes-cli"
)

func listDepends(config *notes.Config, args []string) error {
	if len(args) < 1 {
		return errInvalidArgument
	}
	note := args[0]
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

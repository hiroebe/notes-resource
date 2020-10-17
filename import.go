package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rhysd/notes-cli"
)

func importResources(config *notes.Config, args []string) error {
	if len(args) < 2 {
		return errInvalidArgument
	}
	resources, target := args[:len(args)-1], args[len(args)-1]
	for _, r := range resources {
		if err := importResource(config, r, target); err != nil {
			return err
		}
	}
	return nil
}

func importResource(config *notes.Config, source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	var dir string
	info, err := os.Stat(target)
	if err != nil {
		return err
	}
	if info.IsDir() {
		dir = target
	} else {
		dir = filepath.Dir(target)
	}

	outFile := createResource(filepath.Base(source)).getPath()
	out := filepath.Join(dir, outFile)
	if _, err := os.Stat(out); err == nil {
		return fmt.Errorf("file already exists: %s", out)
	}

	dstFile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

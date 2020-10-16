package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rhysd/notes-cli"
)

func importResources(config *notes.Config, sources []string, target string) error {
	for _, s := range sources {
		if err := importResource(config, s, target); err != nil {
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

	out := filepath.Join(dir, addID(filepath.Base(source)))
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

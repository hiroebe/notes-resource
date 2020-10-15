package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mattn/go-zglob"
	"github.com/rhysd/notes-cli"
	"github.com/rs/xid"
)

var opts struct {
	Tidy bool `short:"t" long:"tidy" description:"Move resources to the proper directory"`
	// TODO: --depends
	// TODO: --prune
}

func help() error {
	_, err := fmt.Println(`Usage:
  notes resource RESOURCE [RESOURCE...] TARGET
  notes resource --tidy`)
	return err
}

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
					path: filepath.Join(path, "..", fname),
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
		if err := os.Rename(match, t.path); err != nil {
			return err
		}
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

func addID(fname string) string {
	ext := filepath.Ext(fname)
	return fmt.Sprintf("%s-%s%s", strings.TrimSuffix(fname, ext), xid.New().String(), ext)
}

func extractID(fname string) string {
	ext := filepath.Ext(fname)
	endIdx := len(fname) - len(ext)
	return fname[endIdx-20 : endIdx]
}

func run() error {
	args, err := flags.NewParser(&opts, flags.IgnoreUnknown).Parse()
	if err != nil {
		return err
	}

	config, err := notes.NewConfig()
	if err != nil {
		return err
	}

	if opts.Tidy {
		return tidy(config)
	}

	// args are like ["notes" "resource" "RESOURCE" "TARGET"]
	args = args[2:]

	if len(args) < 2 {
		return help()
	}

	resources, target := args[:len(args)-1], args[len(args)-1]
	return importResources(config, resources, target)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
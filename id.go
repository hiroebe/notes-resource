package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/xid"
)

var reFileNameWithID = regexp.MustCompile(`\((\S*[0-9a-v]{20}\.\S+)\)`)

func addID(fname string) string {
	ext := filepath.Ext(fname)
	return fmt.Sprintf("%s-%s%s", strings.TrimSuffix(fname, ext), xid.New().String(), ext)
}

func extractID(fname string) string {
	ext := filepath.Ext(fname)
	endIdx := len(fname) - len(ext)
	return fname[endIdx-20 : endIdx]
}

func findResourcesInFile(path string) ([]string, error) {
	matches, err := findMatchInFile(path, reFileNameWithID)
	if err != nil {
		return nil, err
	}
	resources := make([]string, 0, len(matches))
	for _, m := range matches {
		fname := m[1]
		if strings.HasPrefix(fname, "http") {
			continue
		}
		resources = append(resources, fname)
	}
	return resources, nil
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

package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/xid"
)

var (
	reResourceName           = regexp.MustCompile(`\S*[0-9a-v]{20}\.\S+$`)
	reResourceNameInMarkdown = regexp.MustCompile(`\((\S*[0-9a-v]{20}\.\S+)\)`)

	errInvalidResourceName = errors.New("invalid resource name")
)

type resource struct {
	id     string
	prefix string
	ext    string
}

func (r *resource) getName() string {
	return r.prefix + r.id + r.ext
}

func (r *resource) getPath() string {
	return filepath.FromSlash(r.getName())
}

func newResource(resourceName string) (*resource, error) {
	if !reResourceName.MatchString(resourceName) {
		return nil, errInvalidResourceName
	}
	ext := filepath.Ext(resourceName)
	endIdx := len(resourceName) - len(ext)

	return &resource{
		id:     resourceName[endIdx-20 : endIdx],
		prefix: resourceName[:endIdx-20],
		ext:    ext,
	}, nil
}

func createResource(fname string) *resource {
	ext := filepath.Ext(fname)
	prefix := strings.TrimSuffix(fname, ext) + "-"

	return &resource{
		id:     xid.New().String(),
		prefix: prefix,
		ext:    ext,
	}
}

func findResourcesInMarkdown(path string) ([]*resource, error) {
	matches, err := findMatchInFile(path, reResourceNameInMarkdown)
	if err != nil {
		return nil, err
	}
	resources := make([]*resource, 0, len(matches))
	for _, m := range matches {
		fname := m[1]
		if strings.HasPrefix(fname, "http") {
			continue
		}
		r, err := newResource(fname)
		if err != nil {
			continue
		}
		resources = append(resources, r)
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

package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rs/xid"
)

func addID(fname string) string {
	ext := filepath.Ext(fname)
	return fmt.Sprintf("%s-%s%s", strings.TrimSuffix(fname, ext), xid.New().String(), ext)
}

func extractID(fname string) string {
	ext := filepath.Ext(fname)
	endIdx := len(fname) - len(ext)
	return fname[endIdx-20 : endIdx]
}

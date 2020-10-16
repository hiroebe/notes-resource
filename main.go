package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/rhysd/notes-cli"
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

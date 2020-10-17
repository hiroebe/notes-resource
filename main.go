package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/rhysd/notes-cli"
)

var opts struct {
	Tidy    bool `short:"t" long:"tidy" description:"Move resources to the proper directory"`
	Depends bool `short:"d" long:"depends" description:"List resources which the note depends"`
	// TODO: --prune
}

func help() error {
	_, err := fmt.Println(`Usage:
  notes resource RESOURCE [RESOURCE...] TARGET
  notes resource --tidy
  notes resource --depends NOTE`)
	return err
}

func run() error {
	args, err := flags.NewParser(&opts, flags.IgnoreUnknown).Parse()
	if err != nil {
		return err
	}

	// args are like ["notes" "resource" "ARG1" "ARG2"]
	args = args[2:]

	config, err := notes.NewConfig()
	if err != nil {
		return err
	}

	if opts.Tidy {
		return tidy(config)
	}
	if opts.Depends {
		if len(args) < 1 {
			return help()
		}
		return listDepends(config, args[0])
	}

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

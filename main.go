package main

import (
	"errors"
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/rhysd/notes-cli"
)

var opts struct {
	Tidy    bool `short:"t" long:"tidy" description:"Move resources to the proper directory"`
	Depends bool `short:"d" long:"depends" description:"List resources which the note depends"`
	Unused  bool `short:"u" long:"unused" description:"List unused resources"`
}

var errInvalidArgument = errors.New("invalid argument")

func help() error {
	_, err := fmt.Println(`Usage:
  notes resource RESOURCE [RESOURCE...] TARGET
  notes resource --tidy
  notes resource --depends NOTE
  notes resource --unused`)
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
		return tidy(config, args)
	}
	if opts.Depends {
		return listDepends(config, args)
	}
	if opts.Unused {
		return listUnused(config, args)
	}
	return importResources(config, args)
}

func main() {
	if err := run(); err != nil {
		if err == errInvalidArgument {
			help()
			return
		}
		panic(err)
	}
}

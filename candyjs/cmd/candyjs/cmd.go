package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	parser := flags.NewNamedParser("candyjs", flags.Default)
	parser.AddCommand(
		"import",
		"Generates the code necesary to push a package in a Context", "",
		&CmdImport{},
	)

	_, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrCommandRequired {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}

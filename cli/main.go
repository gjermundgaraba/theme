package main

import (
	"fmt"
	"os"

	"github.com/gjermundgaraba/theme/theme-generator"
	"github.com/gjermundgaraba/theme/themes"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: gg-theme <build|link [--dry-run]|serve>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "build":
		if err := themes.Build(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "link":
		opts, err := parseLinkOptions(os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := themes.LinkWithOptions(opts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "serve":
		generator.Serve()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func parseLinkOptions(args []string) (themes.LinkOptions, error) {
	var opts themes.LinkOptions
	for _, arg := range args {
		switch arg {
		case "--dry-run":
			opts.DryRun = true
		default:
			return opts, fmt.Errorf("unknown link option: %s", arg)
		}
	}
	return opts, nil
}

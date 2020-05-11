package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
	"moul.io/godev"
	"moul.io/pkgman/pkg/ipa"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	root := &ffcli.Command{
		ShortUsage: `pkgman SUBCOMMAND`,
		Subcommands: []*ffcli.Command{
			{
				Name:       "ipa-plist-json",
				ShortUsage: "ipa-plist-json PATH",
				Exec:       runIpaPlistJSON,
			},
		},
		Exec: func(context.Context, []string) error { return flag.ErrHelp },
	}
	return root.ParseAndRun(context.Background(), args[1:])
}

func runIpaPlistJSON(_ context.Context, args []string) error {
	if len(args) < 1 {
		return flag.ErrHelp
	}
	for _, arg := range args {
		pkg, err := ipa.Open(arg)
		if err != nil {
			return err
		}
		defer pkg.Close()

		for _, app := range pkg.Apps() {
			plist, err := app.Plist()
			if err != nil {
				return err
			}
			fmt.Println(godev.PrettyJSON(plist))
		}
	}
	return nil
}

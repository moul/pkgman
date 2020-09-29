package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
	"moul.io/pkgman/pkg/apk"
	"moul.io/pkgman/pkg/ipa"
	"moul.io/u"
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
			{Name: "ipa-plist-json", ShortUsage: "ipa-plist-json PATH", Exec: runIpaPlistJSON},
			{Name: "apk-manifest", ShortUsage: "apk-manifest PATH", Exec: runApkManifest},
			{Name: "apk-manifest-xml", ShortUsage: "apk-manifest-xml PATH", Exec: runApkManifestXML},
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
			fmt.Println(u.PrettyJSON(plist))
		}
	}
	return nil
}

func runApkManifest(_ context.Context, args []string) error {
	if len(args) < 1 {
		return flag.ErrHelp
	}
	for _, arg := range args {
		pkg, err := apk.Open(arg)
		if err != nil {
			return err
		}
		defer pkg.Close()

		manifest, err := pkg.Manifest()
		if err != nil {
			return err
		}
		fmt.Println(u.PrettyJSON(manifest))
	}
	return nil
}

func runApkManifestXML(_ context.Context, args []string) error {
	if len(args) < 1 {
		return flag.ErrHelp
	}
	for _, arg := range args {
		pkg, err := apk.Open(arg)
		if err != nil {
			return err
		}
		defer pkg.Close()

		manifest, err := pkg.ManifestXML()
		if err != nil {
			return err
		}
		fmt.Println(manifest)
	}
	return nil
}

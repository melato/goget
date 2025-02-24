package main

import (
	_ "embed"
	"fmt"
	"os"

	"flag"

	"melato.org/goget"
)

//go:embed usage.txt
var usageData []byte

var version string = "dev"

func mainE() error {
	var app goget.App
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var help bool
	fs.BoolVar(&help, "h", false, "help")
	fs.IntVar(&app.Port, "port", 8080, "port to listen to")
	fs.StringVar(&app.ConfigFile, "c", "", "config file (.yaml)")
	fs.StringVar(&app.Template, "template", "", "template file")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	if help {
		fmt.Printf("%s\n", usageData)
		return nil
	}
	args := fs.Args()
	if len(args) == 0 {
		return nil
	}
	cmd := args[0]
	switch cmd {
	case "version":
		fmt.Printf("%s\n", version)
		return nil
	case "template":
		app.PrintTemplate()
		return nil
	}
	err = app.Configured()
	if err != nil {
		return err
	}
	switch cmd {
	case "server":
		return app.Server()
	case "list":
		app.List()
	}
	return nil
}

func main() {
	err := mainE()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

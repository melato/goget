package main

import (
	_ "embed"
	"fmt"

	"melato.org/command"
	"melato.org/project"
)

//go:embed version
var version string

func main() {
	cmd := &command.SimpleCommand{}
	app := &project.App{}
	cmd.Flags(app)
	cmd.Command("list").RunFunc(app.List)
	cmd.Command("server").RunFunc(app.Server)
	cmd.Command("template").NoConfig().RunMethod(app.PrintTemplate)
	cmd.Command("version").NoConfig().RunMethod(func() { fmt.Println(version) }).Short("print program version")
	command.Main(cmd)
}

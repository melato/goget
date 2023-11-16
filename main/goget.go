package main

import (
	_ "embed"
	"fmt"

	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/goget"
)

//go:embed usage.yaml
var usageData []byte

var version string = "dev"

func main() {
	cmd := &command.SimpleCommand{}
	app := &goget.App{}
	cmd.Flags(app)
	cmd.Command("list").RunFunc(app.List)
	cmd.Command("server").RunFunc(app.Server)
	cmd.Command("template").NoConfig().RunMethod(app.PrintTemplate)
	cmd.Command("version").NoConfig().RunMethod(func() { fmt.Println(version) }).Short("print program version")
	usage.Apply(cmd, usageData)
	command.Main(cmd)
}

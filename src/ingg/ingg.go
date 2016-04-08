package main

import (
    "os"
    "github.com/codegangsta/cli"
    "github.com/fatih/color"
)

var inggAppHelpTemplate = `Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .Flags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}
   {{if .Commands}}
Some usefull {{.HelpName}} commands are:
{{range .Categories}}{{if .Name}}
{{.Name}}{{end}}{{range .Commands}}
    {{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}{{end}}
{{end}}{{end}}{{if .Flags}}
Global Options
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

// The text template for the command help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var inggCommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   {{.HelpName}}{{if .Flags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .Flags}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

func main() {

    cli.AppHelpTemplate = inggAppHelpTemplate
    cli.CommandHelpTemplate = inggCommandHelpTemplate

    app := cli.NewApp()
    app.Name = "ingg"
    app.Version = "1.0.0"
    app.Flags = []cli.Flag{
        cli.BoolFlag{
            Name: "no-color, nc",
            Usage: "Disable colored output",
        },
    }

    app.Before = func(c *cli.Context) error {
        noColors := c.Bool("no-color")

        if (noColors) {
            color.NoColor = true
        }

        return nil
    }

    app.Run(os.Args)
}

package main

import (
    "os"
    "github.com/codegangsta/cli"
    "github.com/fatih/color"
    "ingg/cmds"
    "ingg/utils"
    "fmt"
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
var inggCommandHelpTemplate = `Usage: {{.HelpName}}{{if .Flags}} [command options]{{end}} {{if .ArgsUsage}}{{ .ArgsUsage}}{{else}}[arguments...]{{end}}

{{if .Description}}{{.Description}}{{end}}
{{if .Flags}}
Options:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

func main() {

    cli.AppHelpTemplate = inggAppHelpTemplate
    cli.CommandHelpTemplate = inggCommandHelpTemplate

    app := cli.NewApp()
    app.Name = "ingg"
    app.Version = "1.0.0"
    app.EnableBashCompletion = true
    app.Flags = []cli.Flag{
        cli.BoolFlag{
            Name: "no-color, nc",
            Usage: "Disable colored output",
        },
        cli.BoolFlag{
            Name: "verbose, x",
            Usage: "Print more information of the process",
        },
    }

    app.Before = func(c *cli.Context) error {
        noColors := c.Bool("no-color")

        if (noColors) {
            color.NoColor = true
        }

        if c.Bool("verbose") {
            utils.PrintMsg("Verbose mode enabled")
            cmds.Verbose = true
        }

        return nil
    }

    app.Commands = []cli.Command{
        cmds.SvnToGit,
        cmds.MavenBuild,
        {
            Name:  "complete",
            Aliases: []string{"c"},
            Usage: "complete a task on the list",
            Action: func(c *cli.Context) error {
                fmt.Println("completed task: ", c.Args().First())
                return nil
            },
            BashComplete: func(c *cli.Context) {
                // This will complete if no args are passed
                if c.NArg() > 0 {
                    return
                }
                for _, t := range tasks {
                    fmt.Println(t)
                }
            },
        },
    }

    app.Run(os.Args)
}

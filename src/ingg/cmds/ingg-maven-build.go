package cmds

import (
    "github.com/codegangsta/cli"
    "os"
    Print "ingg/utils"
    "bufio"
    "strings"
    "os/exec"
    "fmt"
)

type ConfigBlock struct {
    key     string
    modules []string
}

var MavenBuild = cli.Command{
    Name: "maven-build",
    Category: "Maven Operations",
    ArgsUsage:"--file <configFile> [--category <category>]",
    Description: "Build all maven projects described from a config file",
    Flags: [] cli.Flag{
        cli.StringFlag{
            Name: "file, f",
            Usage: "The config file",
        },
        cli.StringFlag{
            Name: "block, b",
            Usage: "The config block name. Defaults to all",
            Value: "all",
        },
    },
    Action: func(c *cli.Context) {
        var isValid bool = true

        currentDir, _ := os.Getwd();

        file := c.String("file")
        configBlock := c.String("block")

        if file == "" {
            Print.Danger("* File is required")
            isValid = false

        } else {
            _, err := os.Stat(file)
            if err != nil {
                Print.Danger("* Config file not found: ", file)
                isValid = false
            }
        }

        if !isValid {
            Print.PrintMsg("")
            Print.Danger("Requirements not met!")
            os.Exit(1)
        }

        Print.PrintMsg("Starting Building config block " + configBlock + "\n")
        configMap := getConfigMap(file)

        modules := configMap[configBlock]

        Print.PrintInfo(fmt.Sprintf("Config block modules: %v", modules))

        if (len(modules) > 0 ) {
            for _, module := range modules {

                os.Chdir(currentDir + "/" + module);

                cmd := exec.Command("mvn", "clean", "install")
                if Verbose {
                    cmd.Stdout = os.Stdout
                    cmd.Stderr = os.Stderr
                }

                err := cmd.Run()

                if err == nil {
                    Print.PrintSuccess(fmt.Sprintf("* Module build done: %s", module))

                } else {
                    Print.PrintDanger(fmt.Sprintf("* Error building module: %s", module))
                    Print.PrintError(err)
                }
            }
        }

    },
}

func getConfigMap(file string) map[string][]string {
    var configMap map[string][]string = make(map[string][]string)

    openFile, _ := os.Open(file)
    scanner := bufio.NewScanner(openFile)

    configBlock := new(ConfigBlock)
    var line string

    for scanner.Scan() {
        line = strings.TrimSpace(scanner.Text())

        if line != "" {

            // Starting new configuration block
            if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {

                // Saving current state
                if configBlock.key != "" {
                    configMap[configBlock.key] = configBlock.modules;
                    configBlock = new(ConfigBlock)
                }

                configBlock.key = strings.Trim(line, "[]")

            } else {
                configBlock.modules = append(configBlock.modules, line)
            }

        }
    }

    configMap[configBlock.key] = configBlock.modules;

    return configMap
}


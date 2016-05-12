package utils

import (
    "github.com/fatih/color"
    "os"
    "fmt"
)

var Info = color.New(color.FgBlue).PrintlnFunc()
var Success = color.New(color.FgGreen).PrintlnFunc()
var Danger = color.New(color.FgRed).PrintlnFunc()
var Warning = color.New(color.FgHiYellow).PrintlnFunc()

var PrintMsg = func(txt string) {
    fmt.Println(txt)
}

var PrintInfo = func(txt string) {
    Info(txt);
}

var PrintSuccess = func(txt string) {
    Success(txt);
}

var PrintWarning = func(txt string) {
    Warning(txt);
}

var PrintDanger = func(txt string) {
    Danger(txt);
}

// Control the panic output and write
// in the standard error output
func PrintError(err error) {
    if err != nil {
        os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
    }
}

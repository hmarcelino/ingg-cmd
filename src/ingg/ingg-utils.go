package main

import (
    "github.com/fatih/color"
    "os"
    "fmt"
)

var info = color.New(color.FgBlue).PrintlnFunc()
var danger = color.New(color.FgRed).PrintlnFunc()
var warning = color.New(color.FgHiYellow).PrintlnFunc()

var PrintMsg = func(txt string) {
    fmt.Println(txt)
}

var PrintInfo = func(txt string) {
    info(txt);
}

var PrintWarning = func(txt string) {
    warning(txt);
}

var PrintDanger = func(txt string) {
    danger(txt);
}

// Control the panic output and write
// in the standard error output
func PrintError(err error) {
    if err != nil {
        os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
    }
}

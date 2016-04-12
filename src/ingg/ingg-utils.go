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

// Print the bytes array in the
// standard output
func printOutput(outs []byte) {
    if len(outs) > 0 {
        os.Stdout.WriteString(fmt.Sprintf("==> Output:\n%s\n", string(outs)))
    }
}

// Control the panic output and write
// in the standard error output
func printError(err error) {
    if err != nil {
        //os.Stderr.WriteString(color.Red("==> Error: %s\n", err.Error()))
        color.Red("==> Error:\n%s\n", err.Error())
    }
}


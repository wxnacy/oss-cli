package main

import (
    "fmt"
    "github.com/wxnacy/goss/goss"
	"github.com/c-bata/go-prompt"
    "os"
    "strings"
)

var ENDPOINT = "oss-cn-beijing.aliyuncs.com"
var brand goss.IBrand
var terminal *goss.Terminal

func completer(d prompt.Document) []prompt.Suggest {
    s := terminal.Prompt(d.CurrentLine())
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func executor(t string) {
    goss.Log.Infof("Terminal", t)
	if t == "bash" {
		// cmd := exec.Command("bash")
		// cmd.Stdin = os.Stdin
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Run()
        goss.Log.Info(t)
	}
    ts := strings.Split(t, " ")
    switch ts[0] {
        case "exit": {
            goss.Log.Info("Exit")
            os.Exit(0)
        }
        case "ll": {
            terminal.PrintLL()
        }
        // case "ls": {
            // terminal.PrintLS()
        // }
        case "cd": {
            terminal.Cd(ts[1])
        }
        case "get": {
            name := ""
            if len(ts) > 2 {
                name = ts[2]
            }
            terminal.Get(ts[1], name)
        }
        case "post": {
            name := ""
            if len(ts) > 2 {
                name = ts[2]
            }
            terminal.Post(name, ts[1])
        }
        case "..": {
            terminal.Cd("..")
        }
        case "...": {
            terminal.Cd("...")
        }
        case "pwd": {
            fmt.Println(terminal.PWD())
        }
    }
	return
}

func main() {

    goss.Log.Info("Hello World")
    // terminal = goss.NewTerminal("/Brands/oss/Credentials/wxnacy/Buckets/wxnacy-file/Keys")
    terminal = goss.NewTerminal("/Brands/oss/Credentials")
    // terminal = goss.NewTerminal("/")

    p := prompt.New(
        executor,
        completer,
    )
    p.Run()

}


func HandleError(err error) {
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(-1)
    }
}

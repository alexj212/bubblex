package main

import (
    "bufio"
    "bytes"
    "github.com/alexj212/bubblex"
    "github.com/alexj212/gox/commandr"
    "log"
    "os"
    "strings"

    "github.com/potakhov/loge"
)

var console *bubblex.TuiClient

func inputHandler(s string) {
    loge.Info("SetOnEnteredHandler %s", s)

    var b bytes.Buffer
    out := bufio.NewWriter(&b)

    parsed, err := commandr.NewCommandArgs(s, out)
    if err != nil {

        loge.Info("NewCommandArgs parsed: %s\n", b.String())
        loge.Info("NewCommandArgs err: %v\n", err)
        return
    }

    execErr := commandr.DefaultCommands.Execute(console, parsed)

    if execErr != nil {
        loge.Info("Execute: %s\n", b.String())
        loge.Info("Execute Error: %v\n", err)
        return
    }

    loge.Info("Execute: %s\n", string(b.Bytes()))
}

func outputHandler(b []byte) {
    text := string(b)
    lines := strings.Split(text, "\n")
    for _, line := range lines {

        loge.Info("%s", line)
    }
}

func main() {
    cleanup := bubblex.LogeInit()
    defer cleanup()

    go func() {
        spawnSsh()
    }()

    console = bubblex.NewTuiClient(outputHandler, inputHandler)

    for _, cmd := range commandr.DefaultCommands.Commands() {
        if cmd.Use == "cls" || cmd.Use == "exit" {
            commandr.DefaultCommands.RemoveCommand(cmd)
        }
    }

    commandr.DefaultCommands.AddCommand(ExitCommand)
    commandr.DefaultCommands.AddCommand(EchoCommand)
    commandr.DefaultCommands.AddCommand(DebugCommand)
    commandr.DefaultCommands.AddCommand(LinesCommand)
    commandr.DefaultCommands.AddCommand(AdminLevelCommand)
    commandr.DefaultCommands.AddCommand(&commandr.Command{Use: "cls", Exec: clsCmd, Short: "cls", ExecLevel: commandr.All})

    if err := console.Start(); err != nil {
        log.Fatal(err)
    }
}

func clsCmd(client commandr.Client, args *commandr.CommandArgs) (err error) {

    console.ClsReplContent()
    return
}

func spawnSsh() {

    svc, keys, err := bubblex.CreatesSshService(2222)
    if err != nil {
        loge.Info("Unable to launch ssh server: %v\n", err)
        os.Exit(1)
    }

    svc.RegisterUser("alexj_a", commandr.Admin, keys, nil)
    svc.RegisterUser("alexj_sa", commandr.SuperAdmin, keys, nil)
    svc.RegisterUser("alexj", commandr.User, keys, nil)

    svc.Spawn()
    //utilx.LoopForever(nil)
}

package main

import (
    "bubblex"
    "github.com/alexj212/gox/commandr"
    "log"
    "os"

    "github.com/potakhov/loge"
)

var console *bubblex.TuiClient

func main() {
    cleanup := bubblex.LogeInit()
    defer cleanup()

    go func() {
        spawnSsh()
    }()

    console = bubblex.NewTuiClient()

    for _, cmd := range commandr.DefaultCommands.Commands() {
        if cmd.Use == "cls" || cmd.Use == "exit" {
            commandr.DefaultCommands.RemoveCommand(cmd)
        }
    }

    commandr.DefaultCommands.AddCommand(ExitCommand)
    commandr.DefaultCommands.AddCommand(EchoCommand)
    commandr.DefaultCommands.AddCommand(DebugCommand)
    commandr.DefaultCommands.AddCommand(TldrCmd)
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

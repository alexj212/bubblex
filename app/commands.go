package main

import (
    "github.com/alexj212/gox/commandr"
    "github.com/fatih/color"
    "strings"
)

var ExitCommand = &commandr.Command{Use: "exit", Exec: exitCmd, Short: "exit the session", ExecLevel: commandr.All}

var EchoCommand = &commandr.Command{Use: "echo", Exec: echoCmd, Short: "echo input", ExecLevel: commandr.All}

var DebugCommand = &commandr.Command{Use: "debug", Exec: debugCmd, Short: "debug", ExecLevel: commandr.All}

var LinesCommand = &commandr.Command{Use: "lines", Exec: linesCmd, Short: "lines", ExecLevel: commandr.All}

var AdminLevelCommand = &commandr.Command{Use: "admintest", Exec: adminLevelCmd, Short: "admintest", ExecLevel: commandr.Admin}

func debugCmd(client commandr.Client, args *commandr.CommandArgs) (err error) {
    client.Write([]byte(color.GreenString("args.CmdLine: %v\n", args.CmdLine)))
    client.Write([]byte(color.GreenString("args.Args: %v\n", strings.Join(args.Args, " | "))))
    client.Write([]byte(color.GreenString("args.PealOff: %v\n", args.PealOff(1))))
    client.Write([]byte(color.GreenString("args.Debug: %v\n", args.Debug())))
    return
} // debugCmd

func echoCmd(client commandr.Client, args *commandr.CommandArgs) (err error) {
    //text := args.PealOff(0)
    client.Write([]byte(color.GreenString("%v\n", args.PealOff(0))))
    return
} // echoCmd

func exitCmd(client commandr.Client, args *commandr.CommandArgs) (err error) {
    client.Write([]byte(color.GreenString("Bye bye 👋\n")))
    client.Close()
    return
} // exitCmd

func adminLevelCmd(client commandr.Client, args *commandr.CommandArgs) (err error) {
    //text := args.PealOff(0)
    client.Write([]byte(color.GreenString("admintest\n")))

    return
} // adminLevelCmd

func linesCmd(client commandr.Client, args *commandr.CommandArgs) (err error) {

    cnt := args.FlagSet.Int("cnt", 5, "number of lines to print")
    err = args.Parse()

    if err != nil {
        client.Write([]byte(color.GreenString("lines err: %v\n", err)))
        return
    }

    client.Write([]byte(color.GreenString("lines cnt: %v\n", *cnt)))
    client.Write([]byte(color.GreenString("lines invoked\n")))

    for i := 0; i < *cnt; i++ {
        client.Write([]byte(color.GreenString("line[%d]\n", i)))
    }
    return
} // linesCmd

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/alexj212/bubblex"
	"github.com/alexj212/gox/commandr"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	//"github.com/charmbracelet/bubbletea/termbox"
)

var console *bubblex.Ui
var repl *bubblex.Repl
var logViewer *bubblex.LogViewer
var program *tea.Program

func main() {
	replCmds := &commandr.Command{ExecLevel: commandr.All}
	replCmds.AddCommand(&commandr.Command{Use: "cls",
		Exec: func(client io.Writer, cmd *commandr.Command, args *commandr.CommandArgs) error {
			repl.ClsReplContent()
			return nil
		},
		Short:     "cls",
		ExecLevel: commandr.All})

	replIh := func(s string) {
		var b bytes.Buffer
		out := bufio.NewWriter(&b)

		parsed, err := commandr.NewCommandArgs(s, out)
		if err != nil {
			return
		}

		buf := new(bytes.Buffer)

		execErr := replCmds.Execute(buf, parsed)
		if execErr != nil {
			repl.AddReplContent(color.GreenString(fmt.Sprintf("> %s\n", s)))
			repl.AddReplContent(color.RedString(fmt.Sprintf("%s\n", execErr)))
			return
		}

		repl.AddReplContent(color.GreenString(fmt.Sprintf("> %s\n", s)))
		repl.AddReplContent(color.GreenString(fmt.Sprintf("%s\n", buf.String())))
	}

	logCmds := &commandr.Command{ExecLevel: commandr.All}
	logCmds.AddCommand(&commandr.Command{Use: "cls",
		Exec: func(client io.Writer, cmd *commandr.Command, args *commandr.CommandArgs) error {
			repl.ClsReplContent()
			return nil
		},
		Short:     "cls",
		ExecLevel: commandr.All})

	logIh := func(s string) {
		var b bytes.Buffer
		out := bufio.NewWriter(&b)

		parsed, err := commandr.NewCommandArgs(s, out)
		if err != nil {
			return
		}

		execErr := logCmds.Execute(out, parsed)
		if execErr != nil {
			logViewer.AppendEvent(program, 1, fmt.Sprintf("logIh: %s", execErr))
			return
		}
		out.Flush()
		result := string(out.AvailableBuffer())

		logViewer.AppendEvent(program, 1, s+"AA / AAA"+result)
	}

	repl = bubblex.NewRepl(replIh)
	logViewer = bubblex.NewLogViewer(logIh)

	console = bubblex.NewUi(logViewer, repl)
	program = tea.NewProgram(console)

	defer func() {
		//fmt.Print("\033[H\033[2J")

	}()

	go func() {
		cnt := 0
		for {

			logViewer.AppendEvent(program, cnt%5, fmt.Sprintf("Test %d", cnt))
			cnt++
			time.Sleep(5 * time.Second)
		}
	}()

	_, err := program.Run()
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}

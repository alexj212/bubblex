package bubblex

import (
    "bufio"
    "bytes"
    "fmt"
    "github.com/alexj212/gox/commandr"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/potakhov/loge"
)

type Ui struct {
    u         *TuiClient
    r         *Repl
    logViewer *LogViewer

    activeScreen tea.Model
}

func NewUi(u *TuiClient) *Ui { // tea.Model

    i := &Ui{
        u:         u,
        logViewer: NewLogViewer(),
        r:         NewRepl2(),
    }
    i.activeScreen = i.logViewer
    i.r.InputHandler = func(line string) {

        u.AddReplContent(fmt.Sprintf("> %s\n", line))

        var b bytes.Buffer
        s := bufio.NewWriter(&b)

        parsed, err := commandr.NewCommandArgs(line, s)
        if err != nil {
            s.Flush()
            u.AddReplContent(fmt.Sprintf("%s\n", b.String()))
            u.AddReplContent(fmt.Sprintf("%v\n", err))
            return
        }

        execErr := commandr.DefaultCommands.Execute(i.u, parsed)

        if execErr != nil {
            s.Flush()
            u.AddReplContent(fmt.Sprintf("%s\n", b.String()))
            u.AddReplContent(fmt.Sprintf("%v\n", err))
            return
        }

        s.Flush()
        u.AddReplContent(fmt.Sprintf("%s\n", string(b.Bytes())))

        if u.user != nil {
            u.user.history = append(u.user.history, line)
        }

        return
    }
    return i
}

func (m *Ui) Init() tea.Cmd {
    return nil
}

func (m *Ui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var (
        cmd  tea.Cmd
        cmds []tea.Cmd
    )

    cmds = append(cmds, cmd)

    switch msg := msg.(type) {

    case tea.KeyMsg:

        switch msg.Type {
        case tea.KeyF1:

            if m.activeScreen == m.r {
                loge.Info("switching active screen to LogViewer")
                cmds = append(cmds, tea.EnterAltScreen)
                m.activeScreen = m.logViewer
            }

            break
        case tea.KeyF2:
            if m.activeScreen == m.logViewer {
                loge.Info("switching active screen to repl")
                cmds = append(cmds, tea.ExitAltScreen)
                m.activeScreen = m.r
            }

            break
        case tea.KeyCtrlC, tea.KeyEscape:
            cmds = append(cmds, tea.Quit)
            break
        }
        break

    case tea.WindowSizeMsg:
        m.r.Update(msg)
        m.logViewer.Update(msg)
        break
    }

    m.activeScreen.Update(msg)

    return m, tea.Batch(cmds...)
}

func (m *Ui) View() string {
    return m.activeScreen.View()
}

package bubblex

import (
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/knipferrc/teacup/statusbar"
)

func NewRepl2() *Repl { // tea.Model
    ti := NewTextInput()
    ti.Placeholder = "Enter Command"
    ti.Focus()
    //ti.CharLimit = 156
    //ti.Width = 20

    r := &Repl{}
    sb := statusbar.New(
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
            Background: lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F25D94"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
            Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
            Background: lipgloss.AdaptiveColor{Light: "#A550DF", Dark: "#A550DF"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
            Background: lipgloss.AdaptiveColor{Light: "#6124DF", Dark: "#6124DF"},
        },
    )
    r.statusbar = &sb
    r.textInput = ti
    return r
}

type Repl struct {
    content string

    textInput    TextInput
    totalWidth   int
    totalHeight  int
    statusbar    *statusbar.Bubble
    InputHandler func(string)
}

type cls struct {
}

type setContent struct {
    content string
}
type addLine struct {
    line string
}

func (m *Repl) Init() tea.Cmd {
    return nil
}

func (m *Repl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var (
        cmd  tea.Cmd
        cmds []tea.Cmd
    )

    m.textInput, cmd = m.textInput.Update(msg)
    cmds = append(cmds, cmd)

    switch msg := msg.(type) {
    case cls:
        m.content = ""
        break
    case addLine:
        m.content = m.content + "\n" + msg.line
        break
    case setContent:
        m.content = msg.content
        break

    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyHome:
            m.textInput.CursorStart()

        case tea.KeyEnter:
            input := m.textInput.Value()
            m.textInput.SetValue("")
            if m.InputHandler != nil {
                m.InputHandler(input)
            }

        case tea.KeyEnd:
            m.textInput.CursorEnd()

        case tea.KeyEscape, tea.KeyCtrlC, tea.KeyCtrlQ:
            return m, tea.Quit
        }
    case tea.WindowSizeMsg:
        m.totalWidth = msg.Width - 2
        m.totalHeight = msg.Height
        m.statusbar.SetSize(msg.Width - 2)

        headerHeight := lipgloss.Height(m.headerView())
        footerHeight := lipgloss.Height(m.footerView())
        verticalMarginHeight := headerHeight + footerHeight
        verticalMarginHeight = verticalMarginHeight + 5
        m.statusbar.SetContent("Command Shell", "", "", "Active")
    }

    return m, tea.Batch(cmds...)
}

func (m *Repl) View() string {
    /*
       return lipgloss.NewStyle().Height(m.totalHeight-2).MaxHeight(m.totalHeight).Margin(0).Width(m.totalWidth).Border(lipgloss.ThickBorder(), true).Render(
           lipgloss.JoinVertical(lipgloss.Left,
               m.headerView(),
               m.content,
               m.footerView(),
           ),
       )
    */

    return lipgloss.JoinVertical(lipgloss.Left,
        m.headerView(),
        m.content,
        m.footerView(),
    )
}

func (m *Repl) headerView() string {
    //title := titleStyle.Render("Command Execution")
    //line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
    //return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
    return ""
}

func (m *Repl) footerView() string {
    //info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
    //line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
    //line := strings.Repeat("─", max(0, m.viewport.Width))
    // line,

    //footer := lipgloss.JoinHorizontal(lipgloss.Center, info)

    return lipgloss.JoinVertical(lipgloss.Left,
        //footer,
        m.statusbar.View(),
        styleSubtle.Render("(esc to quit)"),
        styleSubtle.Render(m.textInput.View()))
}

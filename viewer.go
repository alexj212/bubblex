package bubblex

import (
    "fmt"
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/evertras/bubble-table/table"
    "github.com/knipferrc/teacup/statusbar"
    "github.com/potakhov/loge"
)

const (
    columnKeyTime    = "time"
    columnKeyLevel   = "level"
    columnKeyMessage = "message"
    // This is not a visible column, but is used to attach useful reference data
    // to the row itself for easier retrieval
    columnKeyData = "data"
)

var rows []table.Row

func init() {
    rows = make([]table.Row, 0)
}

type LogViewer struct {
    logTable    table.Model
    textInput   TextInput
    totalWidth  int
    totalHeight int
    statusbar   *statusbar.Bubble
}

func NewLogViewer() *LogViewer { // tea.Model

    ti := NewTextInput()
    ti.Placeholder = "Enter Command"
    ti.Focus()
    //ti.CharLimit = 156
    //textti.Width = 20

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

    m := &LogViewer{
        textInput: ti,
        statusbar: &sb,
        logTable: table.New([]table.Column{
            table.NewColumn(columnKeyTime, "Time", 30).WithStyle(lipgloss.NewStyle().Align(lipgloss.Left)),
            table.NewColumn(columnKeyLevel, "Level", 10).WithStyle(lipgloss.NewStyle().Align(lipgloss.Left)),
            table.NewFlexColumn(columnKeyMessage, ":( %", 1).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#c88"))).WithStyle(lipgloss.NewStyle().Align(lipgloss.Left)),
        }).WithRows(rows).
            Border(emptyBorder).
            WithBaseStyle(styleBase).
            WithFooterVisibility(false).
            //WithNoPagination().
            //WithPageSize(6).
            //SortByDesc(columnKeyConversations).
            Focused(true),
    }

    return m
}

func (m *LogViewer) Init() tea.Cmd {
    return tea.Batch(tea.EnterAltScreen)
}

func (m *LogViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var (
        cmd  tea.Cmd
        cmds []tea.Cmd
    )

    m.logTable, cmd = m.logTable.Update(msg)
    cmds = append(cmds, cmd)

    m.textInput, cmd = m.textInput.Update(msg)
    cmds = append(cmds, cmd)

    switch msg := msg.(type) {
    case []table.Row:
        m.logTable = m.logTable.WithRows(msg)
        m.recalculateTable()

    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyHome:
            m.textInput.CursorStart()

        case tea.KeyEnd:
            m.textInput.CursorEnd()
        case tea.KeyEnter:
            loge.Info("Entered cmd: %v", m.textInput.Value())
            m.textInput.SetValue("")

        case tea.KeyEscape, tea.KeyCtrlC, tea.KeyCtrlQ:
            cmds = append(cmds, tea.Quit)
        }

    case tea.WindowSizeMsg:
        m.totalWidth = msg.Width - 2
        m.totalHeight = msg.Height
        m.statusbar.SetSize(msg.Width - 2)
        m.recalculateTable()

    }

    return m, tea.Batch(cmds...)
}
func (m *LogViewer) recalculateTable() {
    m.logTable = m.logTable.WithMaxTotalWidth(m.totalWidth).WithTargetWidth(m.totalWidth).WithPageSize(m.totalHeight - 16) // .PageLast()
    m.statusbar.SetContent("Log Viewer", "", fmt.Sprintf("%d/%d", m.logTable.CurrentPage(), m.logTable.MaxPages()), "Active")
}

func (m *LogViewer) View() string {

    return lipgloss.NewStyle().Height(m.totalHeight-2).MaxHeight(m.totalHeight).Margin(0).Width(m.totalWidth).Border(lipgloss.ThickBorder(), true).Render(
        lipgloss.JoinVertical(lipgloss.Left,
            //lipgloss.NewStyle().MaxHeight(2).Height(2).Width(m.totalWidth+4).Render(
            //    lipgloss.JoinVertical(lipgloss.Left,
            //        styleSubtle.Render("Press q or ctrl+c to quit"),
            //    ),
            //),

            lipgloss.NewStyle().MaxHeight(m.totalHeight-10).Height(m.totalHeight-10).MaxWidth(m.totalWidth).Render(m.logTable.View()),

            lipgloss.NewStyle().MaxHeight(5).Height(5).Width(m.totalWidth+4).Render(
                lipgloss.JoinVertical(lipgloss.Left,
                    m.statusbar.View(),
                    styleSubtle.Render("(esc to quit)"),
                    styleSubtle.Render(m.textInput.View()),

                ),
            ),
        ),
    )
}

func Color(l uint32) string {
    switch l {
    case loge.LogLevelTrace:
        return "#44f"
    case loge.LogLevelDebug:
        return "#44f"
    case loge.LogLevelInfo:
        return "#fa0"
    case loge.LogLevelWarning:
        return "#f64"
    case loge.LogLevelError:
        return "#ff0"
    default:
        return "#44f"
    }
}
func Name(l uint32) string {
    switch l {
    case loge.LogLevelTrace:
        return "trace"
    case loge.LogLevelDebug:
        return "debug"
    case loge.LogLevelInfo:
        return "info"
    case loge.LogLevelWarning:
        return "warn"
    case loge.LogLevelError:
        return "error"
    default:
        return "unknown"
    }
}

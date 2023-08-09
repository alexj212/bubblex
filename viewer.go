package bubblex

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/knipferrc/teacup/statusbar"
)

const (
	columnKeyTime    = "time"
	columnKeyLevel   = "level"
	columnKeyMessage = "message"
)

type LogViewer struct {
	logTable         table.Model
	textInput        TextInput
	statusbar        statusbar.Model
	rows             []table.Row
	totalWidth       int
	totalHeight      int
	onCommandEntered func(string)
}

var _ tea.Model = (*LogViewer)(nil)

type responseMsg struct{}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}
func NewLogViewer(onCommandEntered func(string)) *LogViewer { // tea.Model

	m := &LogViewer{
		onCommandEntered: onCommandEntered,
	}
	m.statusbar = statusbar.New(
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

	m.textInput = NewTextInput()
	m.textInput.Placeholder = "Enter Command"
	m.textInput.Focus()

	m.rows = make([]table.Row, 0)

	columns := []table.Column{
		table.NewColumn(columnKeyTime, "Time", 30).
			WithStyle(lipgloss.NewStyle().Align(lipgloss.Left)),
		table.NewColumn(columnKeyLevel, "Level", 10).
			WithStyle(lipgloss.NewStyle().
				Align(lipgloss.Left)),
		table.NewFlexColumn(columnKeyMessage, "", 1).
			WithStyle(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#c88"))).
			WithStyle(lipgloss.NewStyle().
				Align(lipgloss.Left)),
	}

	m.logTable = table.New(columns).
		WithRows(m.rows).
		Border(emptyBorder).
		WithBaseStyle(styleBase).
		WithFooterVisibility(false).
		//WithNoPagination().
		//WithPageSize(6).
		//SortByDesc(columnKeyConversations).
		Focused(true)

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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width - 2
		m.totalHeight = msg.Height
		m.statusbar.SetSize(msg.Width - 2)
		m.recalculateTable()

	case table.Row:
		m.rows = append(m.rows, msg)
		m.recalculateTable()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyHome:
			m.textInput.CursorStart()

		case tea.KeyEnd:
			m.textInput.CursorEnd()
		case tea.KeyEnter:
			if m.onCommandEntered != nil {
				m.onCommandEntered(m.textInput.Value())
			}
			m.statusbar.SetContent(m.textInput.Value(), "", "", "")
			m.textInput.SetValue("")

		case tea.KeyEscape, tea.KeyCtrlC, tea.KeyCtrlQ:
			return m, tea.Quit
		}
	}
	m.logTable = m.logTable.WithRows(m.rows)

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.logTable, cmd = m.logTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *LogViewer) recalculateTable() {
	m.logTable = m.logTable.WithMaxTotalWidth(m.totalWidth).WithTargetWidth(m.totalWidth).WithPageSize(m.totalHeight - 16) // .PageLast()
	m.statusbar.SetContent("Log Viewer", "", fmt.Sprintf("%d/%d", m.logTable.CurrentPage(), m.logTable.MaxPages()), "Active")
}

func (m *LogViewer) View() string {

	return lipgloss.NewStyle().
		Height(m.totalHeight-2).
		MaxHeight(m.totalHeight).
		Margin(0).
		Width(m.totalWidth).
		Border(lipgloss.ThickBorder(), true).
		Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().
					MaxHeight(m.totalHeight-10).
					Height(m.totalHeight-10).
					MaxWidth(m.totalWidth).
					Render(m.logTable.View()),

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

func Color(l int) string {
	switch l {
	case 1:
		return "#44f"
	case 2:
		return "#44f"
	case 3:
		return "#fa0"
	case 4:
		return "#f64"
	case 5:
		return "#ff0"
	default:
		return "#44f"
	}
}
func Name(l int) string {
	switch l {
	case 1:
		return "trace"
	case 2:
		return "debug"
	case 3:
		return "info"
	case 4:
		return "warn"
	case 5:
		return "error"
	default:
		return fmt.Sprintf("%d", l)
	}
}

func (m *LogViewer) AppendEvent(t *tea.Program, Level int, Message string) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("15:04:05.000 MST")

	row := table.NewRow(
		table.RowData{
			columnKeyTime:  formattedTime,
			columnKeyLevel: table.NewStyledCell(Name(Level), lipgloss.NewStyle().Foreground(lipgloss.Color(Color(Level)))),
			//columnKeyLevel:   Name(Level),
			columnKeyMessage: Message,
		})

	//m.Update(m.rows)
	go t.Send(row)
}

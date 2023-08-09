package bubblex

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/statusbar"
)

type Repl struct {
	content      string
	title        string
	textInput    TextInput
	totalWidth   int
	totalHeight  int
	statusbar    statusbar.Model
	inputHandler func(string)
}

var _ tea.Model = (*Repl)(nil)

func NewRepl(title string, inputHandler func(string)) *Repl { // tea.Model
	r := &Repl{inputHandler: inputHandler}
	r.statusbar = statusbar.New(
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
	r.statusbar.SetContent(title, "", "", "Active")
	r.title = title
	r.textInput = NewTextInput()
	r.textInput.Placeholder = "Enter Command"
	r.textInput.Focus()
	return r
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
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyHome:
			m.textInput.CursorStart()

		case tea.KeyEnter:
			input := m.textInput.Value()
			m.textInput.SetValue("")
			if m.inputHandler != nil {
				m.inputHandler(input)
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
		m.statusbar.SetContent(m.title, "", "", "Active")
		// headerHeight := lipgloss.Height(m.headerView())
		// footerHeight := lipgloss.Height(m.footerView())
		// verticalMarginHeight := headerHeight + footerHeight
		// verticalMarginHeight = verticalMarginHeight + 5
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

func (s *Repl) AddReplContent(line string) {
	s.content = s.content + line + "\n"
}
func (s *Repl) ClsReplContent() {
	s.content = ""
}
func (s *Repl) Printf(format string, a ...any) {
	if s == nil {
		return
	}
	line := fmt.Sprintf(format, a...)
	s.content = s.content + line + "\n"
}

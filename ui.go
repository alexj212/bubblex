package bubblex

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Ui struct {
	activeScreen int
	screens      []tea.Model
}
var _ tea.Model = (*Ui)(nil)
func NewUi(screens ...tea.Model) *Ui { // tea.Model

	i := &Ui{}
	i.screens = make([]tea.Model, 0)
	i.screens = append(i.screens, screens...)

	i.activeScreen = 0
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

			if m.activeScreen != 0 {
				m.activeScreen = 0
			}

		case tea.KeyF2:
			if m.activeScreen != 1 {
				m.activeScreen = 1
			}

		case tea.KeyCtrlLeft:
			m.activeScreen++
			m.activeScreen = m.activeScreen % len(m.screens)

		case tea.KeyCtrlRight:
			m.activeScreen--
			if m.activeScreen < 0 {
				m.activeScreen = len(m.screens) - 1
			}

		case tea.KeyCtrlC, tea.KeyEscape:
			cmds = append(cmds, tea.Quit)
		}

	case tea.WindowSizeMsg:
		for _, screen := range m.screens {
			screen.Update(msg)
		}
	}

	m.screens[m.activeScreen%len(m.screens)].Update(msg)
	return m, tea.Batch(cmds...)
}

func (m *Ui) View() string {
	v := m.screens[m.activeScreen%len(m.screens)].View()
	return v
}

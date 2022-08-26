package bubblex

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/evertras/bubble-table/table"
)

var (
    styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))

    styleBase = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#a7a")).
        BorderForeground(lipgloss.Color("#a38")).
        Align(lipgloss.Right)

    emptyBorder = table.Border{
        Top:    "",
        Left:   "",
        Right:  "",
        Bottom: "",

        TopRight:    "",
        TopLeft:     "",
        BottomRight: "",
        BottomLeft:  "",

        TopJunction:    "",
        LeftJunction:   "",
        RightJunction:  "",
        BottomJunction: "",
        InnerJunction:  "",

        InnerDivider: "",
    }
    customBorder = table.Border{
        Top:    "─",
        Left:   "│",
        Right:  "│",
        Bottom: "─",

        TopRight:    "╮",
        TopLeft:     "╭",
        BottomRight: "╯",
        BottomLeft:  "╰",

        TopJunction:    "┬",
        LeftJunction:   "├",
        RightJunction:  "┤",
        BottomJunction: "┴",
        InnerJunction:  "┼",

        InnerDivider: "│",
    }
)

var clients = make(map[*tea.Program]struct{})

func NotifyAll(r any) {
    //fmt.Printf("NotifyAll rows: %d to clients: %d\n", len(r), len(clients))
    for v, _ := range clients {
        v.Send(r)
    }
}

func UnregisterClient(m *tea.Program) {
    if m != nil {
        delete(clients, m)
    }
}

func RegisterClient(m *tea.Program) {
    if m != nil {
        clients[m] = struct{}{}
    }
    //fmt.Printf("RegisterClient clients: %d\n", len(clients))
}

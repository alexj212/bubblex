package bubblex

import (
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

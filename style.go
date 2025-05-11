package main

import "github.com/charmbracelet/lipgloss"

var (
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("240")).
				BorderBottom(true)

	SelectedRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#262626")).
				Background(lipgloss.Color("#ebbd34")).
				Bold(false)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Border(lipgloss.Border{Right: string('\ue0b0')}, false, true, false, false).
			BorderForeground(lipgloss.Color("#FF5F87")).
			BorderBackground(lipgloss.Color("#353533")).
			Padding(0, 1)

	statusText = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"}).
			PaddingRight(1).
			PaddingLeft(1)

	stateStyle = lipgloss.NewStyle().
			Inherit(statusStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF8531")).
			BorderForeground(lipgloss.Color("#FF8531")).
			BorderBackground(lipgloss.Color("#FF5F87")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Inherit(statusStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF303E")).
			BorderForeground(lipgloss.Color("#FF303E")).
			Padding(0, 1)

	topFileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ebbd34")).
			Bold(true)
)

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
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	stateStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#ff8531")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#ff303e")).
			Padding(0, 1).
			MarginRight(1)

	topFileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ebbd34")).
			Bold(true)
)

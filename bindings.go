package main

import (
	"github.com/charmbracelet/bubbles/key"
)

// bindingKey defines a custom type representing a keyboard's pressed key.
type bindingKey string

func (bk bindingKey) String() string {
	return string(bk)
}

const (
	backspace      bindingKey = "backspace"
	quit           bindingKey = "q"
	cancel         bindingKey = "ctrl+c"
	enter          bindingKey = "enter"
	sort           bindingKey = "s"
	explore        bindingKey = "e"
	sortTotalCap   bindingKey = "alt+t"
	sortTotalUsed  bindingKey = "alt+u"
	sortTotalFree  bindingKey = "alt+f"
	sortTotalUsedP bindingKey = "alt+g"
	toggleTopFiles bindingKey = "ctrl+q"
)

var navigateKeyMap = [][]key.Binding{
	{
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp(
				bindKeyStyle.Render("↑/k"),
				helpDescStyle.Render(" - up"),
			),
		),
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp(
				bindKeyStyle.Render("↓/j"),
				helpDescStyle.Render(" - down"),
			),
		),
		key.NewBinding(
			key.WithKeys("b", "pgup"),
			key.WithHelp(
				bindKeyStyle.Render("b/pgup"),
				helpDescStyle.Render(" - page up"),
			),
		),
		key.NewBinding(
			key.WithKeys("f", "pgdown", " "),
			key.WithHelp(
				bindKeyStyle.Render("f/pgdn"),
				helpDescStyle.Render(" - page down"),
			),
		),
	},
	{
		key.NewBinding(
			key.WithKeys("u", "ctrl+u"),
			key.WithHelp(
				bindKeyStyle.Render("u"),
				helpDescStyle.Render(" - page up"),
			),
		),
		key.NewBinding(
			key.WithKeys("d", "ctrl+d"),
			key.WithHelp(
				bindKeyStyle.Render("d"),
				helpDescStyle.Render(" - page down"),
			),
		),
		key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp(
				bindKeyStyle.Render("g/home"),
				helpDescStyle.Render(" - go to start"),
			),
		),
		key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp(
				bindKeyStyle.Render("G/end"),
				helpDescStyle.Render(" - go to end"),
			),
		),
	},
}

var drivesKeyMap = [][]key.Binding{
	{
		key.NewBinding(
			key.WithKeys(
				sortTotalCap.String(),
				sortTotalUsed.String(),
				sortTotalFree.String(),
			),
			key.WithHelp(
				bindKeyStyle.Render("alt+(t/f/u/g)"),
				helpDescStyle.Render(" - \U000F04BA sort total/free/used/usage"),
			),
		),
		key.NewBinding(
			key.WithKeys(enter.String()),
			key.WithHelp(
				bindKeyStyle.Render(enter.String()),
				helpDescStyle.Render(" - \U000F17A3 open drive"),
			),
		),
		key.NewBinding(
			key.WithKeys(explore.String()),
			key.WithHelp(
				bindKeyStyle.Render(explore.String()),
				helpDescStyle.Render(" - \uF115 explore drive"),
			),
		),
		key.NewBinding(
			key.WithKeys(quit.String(), cancel.String()),
			key.WithHelp(
				bindKeyStyle.Render(quit.String()+"/"+cancel.String()),
				helpDescStyle.Render(" - \U000F0206 quit"),
			),
		),
	},
}

var dirsKeyMap = [][]key.Binding{
	{
		key.NewBinding(
			key.WithKeys(enter.String()),
			key.WithHelp(
				bindKeyStyle.Render(enter.String()),
				helpDescStyle.Render(" - \U000F17A3 open dir"),
			),
		),
		key.NewBinding(
			key.WithKeys(backspace.String()),
			key.WithHelp(
				bindKeyStyle.Render(backspace.String()),
				helpDescStyle.Render(" - \U000F17A7 back"),
			),
		),
		key.NewBinding(
			key.WithKeys(explore.String()),
			key.WithHelp(
				bindKeyStyle.Render(explore.String()),
				helpDescStyle.Render(" - \uF115 explore dir/file"),
			),
		),
		key.NewBinding(
			key.WithKeys(quit.String(), cancel.String()),
			key.WithHelp(
				bindKeyStyle.Render(quit.String()+"/"+cancel.String()),
				helpDescStyle.Render(" - \U000F0206 quit"),
			),
		),
	},
	{
		key.NewBinding(
			key.WithKeys(toggleTopFiles.String()),
			key.WithHelp(
				bindKeyStyle.Render(toggleTopFiles.String()),
				helpDescStyle.Render(" - \U000F028B toggle top files"),
			),
		),
	},
}

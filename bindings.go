package main

import "github.com/charmbracelet/bubbles/key"

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
	toggleTopFiles bindingKey = "alt+q"
)

var drivesKeyMap = [][]key.Binding{
	{
		key.NewBinding(
			key.WithKeys(
				sortTotalCap.String(),
				sortTotalUsed.String(),
				sortTotalFree.String(),
			),
			key.WithHelp(
				"alt+(t/f/u/g)",
				"sort total/free/used/usage",
			),
		),
		key.NewBinding(
			key.WithKeys(enter.String()),
			key.WithHelp(
				enter.String(),
				"open drive",
			),
		),
		key.NewBinding(
			key.WithKeys(explore.String()),
			key.WithHelp(
				explore.String(),
				"explore drive",
			),
		),
		key.NewBinding(
			key.WithKeys(quit.String(), cancel.String()),
			key.WithHelp(
				quit.String()+"/"+cancel.String(),
				"quit",
			),
		),
	},
}

var dirsKeyMap = [][]key.Binding{
	{
		key.NewBinding(
			key.WithKeys(enter.String()),
			key.WithHelp(
				enter.String(),
				"open dir",
			),
		),
		key.NewBinding(
			key.WithKeys(backspace.String()),
			key.WithHelp(
				backspace.String(),
				"back",
			),
		),
		key.NewBinding(
			key.WithKeys(sort.String()),
			key.WithHelp(
				sort.String(),
				"Sort asc/desc",
			),
		),
		key.NewBinding(
			key.WithKeys(explore.String()),
			key.WithHelp(
				explore.String(),
				"explore dir/file",
			),
		),
	},
	{
		key.NewBinding(
			key.WithKeys(quit.String(), cancel.String()),
			key.WithHelp(
				quit.String()+"/"+cancel.String(),
				"quit",
			),
		),
		key.NewBinding(
			key.WithKeys(toggleTopFiles.String()),
			key.WithHelp(
				toggleTopFiles.String(),
				"show/hide top files",
			),
		),
	},
}

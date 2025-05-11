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
	toggleTopFiles bindingKey = "ctrl+q"
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
				"\U000F04BA sort total/free/used/usage",
			),
		),
		key.NewBinding(
			key.WithKeys(enter.String()),
			key.WithHelp(
				enter.String(),
				"\U000F17A3 open drive",
			),
		),
		key.NewBinding(
			key.WithKeys(explore.String()),
			key.WithHelp(
				explore.String(),
				"\uF115 explore drive",
			),
		),
		key.NewBinding(
			key.WithKeys(quit.String(), cancel.String()),
			key.WithHelp(
				quit.String()+"/"+cancel.String(),
				"\U000F0206 quit",
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
				"\U000F17A3 open dir",
			),
		),
		key.NewBinding(
			key.WithKeys(backspace.String()),
			key.WithHelp(
				backspace.String(),
				"\U000F17A7 back",
			),
		),
		key.NewBinding(
			key.WithKeys(explore.String()),
			key.WithHelp(
				explore.String(),
				"\uF115 explore dir/file",
			),
		),
		key.NewBinding(
			key.WithKeys(quit.String(), cancel.String()),
			key.WithHelp(
				quit.String()+"/"+cancel.String(),
				"\U000F0206 quit",
			),
		),
	},
	{
		key.NewBinding(
			key.WithKeys(toggleTopFiles.String()),
			key.WithHelp(
				toggleTopFiles.String(),
				"\U000F028B toggle top files",
			),
		),
	},
}

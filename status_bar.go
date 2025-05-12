package main

import (
	"math"

	"github.com/charmbracelet/lipgloss"
)

const (
	DefaultBarBGColor = "#353533"
	DynamicWidth      = -1
)

// BarItem represents a single status bar item, including its string content,
// background color, and width. The width parameter is optional, and the default
// width equals the content width.
//
// The width value -1 denotes that the item will take all available screen width
// minus the sum of all other elements' widths. If multiple items have a width
// value of -1, the resulting width will be spread equally between them.
type BarItem struct {
	content string
	bgColor string
	width   int
}

// DefaultBarItem returns a new *BarItem instance with default values for
// background color and width.
func DefaultBarItem(content string) *BarItem {
	return &BarItem{content: content, bgColor: DefaultBarBGColor}
}

// NewBarItem returns a new *BarItem instance based on the provided parameters.
// If the background color is an empty string, a default color will be assigned.
func NewBarItem(content, bgColor string, width int) *BarItem {
	if bgColor == "" {
		bgColor = DefaultBarBGColor
	}

	return &BarItem{
		content: content,
		bgColor: bgColor,
		width:   width,
	}
}

// NewStatusBar builds a new status bar based on the provided list of *BarItem
// instances. The total bar width is defined by the totalWidth parameter and all
// bar items will be fit in that width according to their parameters or evenly
// spread for the available width.
//
// NOTE: This implementation does not guarantee that the manually defined element
// sizes will not exceed the totalWidth value.
func NewStatusBar(items []*BarItem, totalWidth int) string {
	styles := make([]lipgloss.Style, 0, len(items))
	renderItems := make([]string, 0, len(items))
	toMaxWidth := make(map[int]struct{}, len(items))

	for i := range items {
		item := items[i]

		itemStyle := newBarBlock(item.bgColor, i != len(items)-1)

		if item.width > 0 {
			itemStyle = itemStyle.Width(item.width)
		}

		if i+1 < len(items) {
			itemStyle = itemStyle.BorderBackground(
				lipgloss.Color(items[i+1].bgColor),
			)
		}

		widthDiff := lipgloss.Width(itemStyle.Render(item.content))

		if item.width == DynamicWidth {
			toMaxWidth[i] = struct{}{}
			widthDiff = 1
		}

		totalWidth -= widthDiff
		styles = append(styles, itemStyle)
	}

	var maxItemWidth int

	if len(toMaxWidth) > 0 {
		maxItemWidth = int(
			math.Ceil(float64(totalWidth) / float64(len(toMaxWidth))),
		)
	}

	for i := range items {
		style := styles[i]

		if _, ok := toMaxWidth[i]; ok {
			style = style.Width(min(totalWidth, maxItemWidth))

			totalWidth -= style.GetWidth()
		}

		renderItems = append(renderItems, style.Render(items[i].content))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderItems...)
}

func newBarBlock(bgColor string, border bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color(bgColor)).
		Padding(0, 1)

	if border {
		style = style.Border(
			lipgloss.Border{Right: string('\ue0b0')}, false, true, false, false).
			BorderForeground(lipgloss.Color(bgColor))
	}

	return style
}

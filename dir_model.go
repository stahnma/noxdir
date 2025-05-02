package main

import (
	"container/heap"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"runtime"
	"strconv"
	"strings"
)

const (
	topFilesTableHeight = 16
	colWidthRatio       = 0.13
)

type DirModel struct {
	columns       []Column
	lastErr       []error
	dirsTable     *table.Model
	topFilesTable *table.Model
	nav           *Navigation
	showTopFiles  bool
	width         int
	height        int
}

func NewDirModel(nav *Navigation) *DirModel {
	dm := &DirModel{
		columns: []Column{
			{Title: ""},
			{Title: ""},
			{Title: "Name"},
			{Title: "Size"},
			{Title: "Total Dirs"},
			{Title: "Total Files"},
			{Title: "Last Change"},
			{Title: "Parent usage"},
			{Title: ""},
		},
		dirsTable:     buildTable(),
		topFilesTable: buildTable(),
		nav:           nav,
	}

	style := table.DefaultStyles()
	style.Header = TableHeaderStyle
	style.Cell = lipgloss.NewStyle()
	style.Selected = lipgloss.NewStyle()

	dm.topFilesTable.SetStyles(style)
	dm.topFilesTable.SetHeight(topFilesTableHeight)

	return dm
}

func (dm *DirModel) Init() tea.Cmd {
	return nil
}

func (dm *DirModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case updateDirState:
		runtime.GC()
		dm.nav.Entry().CalculateSize()

		dm.updateTableData()
	case scanFinished:
		runtime.GC()
		dm.nav.Entry().CalculateSize()
		dm.updateTableData()

		dm.topFilesTable.SetRows(nil)
		dm.fillTopFiles()
	case tea.WindowSizeMsg:
		dm.updateSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch bindingKey(msg.String()) {
		case explore:
			if dm.nav.State() == Drives {
				return dm, nil
			}

			sr := dm.dirsTable.SelectedRow()
			if len(sr) < 2 {
				return dm, nil
			}

			if err := dm.nav.Explore(sr[1]); err != nil {
				return dm, nil
			}
		case toggleTopFiles:
			dm.showTopFiles = !dm.showTopFiles
			dm.updateSize(dm.width, dm.height)
		}
	}

	if dm.nav.State() == Drives {
		return dm, nil
	}

	t, _ := dm.dirsTable.Update(msg)
	dm.dirsTable = &t

	return dm, tea.Batch(cmd)
}

func (dm *DirModel) View() string {
	h := lipgloss.Height

	keyBindings := dm.dirsTable.Help.FullHelpView(
		append(dm.dirsTable.KeyMap.FullHelp(), dirsKeyMap...),
	)
	summary := dm.dirsSummary()

	dirsTableHeight := dm.height - h(keyBindings) - (h(summary) * 2)
	dm.dirsTable.SetHeight(dirsTableHeight)

	rows := []string{summary}

	if dm.showTopFiles {
		tft := dm.topFilesTable.View()

		dm.dirsTable.SetHeight(dirsTableHeight - h(tft))
		rows = append(rows, dm.dirsTable.View(), tft)
	} else {
		rows = append(rows, dm.dirsTable.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Top, append(rows, summary, keyBindings)...,
	)
}

func (dm *DirModel) updateTableData() {
	if dm.nav.State() == Drives || dm.nav.Entry() == nil || !dm.nav.Entry().IsDir {
		return
	}

	iconWidth := 5
	nameWidth := (dm.width - iconWidth) / 4

	colWidth := int(float64(dm.width-iconWidth-nameWidth) * colWidthRatio)
	progressWidth := dm.width - (colWidth * 5) - iconWidth - nameWidth

	columns := make([]table.Column, len(dm.columns))

	for i, c := range dm.columns {
		columns[i] = table.Column{Title: c.Title, Width: colWidth}
	}

	columns[0].Width = iconWidth
	columns[1].Width = 0
	columns[2].Width = nameWidth
	columns[len(columns)-1].Width = progressWidth

	dm.dirsTable.SetColumns(columns)
	dm.dirsTable.SetCursor(0)

	fillProgress := NewProgressBar(progressWidth, '🟥', ' ')

	var rows []table.Row
	dm.nav.Entry().SortChild(false)

	for _, child := range dm.nav.Entry().Child {
		totalDirs, totalFiles := "", ""

		if child.IsDir {
			totalDirs = strconv.FormatUint(child.TotalDirs, 10)
			totalFiles = strconv.FormatUint(child.TotalFiles, 10)
		}

		parentUsage := float64(child.Size) / float64(dm.nav.ParentSize())

		pgBar := fillProgress.ViewAs(parentUsage)

		name := lipgloss.NewStyle().MaxWidth(nameWidth - 5).Render(child.Name)
		if lipgloss.Width(name) == nameWidth-5 {
			name += "..."
		}

		rows = append(
			rows,
			table.Row{
				entryIcon(child),
				child.Name,
				name,
				fmtSize(child.Size, true),
				totalDirs,
				totalFiles,
				child.ModTime.Format("2006-01-02 15:04"),
				strconv.FormatFloat(parentUsage*100, 'f', 2, 64) + " %",
				pgBar,
			},
		)
	}

	dm.dirsTable.SetRows(rows)
	dm.dirsTable.SetCursor(0)
}

func (dm *DirModel) dirsSummary() string {
	w := lipgloss.Width

	state := "READY"
	if dm.nav.Locked() {
		state = "PENDING"
	}

	statuses := []string{
		statusStyle.Render("PATH"),
		"",
		stateStyle.Render("STATE"),
		statusText.Padding(0, 1, 0, 1).Render(state),
		statusStyle.Render("SIZE"),
		statusBarStyle.PaddingRight(1).Render(fmtSize(dm.nav.Entry().Size, false)),
		statusStyle.Render("DIRS"),
		unitFmt(dm.nav.Entry().Dirs),
		statusStyle.Render("FILES"),
		unitFmt(dm.nav.Entry().Files),
		errorStyle.Render("ERRORS"),
		unitFmt(uint64(len(dm.lastErr))),
	}

	pathValWidth := dm.width

	for i := range statuses {
		pathValWidth -= w(statuses[i])
	}

	statuses[1] = statusText.Width(pathValWidth).Render(dm.nav.Entry().Path)

	return statusBarStyle.
		Margin(1, 0, 1, 0).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, statuses...))
}

func (dm *DirModel) fillTopFiles() {
	iconWidth := 5
	colSize := int(float64(dm.width-iconWidth) * colWidthRatio)
	nameWidth := dm.width - (colSize * 2) - iconWidth

	columns := []table.Column{
		{Title: "", Width: iconWidth},
		{Title: "", Width: 0},
		{Title: "Name", Width: nameWidth},
		{Title: "Size", Width: colSize},
		{Title: "Last Change", Width: colSize},
	}

	dm.topFilesTable.SetColumns(columns)
	dm.topFilesTable.SetCursor(0)

	if topFilesInstance.Len() == 0 || dm.topFilesTable.Rows() != nil {
		return
	}

	rows := make([]table.Row, 15)
	heap.Pop(&topFilesInstance)

	for i := len(rows) - 1; i >= 0; i-- {
		file := heap.Pop(&topFilesInstance).(*Entry)

		path := strings.TrimSuffix(
			strings.TrimPrefix(file.Path, dm.nav.currentDrive.Path),
			file.Name,
		)

		rows[i] = table.Row{
			entryIcon(file),
			file.Path,
			path + topFileStyle.Render(file.Name),
			fmtSize(file.Size, true),
			file.ModTime.Format("2006-01-02 15:04"),
		}
	}

	dm.topFilesTable.SetRows(rows)
	dm.topFilesTable.SetCursor(0)
}

func (dm *DirModel) updateSize(width, height int) {
	dm.width, dm.height = width, height

	dm.dirsTable.SetWidth(width)
	dm.topFilesTable.SetWidth(width)

	dm.updateTableData()
	dm.fillTopFiles()
}

func entryIcon(e *Entry) string {
	icon := "📄"

	if !e.IsDir {
		return icon
	}

	icon = "📂"

	if !e.HasChild() {
		icon = "📁"
	}

	return icon
}

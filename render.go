package main

import (
	"dirsize/drive"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"log"
	"time"
)

type (
	updateDirState struct{}
	scanFinished   struct{}
)

var teaProg *tea.Program

type ViewModel struct {
	lastErr    []error
	driveModel *DriveModel
	dirModel   *DirModel
	sortOrder  SortOrder
	nav        *Navigation
	width      int
}

func NewViewModel(n *Navigation) *ViewModel {
	return &ViewModel{
		lastErr:    make([]error, 0),
		nav:        n,
		driveModel: NewDriveModel(n),
		dirModel:   NewDirModel(n),
	}
}

func (vm *ViewModel) Init() tea.Cmd {
	return tea.Batch(tea.DisableMouse)
}

func (vm *ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		vm.width = msg.Width
	case tea.KeyMsg:
		switch bindingKey(msg.String()) {
		case quit, cancel:
			return vm, tea.Quit
		case enter:
			vm.levelDown()
		case backspace:
			vm.levelUp()
		}
	}

	vm.driveModel.Update(msg)
	vm.dirModel.Update(msg)

	return vm, tea.Batch(cmd)
}

func (vm *ViewModel) levelDown() {
	if !vm.nav.Lock() {
		return
	}

	sr := vm.dirModel.dirsTable.SelectedRow()

	if vm.nav.State() == Drives {
		sr = vm.driveModel.drivesTable.SelectedRow()
	}

	done, errChan := vm.nav.LevelDown(
		sr[1],
		func(e *Entry, s State) {
			vm.dirModel.updateTableData()
		},
	)

	if done == nil {
		vm.nav.Unlock()

		return
	}

	go func() {
		vm.lastErr = []error{}

		ticker := time.NewTicker(time.Millisecond * 500)
		defer func() {
			ticker.Stop()
			vm.nav.Unlock()
		}()

		teaProg.Send(updateDirState{})

		for {
			select {
			case err := <-errChan:
				if err != nil {
					vm.lastErr = append(vm.lastErr, err)
				}
			case <-ticker.C:
				teaProg.Send(updateDirState{})
			case <-done:
				teaProg.Send(scanFinished{})

				return
			}
		}
	}()
}

func (vm *ViewModel) levelUp() {
	if !vm.nav.Lock() {
		return
	}
	defer vm.nav.Unlock()

	vm.nav.LevelUp()

	if vm.nav.State() == Drives {
		vm.driveModel.resetSort()
		vm.driveModel.updateTableData(drive.TotalUsedP, true)

		return
	}

	vm.dirModel.updateTableData()
}

func (vm *ViewModel) View() string {
	if vm.nav.State() == Drives {
		return vm.driveModel.View()
	}

	return vm.dirModel.View()
}

func NewProgressBar(width int, full, empty rune) progress.Model {
	maxCharLen := max(
		lipgloss.Width(string(full)),
		lipgloss.Width(string(empty)),
	)

	return progress.New(
		progress.WithColorProfile(termenv.Ascii),
		progress.WithWidth(width/maxCharLen),
		progress.WithFillCharacters(full, empty),
		progress.WithoutPercentage(),
	)
}

func Render() {
	drivesList, err := drive.NewList()
	if err != nil {
		log.Fatalf("Error running program: %s", err.Error())
	}

	teaProg = tea.NewProgram(
		NewViewModel(NewNavigation(drivesList)), tea.WithAltScreen(),
	)

	if _, err = teaProg.Run(); err != nil {
		log.Fatalf("Error running program: %s", err.Error())
	}
}

func buildTable() *table.Model {
	tbl := table.New(table.WithFocused(true))

	style := table.DefaultStyles()
	style.Header = TableHeaderStyle
	style.Cell = lipgloss.NewStyle()
	style.Selected = SelectedRowStyle

	tbl.SetStyles(style)

	tbl.Help = help.New()

	return &tbl
}

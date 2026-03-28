package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	stateInputName sessionState = iota
	stateSelectDB
)

type Model struct {
	State       sessionState
	TextInput   textinput.Model
	ProjectName string
	Cursor      int
	Choices     []string
	Selected    string
	Quitting    bool
}

func InitialModel() Model {
	return InitialModelWithName("")
}

func InitialModelWithName(name string) Model {
	ti := textinput.New()
	ti.Placeholder = "my-fastapi-app"
	ti.CharLimit = 32
	ti.Width = 20

	state := stateInputName
	if name != "" {
		state = stateSelectDB
	} else {
		ti.Focus()
	}

	return Model{
		State:       state,
		TextInput:   ti,
		ProjectName: name,
		Choices:     []string{"PostgreSQL (SQLAlchemy)", "MongoDB (Beanie)", "SQLite (Development)"},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit

		case "enter":
			if m.State == stateInputName {
				
				m.ProjectName = m.TextInput.Value()
				if m.ProjectName == "" {
					m.ProjectName = "my-fastapi-app" // Default
				}
				m.State = stateSelectDB
				return m, nil
			} else {
				m.Selected = m.Choices[m.Cursor]
				return m, tea.Quit
			}

		case "up", "k":
			if m.State == stateSelectDB && m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.State == stateSelectDB && m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		}
	}

	if m.State == stateInputName {
		m.TextInput, cmd = m.TextInput.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.Quitting {
		return "Exiting...\n"
	}

	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)

	if m.State == stateInputName {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			headerStyle.Render("Step 1: Project Name"),
			"What is your project called?",
			m.TextInput.View(),
		) + "\n\n(press enter to continue)\n"
	}

	// Database Selection View
	s := headerStyle.Render("Step 2: Database Selection") + "\n\n"
	s += fmt.Sprintf("Project: %s\n\nChoose a DB:\n", m.ProjectName)

	for i, choice := range m.Choices {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return s + "\n(j/k to move, enter to select)\n"
}
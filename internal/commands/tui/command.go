package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gitkumi/snowflake/internal/initialize"
	"github.com/spf13/cobra"
)

type step int

const (
	stepProjectName step = iota
	stepAppType
	stepDatabase
	stepFeatures
	stepConfirm
)

type feature struct {
	name        string
	description string
	enabled     bool
}

type model struct {
	step        step
	projectName textinput.Model
	appTypes    []initialize.AppType
	databases   []initialize.Database
	features    []feature
	cursor      int
	config      *initialize.Config
	err         error
	done        bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "acme"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	features := []feature{
		{name: "Git", description: "Initialize git repository", enabled: true},
		{name: "SMTP", description: "Email sending capability", enabled: true},
		{name: "Storage", description: "S3-compatible storage", enabled: true},
		{name: "Auth", description: "Authentication system", enabled: true},
	}

	return model{
		step:        stepProjectName,
		projectName: ti,
		appTypes:    initialize.AllAppTypes,
		databases:   initialize.AllDatabases,
		features:    features,
		config:      &initialize.Config{},
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			switch m.step {
			case stepProjectName:
				if m.projectName.Value() == "" {
					return m, nil
				}
				m.config.Name = m.projectName.Value()
				m.step = stepAppType
				return m, nil

			case stepAppType:
				m.config.AppType = m.appTypes[m.cursor]
				m.cursor = 0
				m.step = stepDatabase
				return m, nil

			case stepDatabase:
				m.config.Database = m.databases[m.cursor]
				m.cursor = 0
				m.step = stepFeatures
				return m, nil

			case stepFeatures:
				m.step = stepConfirm
				m.config.NoSMTP = !m.features[0].enabled
				m.config.NoStorage = !m.features[1].enabled
				m.config.NoAuth = !m.features[2].enabled
				m.config.NoGit = !m.features[3].enabled
				return m, nil

			case stepConfirm:
				m.done = true
				return m, tea.Quit
			}

		case "up", "k":
			if m.step != stepProjectName && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			switch m.step {
			case stepAppType:
				if m.cursor < len(m.appTypes)-1 {
					m.cursor++
				}
			case stepDatabase:
				if m.cursor < len(m.databases)-1 {
					m.cursor++
				}
			case stepFeatures:
				if m.cursor < len(m.features)-1 {
					m.cursor++
				}
			}

		case " ":
			if m.step == stepFeatures {
				m.features[m.cursor].enabled = !m.features[m.cursor].enabled
			}
		}
	}

	if m.step == stepProjectName {
		var cmd tea.Cmd
		m.projectName, cmd = m.projectName.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit", m.err)
	}

	switch m.step {
	case stepProjectName:
		s.WriteString("Enter project name:\n\n")
		s.WriteString(m.projectName.View())

	case stepAppType:
		s.WriteString("Select application type:\n\n")
		for i, appType := range m.appTypes {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s.WriteString(fmt.Sprintf("%s %s\n", cursor, appType))
		}

	case stepDatabase:
		s.WriteString("Select database:\n\n")
		for i, db := range m.databases {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s.WriteString(fmt.Sprintf("%s %s\n", cursor, db))
		}

	case stepFeatures:
		s.WriteString("Select features (space to toggle):\n\n")
		for i, feature := range m.features {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			checked := " "
			if feature.enabled {
				checked = "x"
			}
			s.WriteString(fmt.Sprintf("%s [%s] %-10s %s\n", cursor, checked, feature.name, feature.description))
		}

	case stepConfirm:
		s.WriteString("Configuration Summary:\n\n")
		s.WriteString(fmt.Sprintf("Project Name: %s\n", m.config.Name))
		s.WriteString(fmt.Sprintf("App Type: %s\n", m.config.AppType))
		s.WriteString(fmt.Sprintf("Database: %s\n", m.config.Database))
		s.WriteString("\nFeatures:\n")
		s.WriteString(fmt.Sprintf("- SMTP: %v\n", !m.config.NoSMTP))
		s.WriteString(fmt.Sprintf("- Storage: %v\n", !m.config.NoStorage))
		s.WriteString(fmt.Sprintf("- Auth: %v\n", !m.config.NoAuth))
		s.WriteString(fmt.Sprintf("- Git: %v\n", !m.config.NoGit))
		s.WriteString("\nPress enter to confirm and generate project")
	}

	s.WriteString("\n\nPress q to quit\n")

	return s.String()
}

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project using the TUI",
		Run: func(cmd *cobra.Command, args []string) {
			m := initialModel()
			finalModel, err := tea.NewProgram(m).Run()
			if err != nil {
				fmt.Printf("error running tui: %v\n", err)
				return
			}

			m, ok := finalModel.(model)
			if !ok || !m.done {
				return
			}

			if err := initialize.Initialize(m.config); err != nil {
				fmt.Printf("error creating project: %v\n", err)
			}
		},
	}

	return cmd
}

package picker

import (
	"fmt"

	"github.com/aldrickdev/dbm-sandbox/internal/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultHeight = 24
	defaultWidth  = 20
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	List     list.Model
	choice   string
	output   *string
	quitting bool
}

// NewPicker is Bubble Tea application that implements the RunnableQuestion interface.
// When the application is ran using the Run method, it will provide the user an
// interface where the user needs to select from predefined options to answer a
// Providers question.
func NewPicker(options []string, optionDesc []string, prompt string, output *string) model {
	items := []list.Item{}

	itemDelegate := list.NewDefaultDelegate()
	if optionDesc == nil {
		itemDelegate.ShowDescription = false

		for _, option := range options {
			items = append(items, item{title: option, desc: ""})
		}
	} else {
		for ix, option := range options {
			items = append(items, item{title: option, desc: optionDesc[ix]})
		}
	}

	itemStyles := list.NewDefaultItemStyles()

	itemStyles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color(styles.DatadogColor)).
		Foreground(lipgloss.Color(styles.DatadogColor)).
		Padding(0, 0, 0, 2)
	itemStyles.SelectedDesc = itemStyles.SelectedTitle.Copy()

	itemDelegate.Styles = itemStyles

	l := list.New(items, itemDelegate, defaultWidth, defaultHeight)
	l.Title = prompt
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	l.Styles.Title = styles.ListItemTitle
	l.Styles.PaginationStyle = styles.ListPagination
	l.Styles.HelpStyle = styles.ListHelpStyle

	return model{List: l, output: output}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.List.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
				*m.output = string(i.title)
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m model) View() string {
	question := styles.Question.Render(fmt.Sprint(m.List.Title))

	if m.choice != "" {
		selection := styles.Answer.Render(fmt.Sprint(m.choice))
		return lipgloss.JoinHorizontal(lipgloss.Bottom, question, selection)
	}

	if m.quitting {
		quitText := styles.Quitting.Render("No option selected, 👋 Bye")
		return lipgloss.JoinVertical(lipgloss.Left, question, quitText)
	}
	return "\n" + m.List.View()
}

func (m model) Run() error {
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return err
	}

	return nil
}

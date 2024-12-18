package textInput

import (
	"fmt"

	"github.com/aldrickdev/dbm-sandbox/internal/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	CHARLIMIT = 156
	WIDTH     = 20
)

var QuitTextStyle = lipgloss.NewStyle().
	Margin(1, 0, 2, 4)

type (
	errMsg error
)

type model struct {
	prompt       string
	textInput    textinput.Model
	defaultValue string
	input        string
	output       *string
	quitting     bool
	err          error
}

// NewTextInput returns a Bubble Tea application that implements the
// RunnableQuestion interface. When the application is ran using the Run
// method, it will provide the user an interface where they can provide
// an answer using text.
func NewTextInput(prompt string, placeholder string, output *string) model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = CHARLIMIT
	ti.Width = WIDTH

	return model{
		prompt:       prompt,
		textInput:    ti,
		defaultValue: placeholder,
		input:        "",
		quitting:     false,
		output:       output,
		err:          nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			if m.textInput.Value() == "" {
				*m.output = m.defaultValue
				m.input = m.defaultValue

				return m, tea.Quit
			}

			*m.output = m.textInput.Value()
			m.input = m.textInput.Value()

			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	question := styles.Question.Render(fmt.Sprint(m.prompt))

	if m.input != "" {
		selection := styles.Answer.Render(fmt.Sprint(*m.output))
		return lipgloss.JoinHorizontal(lipgloss.Bottom, question, selection)
	}

	if m.quitting {
		quitText := styles.Quitting.Render("No value provided, 👋 Bye")
		return lipgloss.JoinVertical(lipgloss.Left, question, quitText)
	}

	return styles.Question.Render(fmt.Sprintf("%s\n\n%s\n", m.prompt, m.textInput.View()))
}

func (m model) Run() error {
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return err
	}

	return nil
}

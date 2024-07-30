package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/list"
)

const (
	// Base TUI Colors
	DatadogColor = "#632CA6"
	WhiteColor   = "#FFF"
)

var (
	DatadogColoredText = lipgloss.NewStyle().
		Foreground(lipgloss.Color(DatadogColor))

	WhiteColoredText = lipgloss.NewStyle().
		Foreground(lipgloss.Color(WhiteColor))


	ProjectTitle = DatadogColoredText.Copy().
		MaxWidth(80).
		Margin(0, 2)

	Question = WhiteColoredText.Copy().
		Width(40).
		Margin(1, 4)

	Quitting = WhiteColoredText.Copy().
			Margin(0, 0, 2, 4)

	Answer = lipgloss.NewStyle().
			Foreground(lipgloss.Color(DatadogColor)).
			Margin(1)

	StatusIndicator = WhiteColoredText.Copy().
		Padding(1, 1, 1, 4)

	Error = StatusIndicator.Copy().
		SetString("🔴  ")

	Success = StatusIndicator.Copy().
		SetString("🟢  ")

	ListItemTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(WhiteColor)).
			MarginLeft(2)

	ListPagination = list.DefaultStyles().
			PaginationStyle.PaddingLeft(4)

	ListHelpStyle = list.DefaultStyles().
			HelpStyle.PaddingLeft(4).
			PaddingBottom(1)
)

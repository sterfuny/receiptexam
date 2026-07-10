package receiptexam

import "charm.land/lipgloss/v2"

var(
	// styleReset = "\x1b[0m"
	Right = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))
	Wrong = lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))
	Bg = lipgloss.NewStyle().
		Background(lipgloss.Color("240"))
)

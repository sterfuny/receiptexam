package questions

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type ChoiceSingleQ struct {
	question string
	options  []string
	cursor   int
	answer   string
	done     bool
}

func NewChoiceSingle(q string, opts []string) *ChoiceSingleQ {
	return &ChoiceSingleQ{
		question: q,
		options:  opts,
		cursor:   0,
		done:     false,
	}
}

func (q *ChoiceSingleQ) Render() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s (选择题)\n\n", q.question))

	for i, opt := range q.options {
		cursor := " "
		if i == q.cursor {
			cursor = ">"
		}
		prefix := "○"
		if i == q.cursor {
			prefix = "●"
		}
		sb.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, prefix, opt))
	}

	if q.done {
		sb.WriteString(fmt.Sprintf("\n✅ 已选择: %s", q.answer))
	} else {
		sb.WriteString("\n↑/↓ 移动 • Enter 确认")
	}
	return sb.String()
}

func (q *ChoiceSingleQ) HandleKey(msg tea.KeyMsg) (done bool, quit bool) {
	if q.done {
		return true, false
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return false, true
	case "up", "k":
		if q.cursor > 0 {
			q.cursor--
		}
	case "down", "j":
		if q.cursor < len(q.options)-1 {
			q.cursor++
		}
	case "enter":
		q.answer = q.options[q.cursor]
		q.done = true
		return true, false
	}
	return false, false
}

func (q *ChoiceSingleQ) GetQuestionText() string {
	return q.question
}

func (q *ChoiceSingleQ) GetAnswer() any {
	return q.answer
}

func (q *ChoiceSingleQ) IsDone() bool {
	return q.done
}

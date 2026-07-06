package questions

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type InputSinceQ struct {
	question  string
	textInput textinput.Model
	answer    string
	done      bool
}

func NewInputSince(q string) *InputSinceQ {
	ti := textinput.New()
	ti.Placeholder = "请输入答案"
	ti.SetVirtualCursor(true)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)

	return &InputSinceQ{
		question:  q,
		textInput: ti,
		done:      false,
	}
}

func (q *InputSinceQ) Init() tea.Cmd {
	return textinput.Blink
}

func (q *InputSinceQ) HandleKey(msg tea.KeyMsg) (done bool, quit bool) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return true, true
	case "enter":
		q.answer = q.textInput.Value()
		q.done = true
		return true, false
	}

	updated, _ := q.textInput.Update(msg)
	q.textInput = updated
	return false, false
}

func (q *InputSinceQ) Render() string {
	header := q.question + "\n"
	str := lipgloss.JoinVertical(lipgloss.Top, header, q.textInput.View(), "(esc to quit)")
	if q.done {
		str += "\n"
	}
	return str
}

func (q *InputSinceQ) GetQuestionText() string {
	return q.question
}

func (q *InputSinceQ) GetAnswer() any {
	return q.answer
}

func (q *InputSinceQ) IsDone() bool {
	return q.done
}

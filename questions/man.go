package questions

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Question 枚举
type QuestionType int

const (
	ChoiceSingle QuestionType = iota + 1
	ChoiceMulti
	InputSince
	InputChunk
)

var typeNames = [...]string{
	ChoiceSingle: "ChoiceSingle",
	ChoiceMulti:  "ChoiceMulti",
	InputSince:   "InputSince",
	InputChunk:   "InputChunk",
}

func (t QuestionType) Name() string {
	if t >= 0 && int(t) < len(typeNames) {
		return typeNames[t]
	}
	return "Unknown"
}

func ParseType(name string) QuestionType {
	for t, n := range typeNames {
		if strings.EqualFold(name, n) {
			return QuestionType(t)
		}
	}
	return -1
}

type Question interface {
	// Render 渲染题目（包括当前输入/选择状态）
	Render() string

	// HandleKey 处理键盘输入，返回 (是否完成, 是否退出)
	HandleKey(msg tea.KeyMsg) (done bool, quit bool)

	// GetQuestionText 获取题目文字
	GetQuestionText() string

	// GetAnswer 获取答案（用于最终汇总）
	GetAnswer() any

	// IsDone 是否已答完
	IsDone() bool
}

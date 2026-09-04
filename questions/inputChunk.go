package questions

import (
	"fmt"
	"strings"
	"strconv"	

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// 块类型枚举int
type ChunkKind int
const (
	textChunk ChunkKind = iota
	editChunk
)

// 以框为单位的块
type Chunk struct {
	Kind  ChunkKind
	Text  string
	Input *textinput.Model
}

func (c Chunk) Width() int {
	switch c.Kind {
	case textChunk:
		return lipgloss.Width(c.Text)
	case editChunk:
		if c.Input == nil {
			return 0
		}
		return lipgloss.Width(c.Input.View())
	}
	return 0
}

// 输入框配置
type InputConfig struct {
	Placeholder string
	CharLimit   int
}

// 多输入框填空题
type InputChunkQ struct {
	question   string
	chunks     []Chunk
	inputs     []*textinput.Model
	focusedIdx int
	answers    []string
	done       bool
	width      int
}

// 创建多输入框题目
func NewInputChunkQ(question string, configs []InputConfig, template string) *InputChunkQ {
	inputs := make([]*textinput.Model, len(configs))
	for i, cfg := range configs {
		ti := newInput(cfg)
		inputs[i] = &ti
	}

	if len(inputs) > 0 {
		inputs[0].Focus()
	}

	chunks := buildParagraph(inputs, template)

	return &InputChunkQ{
		question:   question,
		chunks:     chunks,
		inputs:     inputs,
		focusedIdx: 0,
		answers:    make([]string, len(inputs)),
		done:       false,
		width:      80,
	}
}

// 创建单个输入框
func newInput(cfg InputConfig) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = cfg.Placeholder
	ti.Prompt = ""
	ti.SetVirtualCursor(false)
	ti.CharLimit = cfg.CharLimit
	ti.SetWidth(lipgloss.Width(cfg.Placeholder))
	return ti
}


// 构建类型列表
func buildParagraph(inputs []*textinput.Model, template string) []Chunk {
	// 示例模板："请问您的名字是 {0} ，年龄 {1} 岁，喜欢的颜色是 {2} ，来自 {3}。"
	leafovers := template
	var chunks []Chunk

	for len(leafovers) > 0 {
		// 找 _{ 的位置
		start := strings.Index(leafovers, "_{")
		if start == -1 {
			if leafovers != "" {
				chunks = append(chunks, Chunk{Kind: textChunk,Text: leafovers})
			}
			break
		} else if start > 0 {
			chunks = append(chunks, Chunk{Kind: textChunk,Text: leafovers[:start]})
		}

		// 从 _{ 之后找 }_ 的位置
		rest := leafovers[start+2:] // 跳过 "_{"
		end := strings.Index(rest, "}_")
		if end == -1 {
			// 没有闭合标记
			if leafovers != "" {
				chunks = append(chunks, Chunk{Kind: textChunk,Text: leafovers})
			}
			break
		}

		// 标记前的字符
		if start > 0 {
			chunks = append(chunks, Chunk{Kind: textChunk,Text: leafovers[:start]})
		}

		// 数字部分
		numStr := strings.TrimSpace(rest[:end])
		if num, err := strconv.Atoi(numStr); err == nil {
			chunks = append(
				chunks, Chunk{
					Kind: editChunk,
					Input: inputs[num],
				})
		} else {
			// chunks = append(chunks, numStr)
		}
		leafovers = rest[end+2:] // 跳过 "}_"
	}

	return chunks
}

// 设置渲染宽度
func (q *InputChunkQ) SetWidth(w int) {
	q.width = w
}

// 初始化，返回所有输入框的 Blink 命令
func (q *InputChunkQ) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(q.inputs))
	for i := range q.inputs {
		cmds[i] = textinput.Blink
	}
	return tea.Batch(cmds...)
}

// 处理按键，返回 (done, quit)
func (q *InputChunkQ) HandleKey(msg tea.KeyMsg) (done bool, quit bool) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return true, true
	case "tab":
		if len(q.inputs) == 0 {
			return false, false
		}
		q.inputs[q.focusedIdx].Blur()
		q.focusedIdx = (q.focusedIdx + 1) % len(q.inputs)
		q.inputs[q.focusedIdx].Focus()
		return false, false
	case "shift+tab":
		if len(q.inputs) == 0 {
			return false, false
		}
		q.inputs[q.focusedIdx].Blur()
		q.focusedIdx = (q.focusedIdx - 1 + len(q.inputs)) % len(q.inputs)
		q.inputs[q.focusedIdx].Focus()
		return false, false
	case "enter":
		q.answers = make([]string, len(q.inputs))
		for i, inp := range q.inputs {
			q.answers[i] = inp.Value()
		}
		q.done = true
		return true, false
	default:
		if len(q.inputs) > 0 {
			updated, _ := q.inputs[q.focusedIdx].Update(msg)
			*q.inputs[q.focusedIdx] = updated
		}
		return false, false
	}
}

// 渲染题目内容
func (q *InputChunkQ) Render() string {
	maxWidth := q.width
	if maxWidth <= 0 {
		maxWidth = 80
	}

	lines := q.buildLines(maxWidth)
	var rowStrings []string
	for _, line := range lines {
		var parts []string
		for _, ch := range line {
			switch ch.Kind {
			case textChunk:
				parts = append(parts, ch.Text)
			case editChunk:
				parts = append(parts, ch.Input.View())
			}
		}
		rowStrings = append(rowStrings, lipgloss.JoinHorizontal(lipgloss.Top, parts...))
	}

	body := lipgloss.JoinVertical(lipgloss.Left, rowStrings...)

	if q.done {
		body += "\n\n📋 输入结果:\n"
		for i, ans := range q.answers {
			body += fmt.Sprintf("  %d. %s\n", i+1, ans)
		}
		body += "\n"
	} else {
		body += "\n\n💡 提示: Tab 切换输入框 | Shift+Tab 返回 | Enter 提交 | Esc 退出\n"
	}

	return body
}

// 根据最大宽度将 chunks 分行
func (q *InputChunkQ) buildLines(maxWidth int) [][]Chunk {
	var lines [][]Chunk
	var curLine []Chunk
	curWidth := 0

	for _, ch := range q.chunks {
		w := ch.Width()
		if len(curLine) > 0 && curWidth+w > maxWidth {
			lines = append(lines, curLine)
			curLine = []Chunk{ch}
			curWidth = w
		} else {
			curLine = append(curLine, ch)
			curWidth += w
		}
	}
	if len(curLine) > 0 {
		lines = append(lines, curLine)
	}
	return lines
}

func (q *InputChunkQ) GetQuestionText() string {
	return q.question
}

func (q *InputChunkQ) GetAnswer() any {
	return q.answers
}

func (q *InputChunkQ) IsDone() bool {
	return q.done
}

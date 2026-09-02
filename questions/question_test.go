package questions

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func pressRune(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: r, Text: string(r)}
}

func pressKey(k rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: k}
}

func printRender(t *testing.T, tag string, q Question) {
	t.Helper()
	t.Logf("\n===== %s =====\n%s\n", tag, q.Render())
}

func TestNewChoiceSingle(t *testing.T) {

	q := NewChoiceSingle("Go 是编译型语言吗?", []string{"是", "否", "不确定"})
	t.Run("初始状态与渲染", func(t *testing.T) {
		// 初始状态
		if q.cursor != 0 {
			t.Errorf("初始 cursor = %d, 期望 0", q.cursor)
		}
		if q.IsDone() {
			t.Error("初始不应为 done")
		}
		if got := q.GetAnswer(); got != "" {
			t.Errorf("初始 answer = %q, 期望空", got)
		}
		printRender(t, "初始", q)
	})

	t.Run("选择并确认", func(t *testing.T) {
		// 边界:cursor=0 时按 up 不应越界
		done, quit := q.HandleKey(pressKey(tea.KeyUp))
		if done || quit {
			t.Errorf("up 在顶部不应 done/quit, got done=%v quit=%v", done, quit)
		}
		if q.cursor != 0 {
			t.Errorf("up 在顶部后 cursor = %d, 期望仍为 0", q.cursor)
		}
	
		// down down down: 最多到最后一个
		for i := 0; i < 3; i++ {
			if done, quit := q.HandleKey(pressKey(tea.KeyDown)); done || quit {
				t.Fatalf("down 不应 done/quit, got done=%v quit=%v", done, quit)
			}
		}
		if q.cursor != len(q.options)-1 {
			t.Errorf("连续 down 后 cursor = %d, 期望 %d", q.cursor, len(q.options)-1)
		}

		// up 一次
		if _, quit := q.HandleKey(pressKey(tea.KeyUp)); quit {
			t.Fatal("up 不应 quit")
		}
		if q.cursor != 1 {
			t.Errorf("up 后 cursor = %d, 期望 1", q.cursor)
		}
		printRender(t, "移动光标后", q)

		// enter 确认
		done, quit = q.HandleKey(pressKey(tea.KeyEnter))
		if !done || quit {
			t.Errorf("enter 应 done=true quit=false, got done=%v quit=%v", done, quit)
		}
		if !q.IsDone() {
			t.Error("enter 后应 IsDone()")
		}
		if got := q.GetAnswer(); got != "否" {
			t.Errorf("answer = %q, 期望 %q", got, "否")
		}
		printRender(t, "确认后", q)
	})

	t.Run("退出后响应", func(t *testing.T) {
		// done 后按键不再生效
		done, quit := q.HandleKey(pressKey(tea.KeyDown))
		if !done || quit {
			t.Errorf("done 后按键应无效果, got done=%v quit=%v", done, quit)
		}
		if q.cursor != 1 {
			t.Errorf("done 后 cursor 不应变化, 仍为 %d", q.cursor)
		}
	})
}

func TestNewChoiceSingleQuit(t *testing.T) {
	q := NewChoiceSingle("q 退出测试", []string{"a", "b"})
	done, quit := q.HandleKey(pressRune('q'))
	if done || !quit {
		t.Errorf("q 应 done=false quit=true, got done=%v quit=%v", done, quit)
	}
}

func TestNewInputSince(t *testing.T) {
	t.Run("初始状态与渲染", func(t *testing.T) {
		q := NewInputSince("你的名字是?")

		if cmd := q.Init(); cmd == nil {
			t.Error("Init() 应返回非 nil cmd (textinput.Blink)")
		}
		if q.IsDone() {
			t.Error("初始不应为 done")
		}
		if !q.textInput.Focused() {
			t.Error("textinput 应自动聚焦")
		}
		printRender(t, "初始", q)
	})

	t.Run("输入并确认", func(t *testing.T) {
		q := NewInputSince("你的名字是?")

		for _, r := range "hello 你好" {
			if done, quit := q.HandleKey(pressRune(r)); done || quit {
				t.Fatalf("输入字符 %q 不应 done/quit, got done=%v quit=%v", r, done, quit)
			}
		}
		if got := q.textInput.Value(); got != "hello 你好" {
			t.Errorf("输入值 = %q, 期望 %q", got, "hello 你好")
		}
		printRender(t, "输入后", q)

		done, quit := q.HandleKey(pressKey(tea.KeyEnter))
		if !done || quit {
			t.Errorf("enter 应 done=true quit=false, got done=%v quit=%v", done, quit)
		}
		if got := q.GetAnswer(); got != "hello 你好" {
			t.Errorf("answer = %q, 期望 %q", got, "hello 你好")
		}
		if !strings.Contains(q.Render(), "hello 你好") {
			t.Error("done 后 Render 仍应包含输入内容")
		}
	})

	t.Run("esc 退出", func(t *testing.T) {
		q := NewInputSince("你的名字是?")
		done, quit := q.HandleKey(pressKey(tea.KeyEscape))
		if !done || !quit {
			t.Errorf("esc 应 done=true quit=true, got done=%v quit=%v", done, quit)
		}
	})
}

package cli

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// sendKey は tea.KeyMsg をモデルに送り、更新後のモデルを返す
func sendKey(t *testing.T, m tea.Model, msg tea.KeyMsg) tea.Model {
	t.Helper()
	updated, _ := m.Update(msg)
	return updated
}

func typeRunes(t *testing.T, m tea.Model, s string) tea.Model {
	t.Helper()
	return sendKey(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)})
}

func TestSelectModel_ArrowNavigation(t *testing.T) {
	items := []selectItem{{label: "posts", value: "posts"}, {label: "notes", value: "notes"}}
	var m tea.Model = newSelectModel("prompt", items, "manual")

	// 初期カーソルは先頭
	if got := m.(selectModel).cursor; got != 0 {
		t.Fatalf("initial cursor = %d, want 0", got)
	}

	// ↓ で下に移動
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyDown})
	if got := m.(selectModel).cursor; got != 1 {
		t.Fatalf("after down cursor = %d, want 1", got)
	}

	// 手動入力の項目（末尾）まで移動し、それ以上は進まない
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyDown})
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyDown})
	if got, want := m.(selectModel).cursor, len(items); got != want {
		t.Fatalf("cursor = %d, want %d (clamped at manual item)", got, want)
	}

	// ↑ で上に戻る
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyUp})
	if got := m.(selectModel).cursor; got != 1 {
		t.Fatalf("after up cursor = %d, want 1", got)
	}
}

func TestSelectModel_SelectItem(t *testing.T) {
	items := []selectItem{{label: "posts", value: "posts"}, {label: "notes", value: "notes"}}
	var m tea.Model = newSelectModel("prompt", items, "manual")

	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyDown}) // notes を選ぶ
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	res := m.(selectModel)
	if res.result != "notes" {
		t.Errorf("result = %q, want %q", res.result, "notes")
	}
	if res.aborted {
		t.Error("aborted = true, want false")
	}
	if cmd == nil {
		t.Error("expected quit command on selection")
	}
}

func TestSelectModel_ManualInput(t *testing.T) {
	items := []selectItem{{label: "posts", value: "posts"}}
	var m tea.Model = newSelectModel("prompt", items, "manual")

	// 手動入力の項目へ移動して Enter → 入力モードに入る
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyDown})
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if !m.(selectModel).inputting {
		t.Fatal("expected inputting mode after selecting manual item")
	}

	// 入力して確定
	m = typeRunes(t, m, "my-section")
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if got := m.(selectModel).result; got != "my-section" {
		t.Errorf("result = %q, want %q", got, "my-section")
	}
}

func TestSelectModel_Abort(t *testing.T) {
	items := []selectItem{{label: "posts", value: "posts"}}
	var m tea.Model = newSelectModel("prompt", items, "manual")

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if !m.(selectModel).aborted {
		t.Error("aborted = false, want true after ctrl+c")
	}
	if cmd == nil {
		t.Error("expected quit command on abort")
	}
}

// TestInputModel_Backspace は「文字が消せない」バグが解消されていることを確認する
func TestInputModel_Backspace(t *testing.T) {
	var m tea.Model = newInputModel("title")

	m = typeRunes(t, m, "abcd")
	// バックスペースで2文字消す
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyBackspace})
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyBackspace})

	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if got := m.(inputModel).result; got != "ab" {
		t.Errorf("result = %q, want %q (backspace should delete characters)", got, "ab")
	}
}

func TestInputModel_DeleteAll(t *testing.T) {
	var m tea.Model = newInputModel("title")

	m = typeRunes(t, m, "abc")
	for range 3 {
		m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyBackspace})
	}
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if got := m.(inputModel).result; got != "" {
		t.Errorf("result = %q, want empty string after deleting all", got)
	}
}

func TestInputModel_TrimsWhitespace(t *testing.T) {
	var m tea.Model = newInputModel("title")
	m = typeRunes(t, m, "  hello  ")
	m = sendKey(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if got := m.(inputModel).result; got != "hello" {
		t.Errorf("result = %q, want %q", got, "hello")
	}
}

func TestInputModel_Abort(t *testing.T) {
	var m tea.Model = newInputModel("title")
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !m.(inputModel).aborted {
		t.Error("aborted = false, want true after esc")
	}
	if cmd == nil {
		t.Error("expected quit command on abort")
	}
}

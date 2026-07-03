package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ErrAborted はユーザが操作を中断（Ctrl+C など）したときに返される
var ErrAborted = errors.New("operation aborted by user")

var (
	promptStyle = lipgloss.NewStyle().Bold(true)
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	helpStyle   = lipgloss.NewStyle().Faint(true)
)

// このファイルは bubbletea による TUI 実装をまとめている。
//
// bubbletea は The Elm Architecture という設計で、UIを次の3つに分解する:
//   - Model : 「今の状態」を持つ構造体（例: 選択中のカーソル位置）
//   - Update: キー入力などのメッセージを受け、新しい Model を返す（状態遷移）
//   - View  : Model を見て画面に描く文字列を組み立てる
// 入力待ちのループや再描画は bubbletea 側が回すため、こちらは自前で
// ReadString や for ループを書かない。従来の自作TUIをやめてこの方式にしたのは、
// テキスト編集や選択カーソルの管理を実績あるフレームワークに委ねるためである。

// bubbleTeaPrompter は bubbletea を使った対話的な入力を提供する Prompter 実装
type bubbleTeaPrompter struct {
	out io.Writer
}

func newBubbleTeaPrompter(out io.Writer) *bubbleTeaPrompter {
	return &bubbleTeaPrompter{out: out}
}

func (p *bubbleTeaPrompter) SelectSection(sections []string) (string, error) {
	items := make([]selectItem, 0, len(sections))
	for _, s := range sections {
		items = append(items, selectItem{label: s, value: s})
	}
	return p.runSelect(
		"Please select or enter the section name for the new article.",
		items,
		"Enter a new section name",
	)
}

func (p *bubbleTeaPrompter) SelectSlug(candidates []string) (string, error) {
	items := make([]selectItem, 0, len(candidates)+1)
	for _, c := range candidates {
		items = append(items, selectItem{label: c, value: c})
	}
	items = append(items, selectItem{label: "🔄 Regenerate candidates", value: redoCommand})
	return p.runSelect(
		"Select a slug for the new article.",
		items,
		"Enter a custom slug",
	)
}

func (p *bubbleTeaPrompter) InputTitle() (string, error) {
	return p.runInput("Please enter the title of the article.")
}

// runSelect は selectModel を起動し、確定した値を通常の string として返す。
// bubbletea の世界（Model/Update/View）をこの関数の内側に閉じ込めることで、
// 呼び出し元の Prompter からは「文字列を1つ返す関数」にしか見えなくなる。
func (p *bubbleTeaPrompter) runSelect(prompt string, items []selectItem, manualLabel string) (string, error) {
	// Run() はユーザが確定/中断するまで戻らず、終了時点の Model を返す。
	m := newSelectModel(prompt, items, manualLabel)
	final, err := tea.NewProgram(m, tea.WithOutput(p.out)).Run()
	if err != nil {
		return "", err
	}
	res := final.(selectModel)
	if res.aborted {
		return "", ErrAborted
	}
	return res.result, nil
}

func (p *bubbleTeaPrompter) runInput(prompt string) (string, error) {
	m := newInputModel(prompt)
	final, err := tea.NewProgram(m, tea.WithOutput(p.out)).Run()
	if err != nil {
		return "", err
	}
	res := final.(inputModel)
	if res.aborted {
		return "", ErrAborted
	}
	return res.result, nil
}

// selectItem は選択肢1つを表す（表示ラベルと確定時の値）
type selectItem struct {
	label string
	value string
}

// selectModel は上下矢印で選択でき、末尾の項目で手動入力に切り替えられるリストモデル。
//
// 「一覧から選ぶ」と「自由入力する」を1つのモデルに統合しているのは、
// セクション名もSlugも "候補から選ぶ or 自分で打つ" が同じ体験だからである。
// 末尾に手動入力用の項目を1つ足し、そこで Enter されたら入力モードへ遷移する。
type selectModel struct {
	prompt      string
	items       []selectItem
	manualLabel string

	// cursor は今どの行を指しているか。矢印キーでこの整数を増減させるだけで
	// 「選択中の行が動く」挙動を表現する（View 側で cursor の行に "> " を付ける）。
	cursor    int
	inputting bool // true の間は手動入力モード（textinput に委譲）
	ti        textinput.Model

	result  string // 確定した値
	aborted bool   // Ctrl+C 等で中断されたか
}

func newSelectModel(prompt string, items []selectItem, manualLabel string) selectModel {
	ti := textinput.New()
	ti.Placeholder = manualLabel
	ti.Prompt = "> "
	return selectModel{
		prompt:      prompt,
		items:       items,
		manualLabel: manualLabel,
		ti:          ti,
	}
}

// manualIndex は「手動入力」項目のカーソル位置（リストの末尾）
func (m selectModel) manualIndex() int {
	return len(m.items)
}

func (m selectModel) Init() tea.Cmd {
	return nil
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	// 手動入力モード中は、確定/中断/戻る以外のキーはすべて textinput に委譲する。
	// 文字入力・カーソル移動・バックスペースの面倒な処理を自作せず部品に任せることで、
	// 旧実装にあった「文字を消しても消えきらない」バグを構造的に排除している。
	if m.inputting {
		switch keyMsg.Type {
		case tea.KeyCtrlC:
			m.aborted = true
			return m, tea.Quit
		case tea.KeyEsc:
			// 入力をやめてリスト選択に戻る
			m.inputting = false
			m.ti.Blur()
			return m, nil
		case tea.KeyEnter:
			m.result = strings.TrimSpace(m.ti.Value())
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.ti, cmd = m.ti.Update(msg)
			return m, cmd
		}
	}

	// リスト選択モード: 矢印キーで cursor を動かし、Enter で確定する。
	switch keyMsg.String() {
	case "ctrl+c", "q":
		m.aborted = true
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 { // 先頭より上には行かない
			m.cursor--
		}
	case "down", "j":
		if m.cursor < m.manualIndex() { // 末尾（手動入力項目）より下には行かない
			m.cursor++
		}
	case "enter":
		if m.cursor == m.manualIndex() {
			// 末尾の項目 = 手動入力。入力モードへ遷移する。
			m.inputting = true
			m.ti.Focus()
			return m, textinput.Blink
		}
		m.result = m.items[m.cursor].value
		return m, tea.Quit // tea.Quit を返すと bubbletea のループが終了する
	}
	return m, nil
}

func (m selectModel) View() string {
	var b strings.Builder
	b.WriteString(promptStyle.Render(m.prompt))
	b.WriteString("\n\n")

	for i, item := range m.items {
		b.WriteString(m.renderRow(i, item.label))
	}
	// 手動入力の項目
	b.WriteString(m.renderRow(m.manualIndex(), m.manualLabel))

	if m.inputting {
		b.WriteString("\n")
		b.WriteString(m.ti.View())
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("enter: confirm • esc: back"))
	} else {
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("↑/↓: move • enter: select • q: quit"))
	}
	b.WriteString("\n")
	return b.String()
}

func (m selectModel) renderRow(index int, label string) string {
	if m.cursor == index && !m.inputting {
		return cursorStyle.Render(fmt.Sprintf("> %s", label)) + "\n"
	}
	return fmt.Sprintf("  %s\n", label)
}

// inputModel はタイトル入力用の単一行テキスト入力モデル。
// 行編集は bubbles/textinput に完全に委譲しており、バックスペースやカーソル移動、
// マルチバイト文字（日本語）の扱いを自前実装しないことでバグの温床をなくしている。
type inputModel struct {
	prompt  string
	ti      textinput.Model
	result  string
	aborted bool
}

func newInputModel(prompt string) inputModel {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.Focus()
	return inputModel{
		prompt: prompt,
		ti:     ti,
	}
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 確定(Enter)と中断(Ctrl+C/Esc)だけこちらで横取りし、
	// それ以外の入力キーはすべて textinput に流して編集させる。
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.aborted = true
			return m, tea.Quit
		case tea.KeyEnter:
			m.result = strings.TrimSpace(m.ti.Value())
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	var b strings.Builder
	b.WriteString(promptStyle.Render(m.prompt))
	b.WriteString("\n\n")
	b.WriteString(m.ti.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("enter: confirm • esc: quit"))
	b.WriteString("\n")
	return b.String()
}

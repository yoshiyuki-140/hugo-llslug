package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/yoshiyuki-140/hugo-llslug/internal/domain"
	"github.com/yoshiyuki-140/hugo-llslug/internal/usecase"
)

// redoCommand はSlug候補を再生成するための特別な選択結果を表す
const redoCommand = "/redo"

// Prompter はユーザとの対話的な入力を抽象化するインターフェース。
//
// なぜ切り出すか:
// Run() は「セクション決定 → タイトル入力 → 候補生成 → Slug選択 → hugo実行」という
// 手順だけを持ち、入力を"どう取るか"は関心の外に置きたい。矢印キー操作は
// 従来の bufio.Reader への文字列流し込みでは表現しづらくテストしにくいため、
// 入力手段をこのインターフェース越しに差し替え可能にしている。
//   - 本番: bubbleTeaPrompter（矢印キーのTUI）
//   - テスト: モック（決め打ちの値を返す）
//
// これにより端末に触れずに Run() のオーケストレーションだけを検証できる。
type Prompter interface {
	// SelectSection はセクションを一覧から選択、または手動入力して返す
	SelectSection(sections []string) (string, error)
	// InputTitle は記事タイトルの入力を受け取る
	InputTitle() (string, error)
	// SelectSlug はSlug候補を一覧から選択、または手動入力して返す。
	// 再生成が選ばれた場合は redoCommand ("/redo") を返す。
	SelectSlug(candidates []string) (string, error)
}

type Runner struct {
	uc          *usecase.SlugUsecase
	writer      io.Writer
	getSections func() []string
	prompter    Prompter
}

func NewRunner(uc *usecase.SlugUsecase) *Runner {
	return NewRunnerWithDeps(uc, os.Stdout, getContentSections, newBubbleTeaPrompter(os.Stdout))
}

// NewRunnerWithDeps はテスト時に依存を差し替えるためのコンストラクタ
func NewRunnerWithDeps(uc *usecase.SlugUsecase, out io.Writer, getSections func() []string, prompter Prompter) *Runner {
	return &Runner{
		uc:          uc,
		writer:      out,
		getSections: getSections,
		prompter:    prompter,
	}
}

func (r *Runner) Run() error {
	section, err := r.prompter.SelectSection(r.getSections())
	if err != nil {
		return err
	}

	title, err := r.prompter.InputTitle()
	if err != nil {
		return err
	}

	// Slugを再生成する処理。
	// ユーザが「再生成」を選ぶ限り候補生成→選択を繰り返したいので、
	// SelectSlug が redoCommand を返す間はループを続ける。
	var slug string
	for {
		candidates, err := r.generateCandidates(title)
		if err != nil {
			return err
		}

		slug, err = r.prompter.SelectSlug(candidates)
		if err != nil {
			return err
		}
		if slug != redoCommand {
			break
		}
	}

	fmt.Fprintln(r.writer, "Executing Hugo Command ...")
	fmt.Fprintf(r.writer, "`hugo new %s/%s/index.md`\n", section, slug)

	if err := r.uc.RunHugoNew(section, slug); err != nil {
		return fmt.Errorf("Hugo Execution Error: %w", err)
	}

	fmt.Fprintln(r.writer, "Completed !")
	return nil
}

// generateCandidates はSlug候補を生成する。LLMの応答が不正な場合は最大回数までリトライする。
func (r *Runner) generateCandidates(title string) ([]string, error) {
	const maxRetries = 10
	var candidates []string
	var err error
	for i := range maxRetries {
		fmt.Fprintln(r.writer, "Generating ...")
		candidates, err = r.uc.GetCandidates(title)
		if err == nil {
			return candidates, nil
		}
		if !errors.Is(err, domain.ErrLLMResponseParse) && !errors.Is(err, domain.ErrInvalidSlugFormat) {
			return nil, fmt.Errorf("Slug generate Error: %w", err)
		}
		if i < maxRetries-1 {
			fmt.Fprintf(r.writer, "The Result is Invalid. Regenerating... (%d/%d)\n", i+1, maxRetries)
		}
	}
	return nil, fmt.Errorf("Slug generate Error: %w", err)
}

// getContentSections は ./content/ 配下のトップレベルディレクトリ一覧を返す
func getContentSections() []string {
	entries, err := os.ReadDir("content")
	if err != nil {
		return nil
	}
	var sections []string
	for _, e := range entries {
		if e.IsDir() {
			sections = append(sections, e.Name())
		}
	}
	return sections
}

package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/yoshiyuki-140/hugo-llslug/internal/domain"
	"github.com/yoshiyuki-140/hugo-llslug/internal/usecase"
)

type Runner struct {
	uc          *usecase.SlugUsecase
	reader      *bufio.Reader
	writer      io.Writer
	getSections func() []string
}

func NewRunner(uc *usecase.SlugUsecase) *Runner {
	return NewRunnerWithDeps(uc, os.Stdin, os.Stdout, getContentSections)
}

// NewRunnerWithDeps はテスト時に依存を差し替えるためのコンストラクタ
func NewRunnerWithDeps(uc *usecase.SlugUsecase, in io.Reader, out io.Writer, getSections func() []string) *Runner {
	return &Runner{
		uc:          uc,
		reader:      bufio.NewReader(in),
		writer:      out,
		getSections: getSections,
	}
}

func (r *Runner) Run() error {
	section, err := r.selectSection()
	if err != nil {
		return err
	}

	fmt.Fprintln(r.writer, "追加したい記事のタイトルを入力してください．")
	fmt.Fprint(r.writer, "> ")
	title, _ := r.reader.ReadString('\n')
	title = strings.TrimSpace(title)

	const maxRetries = 3
	var candidates []string
	for i := range maxRetries {
		fmt.Fprintln(r.writer, "Generating ...")
		candidates, err = r.uc.GetCandidates(title)
		if err == nil {
			break
		}
		if !errors.Is(err, domain.ErrLLMResponseParse) && !errors.Is(err, domain.ErrInvalidSlugFormat) {
			return fmt.Errorf("スラッグ生成エラー: %w", err)
		}
		if i < maxRetries-1 {
			fmt.Fprintf(r.writer, "生成結果が不正でした。再生成します... (%d/%d)\n", i+1, maxRetries)
		}
	}
	if err != nil {
		return fmt.Errorf("スラッグ生成エラー: %w", err)
	}

	slug, err := r.selectSlug(candidates)
	if err != nil {
		return err
	}

	fmt.Fprintln(r.writer, "Executing Hugo Command ...")
	fmt.Fprintf(r.writer, "`hugo new %s/%s/index.md`\n", section, slug)

	if err := r.uc.RunHugoNew(section, slug); err != nil {
		return fmt.Errorf("Hugo実行エラー: %w", err)
	}

	fmt.Fprintln(r.writer, "Completed !")
	return nil
}

func (r *Runner) selectSection() (string, error) {
	sections := r.getSections()

	fmt.Fprintln(r.writer, "追加したい記事のセクション名を選択もしくは入力してください．")
	for i, s := range sections {
		fmt.Fprintf(r.writer, "    %d. %s\n", i+1, s)
	}
	fmt.Fprint(r.writer, "(入力する) > ")

	input, _ := r.reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if num, err := strconv.Atoi(input); err == nil {
		if num >= 1 && num <= len(sections) {
			return sections[num-1], nil
		}
	}

	return input, nil
}

func (r *Runner) selectSlug(candidates []string) (string, error) {
	fmt.Fprintln(r.writer, "Select following Slugs")
	for i, s := range candidates {
		fmt.Fprintf(r.writer, "    %d. %s\n", i+1, s)
	}
	fmt.Fprint(r.writer, "(入力する) > ")

	input, _ := r.reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if num, err := strconv.Atoi(input); err == nil {
		if num >= 1 && num <= len(candidates) {
			return candidates[num-1], nil
		}
	}

	return input, nil
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

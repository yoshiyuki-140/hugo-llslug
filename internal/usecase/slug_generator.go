package usecase

import (
	_ "embed"
	"strings"

	"github.com/yoshiyuki-140/hugo-llslug/internal/domain"
)

//go:embed prompts/slug_generate_instruction.txt
var slugInstructionTemplate string

type SlugUsecase struct {
	llmClient    domain.LLMClient
	hugoExecutor domain.HugoExecutor
}

// 外部からLLMClientをDIしてインスタンスを作成する
func NewSlugUsecase(client domain.LLMClient, hugo domain.HugoExecutor) *SlugUsecase {
	return &SlugUsecase{llmClient: client, hugoExecutor: hugo}
}

// タイトルから候補を5つ生成する
func (u *SlugUsecase) GetCandidates(title string) ([]string, error) {
	// バリデーション
	if title == "" {
		return nil, domain.ErrEmptyTitle
	}
	systemPrompt := strings.Replace(slugInstructionTemplate, "{{.ArticleTitle}}", title, 1)
	// TODO: タイトルのトークン数がLLMのコンテキスト長で耐えれるかをバリデーション
	return u.llmClient.GenerateSlugCandidates(systemPrompt)
}

// ユーザが選んだ最終結果を元にHugoを叩く
func (u *SlugUsecase) RunHugoNew(section string, selectedSlug string) error {
	// バリデーション
	if section == "" {
		return domain.ErrEmptySectionName
	}
	if selectedSlug == "" {
		return domain.ErrEmptySelectedSlug
	}
	if strings.Contains(section, "/") {
		return domain.ErrInvalidSectionName
	}
	if strings.Contains(selectedSlug, "/") {
		return domain.ErrInvalidSelectedSlug
	}
	return u.hugoExecutor.CreateNewPost(section, selectedSlug)
}

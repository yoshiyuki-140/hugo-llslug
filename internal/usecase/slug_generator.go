package usecase

import "github.com/yoshiyuki-140/hugo-llslug/internal/domain"

type SlugUsecase struct {
	LLMClient domain.LLMClient
}

// 外部からLLMClientをDIしてインスタンスを作成する
func NewSlugUsecase(client domain.LLMClient) *SlugUsecase {
	return &SlugUsecase{LLMClient: client}
}

// Executeはタイトルからスラッグを生成するビジネスロジック
func (u *SlugUsecase) Execute(title string) (string, error) {
	// タイトルのバリデーション
	// TODO: nilではなくて適切なエラーを返却させる
	if title == "" {
		return "", nil
	}
	// LLMClientを使ってスラッグを生成
	slug, err := u.LLMClient.GenerateSlug(title)
	if err != nil {
		return "", err
	}

	return slug, err
}

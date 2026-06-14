package domain

// LLClientはローカルLLMプラットフォームを抽象化するインターフェース

type LLMClient interface {
	GenerateSlug(title string) (string, error)
}

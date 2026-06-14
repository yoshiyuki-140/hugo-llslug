package domain

// Hugoの外部コマンド実行のためのインターフェース
type HugoExecutor interface {
	CreateNewPost(section string, slug string) error
}

// LLClientはローカルLLMプラットフォームを抽象化するインターフェース
type LLMClient interface {
	GenerateSlugCandidates(systemPrompt string) ([]string, error)
}

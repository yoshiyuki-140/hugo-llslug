package domain

// Hugoの外部コマンド実行のためのインターフェース
type HugoExecutor interface {
	CreateNewPost(section string, slug string) error
}

// LLMClientはローカルLLMプラットフォームを抽象化するインターフェース
type LLMClient interface {
	GenerateSlugCandidates(systemPrompt string) ([]string, error)
}

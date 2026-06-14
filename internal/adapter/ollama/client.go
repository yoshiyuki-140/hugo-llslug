package ollama

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ollama/ollama/api"
)

// Ollamaのレスポンスを受け取るための構造体
type slugResponse struct {
	Slugs []string `json:"slugs"`
}

type LLMClient struct {
	apiClient *api.Client
	modelName string
}

// Ollama APIクライアントのインスタンスを生成する
func NewClient(modelName string) (*LLMClient, error) {
	apiClient, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama client: %w", err)
	}

	if modelName == "" {
		modelName = "qwen3.5:0.8b" // default model name
	}

	return &LLMClient{
		apiClient: apiClient,
		modelName: modelName,
	}, nil
}

// domain.LLMClientインターフェースを実装する
func (c *LLMClient) GenerateSlugCandidates(systemPrompt string) ([]string, error) {
	ctx := context.Background()

	// Ollamaへのリクエスト組み立て
	req := &api.GenerateRequest{
		Model:  c.modelName,
		Prompt: systemPrompt,
		Format: json.RawMessage("json"), // 強制的にJSONを出力させるモードを有効化
		Options: map[string]interface{}{
			"temperature": 0.2, // 遊びを減らして，指示通りのフォーマットに固める
		},
	}

	var rawResponse string
	// OllamaのAPIを呼び出す
	respFunc := func(resp api.GenerateResponse) error {
		rawResponse += resp.Response
		return nil
	}

	if err := c.apiClient.Generate(ctx, req, respFunc); err != nil {
		// 接続エラー
		return nil, fmt.Errorf("ollama api error: %w", err)
	}

	// 返却されたJSON文字列を構造体にパース
	var slugRes slugResponse
	if err := json.Unmarshal(json.RawMessage(rawResponse), &slugRes); err != nil {
		return nil, fmt.Errorf("failed to parse llm json response: %w", err)
	}
	return slugRes.Slugs, nil
}

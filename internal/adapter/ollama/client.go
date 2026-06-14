package ollama

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type slugResponse struct {
	Slugs []string `json:"slugs"`
}

// CommandRunner abstracts exec.Command to allow injection in tests.
type CommandRunner func(name string, args ...string) ([]byte, error)

func defaultRunner(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

type LLMClient struct {
	modelName string
	runner    CommandRunner
}

func NewClient(modelName string) *LLMClient {
	return NewClientWithRunner(modelName, defaultRunner)
}

// テスト時に CommandRunner を差し替えるためのコンストラクタ
func NewClientWithRunner(modelName string, runner CommandRunner) *LLMClient {
	if modelName == "" {
		modelName = "qwen3.5:0.8b"
	}
	return &LLMClient{
		modelName: modelName,
		runner:    runner,
	}
}

// domain.LLMClientインターフェースを実装する
func (c *LLMClient) GenerateSlugCandidates(systemPrompt string) ([]string, error) {
	out, err := c.runner("ollama", "run", c.modelName, "--think=false", "--format", "json", systemPrompt)
	if err != nil {
		return nil, fmt.Errorf("ollama command error: %w", err)
	}

	var slugRes slugResponse
	if err := json.Unmarshal(out, &slugRes); err != nil {
		return nil, fmt.Errorf("failed to parse llm response: %w", err)
	}
	return slugRes.Slugs, nil
}

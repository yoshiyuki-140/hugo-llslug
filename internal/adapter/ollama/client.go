package ollama

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/yoshiyuki-140/hugo-llslug/internal/domain"
)

var kebabCaseRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// IsKebabCase はスラッグがケバブケース形式か否かを返す。
// 有効例: "slug", "slug-one", "hello-world-123"
// 無効例: "Slug", "slug_one", "-slug", "slug-"
func IsKebabCase(s string) bool {
	return kebabCaseRegex.MatchString(s)
}

// Goでは「代入先がインターフェース型のとき、右辺の型がそのインターフェースを満たしていなければコンパイルエラー」という規則がある。
var _ domain.LLMClient = (*LLMClient)(nil)

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
		modelName = "liquidai/lfm2.5-350m:q4_0"
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
		return nil, fmt.Errorf("%w: %w", domain.ErrLLMResponseParse, err)
	}

	for _, s := range slugRes.Slugs {
		if !IsKebabCase(s) {
			return nil, fmt.Errorf("%w: %q", domain.ErrInvalidSlugFormat, s)
		}
	}
	return slugRes.Slugs, nil
}

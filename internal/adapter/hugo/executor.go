package hugo

import (
	"fmt"
	"os/exec"

	"github.com/yoshiyuki-140/hugo-llslug/internal/domain"
)

// Goでは「代入先がインターフェース型のとき、右辺の型がそのインターフェースを満たしていなければコンパイルエラー」という規則がある。
var _ domain.HugoExecutor = (*HugoExecutor)(nil)

type CommandRunner func(name string, args ...string) ([]byte, error)

func defaultRunner(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

type HugoExecutor struct {
	runner CommandRunner
}

func NewExecutor() *HugoExecutor {
	return NewExecutorWithRunner(defaultRunner)
}

func NewExecutorWithRunner(runner CommandRunner) *HugoExecutor {
	return &HugoExecutor{runner: runner}
}

// domain.HugoExecutor インターフェースを実装する
func (e *HugoExecutor) CreateNewPost(section string, slug string) error {
	contentPath := fmt.Sprintf("%s/%s/index.md", section, slug)
	out, err := e.runner("hugo", "new", contentPath)
	if err != nil {
		return fmt.Errorf("hugo command error: %w\n%s", err, string(out))
	}
	return nil
}

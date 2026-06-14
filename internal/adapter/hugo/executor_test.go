package hugo_test

import (
	"errors"
	"testing"

	"github.com/yoshiyuki-140/hugo-llslug/internal/adapter/hugo"
)

func TestExecutor_CreateNewPost(t *testing.T) {
	tests := []struct {
		name      string
		section   string
		slug      string
		runnerErr error
		wantCmd   string
		wantArgs  []string
		wantErr   bool
	}{
		{
			name:     "正常にhugoコマンドを実行する",
			section:  "posts",
			slug:     "my-article",
			wantCmd:  "hugo",
			wantArgs: []string{"new", "posts/my-article/index.md"},
		},
		{
			name:      "コマンド失敗時はエラーを返す",
			section:   "posts",
			slug:      "my-article",
			runnerErr: errors.New("hugo not found"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotCmd string
			var gotArgs []string
			mockRunner := func(name string, args ...string) ([]byte, error) {
				gotCmd = name
				gotArgs = args
				return nil, tt.runnerErr
			}
			e := hugo.NewExecutorWithRunner(mockRunner)

			err := e.CreateNewPost(tt.section, tt.slug)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateNewPost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if gotCmd != tt.wantCmd {
					t.Errorf("cmd = %q, want %q", gotCmd, tt.wantCmd)
				}
				if len(gotArgs) != len(tt.wantArgs) {
					t.Fatalf("args = %v, want %v", gotArgs, tt.wantArgs)
				}
				for i, a := range gotArgs {
					if a != tt.wantArgs[i] {
						t.Errorf("args[%d] = %q, want %q", i, a, tt.wantArgs[i])
					}
				}
			}
		})
	}
}

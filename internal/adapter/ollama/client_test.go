package ollama_test

import (
	"errors"
	"testing"

	"github.com/yoshiyuki-140/hugo-llslug/internal/adapter/ollama"
	"github.com/yoshiyuki-140/hugo-llslug/internal/domain"
)

func TestIsKebabCase(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"slug", true},
		{"slug-one", true},
		{"hello-world-123", true},
		{"a1b2-c3d4", true},
		{"", false},
		{"Slug", false},
		{"SLUG", false},
		{"slug_one", false},
		{"-slug", false},
		{"slug-", false},
		{"slug--one", false},
		{"slug one", false},
		{"スラッグ", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ollama.IsKebabCase(tt.input); got != tt.want {
				t.Errorf("IsKebabCase(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestLLMClient_GenerateSlugCandidates(t *testing.T) {
	tests := []struct {
		name      string
		runnerOut []byte
		runnerErr error
		wantSlugs []string
		wantErr   bool
		wantErrIs error
	}{
		{
			name:      "正常にスラッグ候補を返す",
			runnerOut: []byte(`{"slugs":["slug-one","slug-two","slug-three","slug-four","slug-five"]}`),
			wantSlugs: []string{"slug-one", "slug-two", "slug-three", "slug-four", "slug-five"},
		},
		{
			name:      "コマンド失敗時はエラーを返す",
			runnerErr: errors.New("ollama not found"),
			wantErr:   true,
		},
		{
			name:      "JSONパース失敗時はErrLLMResponseParseを返す",
			runnerOut: []byte(`not json`),
			wantErr:   true,
			wantErrIs: domain.ErrLLMResponseParse,
		},
		{
			name:      "ケバブケース以外のスラッグが含まれる場合はErrInvalidSlugFormatを返す",
			runnerOut: []byte(`{"slugs":["valid-slug","InvalidSlug"]}`),
			wantErr:   true,
			wantErrIs: domain.ErrInvalidSlugFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRunner := func(_ string, _ ...string) ([]byte, error) {
				return tt.runnerOut, tt.runnerErr
			}
			c := ollama.NewClientWithRunner("test-model", mockRunner)

			got, err := c.GenerateSlugCandidates("test prompt")
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateSlugCandidates() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Fatalf("GenerateSlugCandidates() error = %v, want errors.Is(%v)", err, tt.wantErrIs)
			}
			if !tt.wantErr {
				if len(got) != len(tt.wantSlugs) {
					t.Fatalf("got %v, want %v", got, tt.wantSlugs)
				}
				for i, s := range got {
					if s != tt.wantSlugs[i] {
						t.Errorf("got[%d] = %q, want %q", i, s, tt.wantSlugs[i])
					}
				}
			}
		})
	}
}

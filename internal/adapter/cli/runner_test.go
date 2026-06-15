package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yoshiyuki-140/hugo-llslug/internal/adapter/cli"
	"github.com/yoshiyuki-140/hugo-llslug/internal/usecase"
)

type mockLLMClient struct {
	slugs []string
	err   error
}

func (m *mockLLMClient) GenerateSlugCandidates(_ string) ([]string, error) {
	return m.slugs, m.err
}

type mockHugoExecutor struct {
	capturedSection string
	capturedSlug    string
	err             error
}

func (m *mockHugoExecutor) CreateNewPost(section string, slug string) error {
	m.capturedSection = section
	m.capturedSlug = slug
	return m.err
}

func TestRunner_Run(t *testing.T) {
	defaultSlugs := []string{"slug-one", "slug-two", "slug-three", "slug-four", "slug-five"}
	defaultSections := []string{"posts", "notes"}

	tests := []struct {
		name        string
		input       string
		sections    []string
		slugs       []string
		llmErr      error
		hugoErr     error
		wantSection string
		wantSlug    string
		wantErr     bool
	}{
		{
			name:        "セクションを番号で選択し、スラッグを番号で選択する",
			input:       "1\ntest title\n1\n",
			sections:    defaultSections,
			slugs:       defaultSlugs,
			wantSection: "posts",
			wantSlug:    "slug-one",
		},
		{
			name:        "セクションを番号で選択し、スラッグを3番で選択する",
			input:       "2\ntest title\n3\n",
			sections:    defaultSections,
			slugs:       defaultSlugs,
			wantSection: "notes",
			wantSlug:    "slug-three",
		},
		{
			name:        "セクションを直接入力する",
			input:       "my-section\ntest title\n1\n",
			sections:    defaultSections,
			slugs:       defaultSlugs,
			wantSection: "my-section",
			wantSlug:    "slug-one",
		},
		{
			name:        "スラッグを直接入力する",
			input:       "1\ntest title\nmy-custom-slug\n",
			sections:    defaultSections,
			slugs:       defaultSlugs,
			wantSection: "posts",
			wantSlug:    "my-custom-slug",
		},
		{
			name:        "セクション一覧が空のとき直接入力になる",
			input:       "blog\ntest title\n1\n",
			sections:    []string{},
			slugs:       defaultSlugs,
			wantSection: "blog",
			wantSlug:    "slug-one",
		},
		{
			name:     "LLMがエラーを返すときRunはエラーを返す",
			input:    "1\ntest title\n",
			sections: defaultSections,
			slugs:    nil,
			llmErr:   errMockLLM,
			wantErr:  true,
		},
		{
			name:     "hugoコマンドがエラーを返すときRunはエラーを返す",
			input:    "1\ntest title\n1\n",
			sections: defaultSections,
			slugs:    defaultSlugs,
			hugoErr:  errMockHugo,
			wantErr:  true,
		},
		{
			name:        "/redoを入力すると再生成されその後番号で選択できる",
			input:       "1\ntest title\n/redo\n2\n",
			sections:    defaultSections,
			slugs:       defaultSlugs,
			wantSection: "posts",
			wantSlug:    "slug-two",
		},
		{
			name:        "/redoを複数回入力しても最終的に選択できる",
			input:       "1\ntest title\n/redo\n/redo\n1\n",
			sections:    defaultSections,
			slugs:       defaultSlugs,
			wantSection: "posts",
			wantSlug:    "slug-one",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hugo := &mockHugoExecutor{err: tt.hugoErr}
			uc := usecase.NewSlugUsecase(
				&mockLLMClient{slugs: tt.slugs, err: tt.llmErr},
				hugo,
			)
			runner := cli.NewRunnerWithDeps(
				uc,
				strings.NewReader(tt.input),
				&bytes.Buffer{},
				func() []string { return tt.sections },
			)

			err := runner.Run()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if hugo.capturedSection != tt.wantSection {
					t.Errorf("section = %q, want %q", hugo.capturedSection, tt.wantSection)
				}
				if hugo.capturedSlug != tt.wantSlug {
					t.Errorf("slug = %q, want %q", hugo.capturedSlug, tt.wantSlug)
				}
			}
		})
	}
}

var (
	errMockLLM  = mockError("llm error")
	errMockHugo = mockError("hugo error")
)

type mockError string

func (e mockError) Error() string { return string(e) }

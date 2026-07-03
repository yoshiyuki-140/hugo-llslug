package cli_test

import (
	"bytes"
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

// mockPrompter は対話的な入力を差し替えるためのモック。
// SelectSlug は呼び出しごとに slugResponses を順番に返す（/redo のテスト用）。
type mockPrompter struct {
	sectionResponse string
	sectionErr      error
	titleResponse   string
	titleErr        error
	slugResponses   []string
	slugErr         error

	slugCallCount int
}

func (m *mockPrompter) SelectSection(_ []string) (string, error) {
	return m.sectionResponse, m.sectionErr
}

func (m *mockPrompter) InputTitle() (string, error) {
	return m.titleResponse, m.titleErr
}

func (m *mockPrompter) SelectSlug(_ []string) (string, error) {
	if m.slugErr != nil {
		return "", m.slugErr
	}
	resp := m.slugResponses[m.slugCallCount]
	m.slugCallCount++
	return resp, nil
}

func TestRunner_Run(t *testing.T) {
	defaultSlugs := []string{"slug-one", "slug-two", "slug-three", "slug-four", "slug-five"}
	defaultSections := []string{"posts", "notes"}

	tests := []struct {
		name        string
		prompter    *mockPrompter
		sections    []string
		slugs       []string
		llmErr      error
		hugoErr     error
		wantSection string
		wantSlug    string
		wantErr     bool
	}{
		{
			name:        "セクションとスラッグを選択する",
			prompter:    &mockPrompter{sectionResponse: "posts", titleResponse: "test title", slugResponses: []string{"slug-one"}},
			sections:    defaultSections,
			slugs:       defaultSlugs,
			wantSection: "posts",
			wantSlug:    "slug-one",
		},
		{
			name:        "セクションを手動入力する",
			prompter:    &mockPrompter{sectionResponse: "my-section", titleResponse: "test title", slugResponses: []string{"slug-three"}},
			sections:    defaultSections,
			slugs:       defaultSlugs,
			wantSection: "my-section",
			wantSlug:    "slug-three",
		},
		{
			name:        "スラッグを手動入力する",
			prompter:    &mockPrompter{sectionResponse: "posts", titleResponse: "test title", slugResponses: []string{"my-custom-slug"}},
			sections:    defaultSections,
			slugs:       defaultSlugs,
			wantSection: "posts",
			wantSlug:    "my-custom-slug",
		},
		{
			name:     "LLMがエラーを返すときRunはエラーを返す",
			prompter: &mockPrompter{sectionResponse: "posts", titleResponse: "test title", slugResponses: []string{"slug-one"}},
			sections: defaultSections,
			slugs:    nil,
			llmErr:   errMockLLM,
			wantErr:  true,
		},
		{
			name:     "hugoコマンドがエラーを返すときRunはエラーを返す",
			prompter: &mockPrompter{sectionResponse: "posts", titleResponse: "test title", slugResponses: []string{"slug-one"}},
			sections: defaultSections,
			slugs:    defaultSlugs,
			hugoErr:  errMockHugo,
			wantErr:  true,
		},
		{
			name:        "/redoを選ぶと再生成されその後選択できる",
			prompter:    &mockPrompter{sectionResponse: "posts", titleResponse: "test title", slugResponses: []string{"/redo", "slug-two"}},
			sections:    defaultSections,
			slugs:       defaultSlugs,
			wantSection: "posts",
			wantSlug:    "slug-two",
		},
		{
			name:        "/redoを複数回選んでも最終的に選択できる",
			prompter:    &mockPrompter{sectionResponse: "posts", titleResponse: "test title", slugResponses: []string{"/redo", "/redo", "slug-one"}},
			sections:    defaultSections,
			slugs:       defaultSlugs,
			wantSection: "posts",
			wantSlug:    "slug-one",
		},
		{
			name:     "セクション選択でエラーが起きるとRunはエラーを返す",
			prompter: &mockPrompter{sectionErr: errMockAbort},
			sections: defaultSections,
			slugs:    defaultSlugs,
			wantErr:  true,
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
				&bytes.Buffer{},
				func() []string { return tt.sections },
				tt.prompter,
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
	errMockLLM   = mockError("llm error")
	errMockHugo  = mockError("hugo error")
	errMockAbort = mockError("abort error")
)

type mockError string

func (e mockError) Error() string { return string(e) }

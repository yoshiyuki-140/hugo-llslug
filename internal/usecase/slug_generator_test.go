// NOTE: Goのコンパイルの制約を回避しつつレイヤ間の独立性を別パッケージからテストするために(usecase_test)に切り分けてる
package usecase_test

import (
	"testing"

	"github.com/yoshiyuki-140/hugo-llslug/internal/domain"
	"github.com/yoshiyuki-140/hugo-llslug/internal/usecase"
)

type mockLLMClient struct {
	mockResponse string
	mockErr      error
}

func (m *mockLLMClient) GenerateSlug(title string) (string, error) {
	return m.mockResponse, m.mockErr
}

func TestSlugUsecase_Execute(t *testing.T) {
	mockllmClient := &mockLLMClient{
		mockResponse: "test-slug",
		mockErr:      nil,
	}

	// UsecaseにモックをDI
	type fields struct {
		LLMClient domain.LLMClient
	}
	type args struct {
		title string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "日本語の文字列に対応できる",
			fields: fields{
				LLMClient: mockllmClient,
			},
			args: args{
				title: "テストタイトル",
			},
			want:    "test-slug",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// uc:usecase
			uc := usecase.NewSlugUsecase(tt.fields.LLMClient)
			got, err := uc.Execute(tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("SlugUsecase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SlugUsecase.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

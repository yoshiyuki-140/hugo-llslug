// NOTE: Goのコンパイルの制約を回避しつつレイヤ間の独立性を別パッケージからテストするために(usecase_test)に切り分けてる
package usecase_test

import (
	"reflect"
	"testing"

	"github.com/yoshiyuki-140/hugo-llslug/internal/domain"
	"github.com/yoshiyuki-140/hugo-llslug/internal/usecase"
)

type mockLLMClient struct {
	mockResponses []string
	mockErr       error
}

type mockHugoExecutor struct {
	mockErr error
}

// モックの実装
func (mlc *mockLLMClient) GenerateSlugCandidates(title string) ([]string, error) {
	return mlc.mockResponses, mlc.mockErr
}

func (mhe *mockHugoExecutor) CreateNewPost(section string, slug string) error {
	return mhe.mockErr
}

// モックのインスタンス化
var mockllmClient = &mockLLMClient{
	mockResponses: []string{
		"test-slug1",
		"test-slug2",
		"test-slug3",
		"test-slug4",
		"test-slug5",
	},
	mockErr: nil,
}
var mockhugoExecutor = &mockHugoExecutor{
	mockErr: nil,
}

func TestSlugUsecase_GetCandidates(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		client domain.LLMClient
		hugo   domain.HugoExecutor
		// Named input parameters for target function.
		title   string
		want    []string
		wantErr bool
	}{
		{
			name:   "候補を5つ生成できている",
			client: mockllmClient,
			hugo:   mockhugoExecutor,
			title:  "テストタイトル",
			want: []string{
				"test-slug1",
				"test-slug2",
				"test-slug3",
				"test-slug4",
				"test-slug5",
			},
			wantErr: false,
		},
		{
			name:    "タイトルが空文字の時はエラーを返す",
			client:  mockllmClient,
			hugo:    mockhugoExecutor,
			title:   "",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := usecase.NewSlugUsecase(tt.client, tt.hugo)
			got, gotErr := u.GetCandidates(tt.title) // テスト対象メソッドの呼び出し
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetCandidates() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetCandidates() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("GetCandidates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlugUsecase_RunHugoNew(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		client domain.LLMClient
		hugo   domain.HugoExecutor
		// Named input parameters for target function.
		section      string
		selectedSlug string
		wantErr      bool
	}{
		{
			name:         "コマンド実行結果はnilである",
			client:       mockllmClient,
			hugo:         mockhugoExecutor,
			section:      "mock-section",
			selectedSlug: "mock-selected-slug",
			wantErr:      false,
		},
		{
			name:         "セクション名が空文字の時はエラー",
			client:       mockllmClient,
			hugo:         mockhugoExecutor,
			section:      "",
			selectedSlug: "mock-selected-slug",
			wantErr:      true,
		},
		{
			name:         "セクション名にスラッシュ(`/`)を含むときはエラー",
			client:       mockllmClient,
			hugo:         mockhugoExecutor,
			section:      "mock-/-section",
			selectedSlug: "mock-slug",
			wantErr:      true,
		},
		{
			name:         "Slugにスラッシュ(`/`)を含むときはエラー",
			client:       mockllmClient,
			hugo:         mockhugoExecutor,
			section:      "mock-section",
			selectedSlug: "mock-/-slug",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := usecase.NewSlugUsecase(tt.client, tt.hugo)
			gotErr := u.RunHugoNew(tt.section, tt.selectedSlug)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("RunHugoNew() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("RunHugoNew() succeeded unexpectedly")
			}
		})
	}
}

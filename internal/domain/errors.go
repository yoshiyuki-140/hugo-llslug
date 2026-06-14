package domain

import "errors"

var (
	// バリデーション系
	ErrEmptyTitle          = errors.New("記事のタイトルが空です")
	ErrEmptySectionName    = errors.New("セクション名が空です")
	ErrEmptySelectedSlug   = errors.New("Slugが空")
	ErrInvalidSectionName  = errors.New("セクション名にスラッシュ(`/`)等の文字は含められません")
	ErrInvalidSelectedSlug = errors.New("Slugにスラッシュ(`/`)等の文字は含められません")
	// ファイルロード系
	ErrCantLoadSystemPrompt = errors.New("システムプロンプトが読み込めませんでした")
)

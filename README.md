# llslug (local-llm-slug)

`llslug` は、ローカルLLM（Ollama）を活用して、日本語の記事タイトルからURLフレンドリーな英語のケバブケース（Slug）を自動生成する、Hugo専用のCLI拡張ツール（プラグイン）です。

Hugoの外部サブコマンド機構に対応しており、環境変数 `$PATH` に配置することで、お馴染みの `hugo` コマンドのサブコマンドとしてシームレスに実行できます。

## 使用方法

```bash
$ hugo llslug
追加したい記事のセクション名を選択もしくは入力してください．
    (入力する) > 1<Enter>
    1. posts
    2. class

追加したい記事のタイトルを入力してください．
> ローカルLLMを利用してHugoのSlugを自動生成してくれるツールを作ってみた
Generating ...
Select following Slugs
    (入力する) > 1<Enter>
    1. hugo-slug-auto-generator-local-llm
    2. hugo-slug-generator-local-llm
    3. auto-generate-hugo-slug-with-local-llm
    4. local-llm-hugo-slug-tool
    5. hugo-local-llm-slug

Executing Hugo Command ...
`hugo new posts/hugo-slug-auto-generator-local-llm/index.md`
Completed !
$ 
```

## 主な機能

- **爆速のローカル推論**: `llama3.2:1b` や `gemma2:2b` などの軽量ローカルLLMに対応し、1秒未満でレスポンスを返します。
- **Hugoとのシームレスな統合**: `hugo llslug` として実行可能（Cobraのプラグイン機構を利用）。
- **クリーンアーキテクチャ採用**: LLMプラットフォームやプロトコルの変更に強い、疎結合でテスタブルなGoコード。

## 前提条件

1. **Ollama** がローカルマシンで起動していること。
2. 使用する軽量LLMモデル（推奨: `llama3.2:1b`）がダウンロードされていること。
   ```bash
   ollama pull llama3.2:1b

## ディレクトリ構造

```bash
hugo-llslug/
├── .gitignore
├── go.mod
├── go.sum
├── main.go                     # エントリポイント
├── cmd/                        # CobraのCLI層（Delivery / Presentation）
│   ├── root.go                 # `hugo-llslug` ルートコマンド定義
│   └── version.go              # サブコマンド例（バージョン表示など）
│
└── internal/                   # 外部にエクスポートしないプライベートロジック
    ├── domain/                 # 【Domain層】ビジネスルール・インターフェース定義
    │   ├── slug.go             # スラッグ生成に関するインターフェース（LLMClient等）
    │   └── model.go            # ドメインモデル（必要なら）
    │
    ├── usecase/                # 【Usecase層】アプリケーションロジック（依存性を持たない）
    │   ├── slug_generator.go   # タイトルからスラッグを生成する一連の手順
    │   └── slug_generator_test.go # ➔ ★ Usecaseの単体テスト
    │
    └── adapter/                # 【Adapter層】外部インフラへの接続（依存性の外側）
        └── ollama/             # Ollamaクライアントの実装
            ├── client.go
            └── client_test.go  # ➔ ★ API通信のモックテスト
```

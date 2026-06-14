# llslug (local-llm-slug)

`llslug` は、ローカルLLM（Ollama）を活用して、日本語の記事タイトルからURLフレンドリーな英語のケバブケース（Slug）を自動生成する、Hugo専用のCLI拡張ツール（プラグイン）です。

Hugoの外部サブコマンド機構に対応しており、環境変数 `$PATH` に配置することで、 `hugo` のサブコマンドとして実行できます。

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

- ローカル推論: 軽量ローカルLLMを使用し、レスポンスを返します。
- Hugoとの統合**: `hugo llslug` として実行可能（Cobraのプラグイン機構を利用）。

## 前提条件

1. Ollamaがローカルマシンにインストールされていること．
2. 使用する軽量LLMモデル（推奨: `qwen3.5:0.8b`）がダウンロードされていること。

    もしダウンロードされていなければダウンロードしてください．
   ```bash
   ollama pull qwen3.5:0.8b
   ```

## ディレクトリ構造

```bash
hugo-llslug/
├── LICENSE
├── Makefile                            # ビルド・インストール用タスク定義
├── README.md
├── go.mod
├── go.sum
├── main.go                             # エントリポイント（cmd.Execute() を呼び出すのみ）
│
├── cmd/                                # Cobraコマンド定義（Presentation層）
│   └── root.go                         # ルートコマンド定義・フラグ設定・依存性の組み立て
│
└── internal/                           # 外部にエクスポートしないプライベートロジック
    ├── domain/                         # ビジネスルール・インターフェース定義
    │   ├── errors.go                   # ドメインエラー定数（バリデーション・ファイル読込系）
    │   └── hugo.go                     # HugoExecutor / LLMClient インターフェース定義
    │
    ├── usecase/                        # アプリケーションロジック
    │   ├── prompts/
    │   │   └── slug_generate_instruction.txt  # LLMへのシステムプロンプトテンプレート
    │   ├── slug_generator.go           # タイトル→スラッグ候補生成のユースケース実装
    │   └── slug_generator_test.go      # ユースケースの単体テスト
    │
    └── adapter/                        # 外部インフラへの接続実装
        ├── cli/                        # インタラクティブCLI入出力の処理
        │   ├── runner.go               # ユーザーとの対話フロー制御（セクション・タイトル・スラッグ選択）
        │   └── runner_test.go          # runner のテスト（stdin/stdout をモック）
        ├── hugo/                       # Hugo コマンド実行アダプタ
        │   ├── executor.go             # `hugo new` コマンドを実行する HugoExecutor 実装
        │   └── executor_test.go        # executor のテスト
        └── ollama/                     # Ollama クライアント実装
            ├── client.go               # Ollamaコマンドを呼び出す LLMClient 実装
            └── client_test.go          # API通信のモックテスト
```

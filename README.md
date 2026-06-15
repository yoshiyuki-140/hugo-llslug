# llslug (local-llm-slug)

`llslug` は，ローカルLLMを活用して，日本語の記事タイトルからURLフレンドリーな英語のケバブケース（Slug）を自動生成する，Hugo専用のCLI拡張です．
`hugo-llslug` として実行できます．

## 前提

1. [Ollama](https://github.com/ollama/ollama) がローカルマシンにインストールされ，起動していること．
2. [hugo](https://github.com/gohugoio/hugo)がインストールされていること．
3. 使用する軽量LLMモデル（推奨: `qwen3.5:0.8b` など）がダウンロードされていること．

```bash
# 推奨モデルのダウンロード
ollama pull qwen3.5:0.8b
```

## インストール

### Linux / macOS

```bash
curl -fsSL [https://raw.githubusercontent.com/yoshiyuki-140/hugo-llslug/main/install.sh](https://raw.githubusercontent.com/yoshiyuki-140/hugo-llslug/main/install.sh) | sh
```

### Windows (PowerShell)
環境のアーキテクチャ（デフォルトは `x86_64`）に合わせて実行してください．

```powershell
# x86_64 (標準)
powershell -Command "Invoke-WebRequest -Uri [https://raw.githubusercontent.com/yoshiyuki-140/hugo-llslug/main/install.ps1](https://raw.githubusercontent.com/yoshiyuki-140/hugo-llslug/main/install.ps1) -OutFile install.ps1; .\install.ps1"

# arm64 / i386 を指定する場合（-Arch オプションを追加）
# .\install.ps1 -Arch arm64
```

## 使い方

Hugoプロジェクトのルートディレクトリで実行します．

```bash
$ hugo-llslug
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
```

## ディレクトリ構造

```bash
hugo-llslug/
├── LICENSE
├── Makefile                           # ビルド・インストール用タスク定義
├── README.md
├── go.mod
├── go.sum
├── main.go                            # エントリポイント（cmd.Execute() の呼び出し）
│
├── cmd/                               # Cobraコマンド定義（Presentation層）
│   └── root.go                        # ルートコマンド定義・フラグ設定・DI
│
└── internal/                          # 外部非公開のプライベートロジック
    ├── domain/                        # ビジネスルール・インターフェース定義
    │   ├── errors.go                  # ドメインエラー定数
    │   └── hugo.go                    # HugoExecutor / LLMClient のインターフェース
    │
    ├── usecase/                       # アプリケーションロジック
    │   ├── prompts/
    │   │   └── slug_generate_instruction.txt # LLM用システムプロンプト
    │   ├── slug_generator.go          # タイトルからSlug候補を生成するユースケース
    │   └── slug_generator_test.go     # ユースケースの単体テスト
    │
    └── adapter/                       # 外部インフラ・技術駆動コードの実装
        ├── cli/                       # インタラクティブUI（入出力）処理
        │   ├── runner.go              # ユーザーとの対話フロー制御
        │   └── runner_test.go         # 擬似入出力を用いたテスト
        ├── hugo/                      # Hugo コマンド実行アダプタ
        │   ├── executor.go            # `hugo new` を叩く実体
        │   └── executor_test.go       # executor のテスト
        └── ollama/                    # Ollama クライアント実装
            ├── client.go              # OllamaのAPI/CLIを呼び出す実体
            └── client_test.go         # API通信のモックテスト
```
# llslug (local-llm-slug)

`llslug` は，ローカルLLMを活用して，日本語の記事タイトルからURLフレンドリーな英語のケバブケース（Slug）を自動生成する，Hugo専用のCLI拡張です．
`hugo-llslug` として実行できます．

## 前提

1. [Ollama](https://github.com/ollama/ollama) がローカルマシンにインストールされ，起動していること．
2. [hugo](https://github.com/gohugoio/hugo)がインストールされていること．
3. 使用する軽量LLMモデル（推奨: `liquidai/lfm2.5-350m:q4_0` など）がダウンロードされていること．

```bash
# 推奨モデルのダウンロード
ollama pull liquidai/lfm2.5-350m:q4_0
```

## インストール

### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/yoshiyuki-140/hugo-llslug/main/install.sh | sh
```

### Windows (PowerShell)
環境のアーキテクチャ（デフォルトは `x86_64`）に合わせて実行してください．

```powershell
# x86_64 (標準)
powershell -Command "Invoke-WebRequest -Uri https://raw.githubusercontent.com/yoshiyuki-140/hugo-llslug/main/install.ps1 -OutFile install.ps1; .\install.ps1"

# arm64 / i386 を指定する場合（-Arch オプションを追加）
# .\install.ps1 -Arch arm64
```

## 使い方

![demo video](./demo.gif)

Hugoプロジェクトのルートディレクトリで実行します．

```bash
$ hugo-llslug 
Please select or enter the section name for the new article.
    1. class
    2. posts
(Manual Input) > 2
Please enter the title of the article.
> ローカルLLM を使ったSlug生成ツールを作ってみた
Generating ...
Select following Slugs
    1. local-llm-used-slug-generator-tool
    2. create-global-slug-generation-tutorial
    3. my-first-local-generate-slugs
    4. global-slug-conversion-methods
    5. hugo-cli-using-local-languages
(`/redo` or Select Number) > 1
Executing Hugo Command ...
`hugo new posts/local-llm-used-slug-generator-tool/index.md`
Completed !
```

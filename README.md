# llslug (local-llm-slug)

[日本語READMEはこちら](./docs/README.ja.md)

`llslug` is a CLI extension designed specifically for Hugo. It leverages local LLMs to automatically generate URL-friendly, English kebab-case slugs from Japanese article titles. 
It can be executed as `hugo-llslug`.

## Prerequisites

1. [Ollama](https://github.com/ollama/ollama) must be installed and running on your local machine.
2. [Hugo](https://github.com/gohugoio/hugo) must be installed.
3. A lightweight LLM model (Recommended: `liquidai/lfm2.5-350m:q4_0`, etc.) must be downloaded.

```bash
# Download the recommended model
ollama pull liquidai/lfm2.5-350m:q4_0

```

## Installation

### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/yoshiyuki-140/hugo-llslug/main/install.sh | sh

```

### Windows (PowerShell)

Run the script matching your environment's architecture (Default is `x86_64`).

```powershell
# x86_64 (Standard)
powershell -Command "Invoke-WebRequest -Uri https://raw.githubusercontent.com/yoshiyuki-140/hugo-llslug/main/install.ps1 -OutFile install.ps1; .\install.ps1"

# To specify arm64 / i386 (Append the -Arch option)
# .\install.ps1 -Arch arm64
```

## Usage

![demo video](./docs/demo.gif)

Run the command in the root directory of your Hugo project.

```bash
$ hugo-llslug 
Please select or enter the section name for the new article.
    1. class
    2. posts
(Manual Input) > 2
Please enter the title of the article.
> Building a Local LLM-Powered Slug Generator
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

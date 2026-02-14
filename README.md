# Comic Translator CLI

這是一個基於 Golang 的漫畫翻譯 CLI 工具，支援 OCR 文字識別、機器翻譯以及嵌字輸出。

## 功能

*   **CLI 介面**：簡單易用的指令行工具。
*   **OCR**：使用 Tesseract (gosseract) 進行本地文字識別。
*   **翻譯**：
    *   `mock` (預設)：本地模擬翻譯，無需網路。
    *   `openai`：支援 OpenAI GPT-3.5/4 進行高質量翻譯 (需要 API Key)。
*   **嵌字**：自動抹除原文並填入翻譯後的文字。

## 安裝需求

### 1. 系統套件

請先安裝 Tesseract OCR 和開發庫：

**Ubuntu/Debian:**
```bash
sudo apt-get install -y tesseract-ocr libtesseract-dev tesseract-ocr-jpn tesseract-ocr-chi-tra fonts-noto-cjk
```

**macOS:**
```bash
brew install tesseract
# 需手動下載語言包或使用 brew 安裝
```

### 2. 編譯專案

```bash
go mod tidy
go build -o comic-translator main.go
```

## 使用方法

### 基本翻譯 (使用 Mock 引擎)

```bash
./comic-translator input.jpg
# 輸出: input_translated.jpg
```

### 使用 OpenAI 翻譯

```bash
./comic-translator input.jpg --engine=openai --api-key="YOUR_OPENAI_API_KEY" --source=jpn --target=chi_tra
```

### 參數說明

*   `--source, -s`: 來源語言代碼 (預設: `jpn`)。例如: `eng`, `jpn`, `chi_tra`。
*   `--target, -t`: 目標語言代碼 (預設: `chi_tra`)。
*   `--engine, -e`: 翻譯引擎 (預設: `mock`)。可選: `mock`, `openai`。
*   `--api-key`: 翻譯服務的 API Key (若使用 openai 則必填)。
*   `--output, -o`: 指定輸出檔案路徑。
*   `--font`: 指定用於嵌字的字型檔路徑 (預設使用系統 DroidSansFallbackFull.ttf)。

## 範例

將 `page01.png` 從日文翻譯成繁體中文：

```bash
./comic-translator page01.png -s jpn -t chi_tra -e openai --api-key="sk-..."
```

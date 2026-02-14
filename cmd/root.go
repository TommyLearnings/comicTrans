package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"comic-translator/pkg/ocr"
	"comic-translator/pkg/translator"
	"comic-translator/pkg/imager"
)

var (
	sourceLang string
	targetLang string
	outputFile string
	engine     string
	apiKey     string
	fontPath   string
)

var rootCmd = &cobra.Command{
	Use:   "comic-translator [image path]",
	Short: "Translate comic images",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		// Determine output path if not specified
		if outputFile == "" {
			ext := filepath.Ext(inputPath)
			name := strings.TrimSuffix(inputPath, ext)
			outputFile = name + "_translated" + ext
		}

		// If engine is openai and apiKey is empty, try to get from env
		if engine == "openai" && apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}

		fmt.Printf("Processing %s...\n", inputPath)

		// 1. Initialize Translator
		trans, err := translator.NewTranslator(engine, apiKey)
		if err != nil {
			return fmt.Errorf("failed to initialize translator: %v", err)
		}

		// 2. OCR Extraction
		fmt.Println("Extracting text...")
		blocks, err := ocr.ExtractText(inputPath, sourceLang)
		if err != nil {
			return fmt.Errorf("OCR failed: %v", err)
		}
		fmt.Printf("Found %d text blocks.\n", len(blocks))

		// 3. Translation
		fmt.Println("Translating text...")
		var translatedBlocks []ocr.TextBlock
		for _, block := range blocks {
			translatedText, err := trans.Translate(block.Text, sourceLang, targetLang)
			if err != nil {
				fmt.Printf("Warning: failed to translate block '%s': %v\n", block.Text, err)
				translatedText = block.Text // Keep original if translation fails
			}

			// Append with translated text
			newBlock := block
			newBlock.Text = translatedText
			translatedBlocks = append(translatedBlocks, newBlock)
		}

		// 4. Image Processing
		fmt.Printf("Generating output image to %s...\n", outputFile)
		if err := imager.ProcessImage(inputPath, outputFile, translatedBlocks, fontPath); err != nil {
			return fmt.Errorf("image processing failed: %v", err)
		}

		fmt.Println("Done!")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&sourceLang, "source", "s", "jpn", "Source language code (e.g., jpn, eng)")
	rootCmd.Flags().StringVarP(&targetLang, "target", "t", "chi_tra", "Target language code (e.g., chi_tra, eng)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	rootCmd.Flags().StringVarP(&engine, "engine", "e", "mock", "Translation engine (mock, openai)")
	rootCmd.Flags().StringVar(&apiKey, "api-key", "", "API Key for translation service")
	rootCmd.Flags().StringVar(&fontPath, "font", "", "Path to font file (optional)")
}

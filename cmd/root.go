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
	Use:   "comic-translator [image path or directory]",
	Short: "Translate comic images",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		info, err := os.Stat(inputPath)
		if err != nil {
			return fmt.Errorf("invalid path: %v", err)
		}

		// If engine is openai and apiKey is empty, try to get from env
		if engine == "openai" && apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}

		// 1. Initialize Translator once
		trans, err := translator.NewTranslator(engine, apiKey)
		if err != nil {
			return fmt.Errorf("failed to initialize translator: %v", err)
		}

		if info.IsDir() {
			// === Directory Mode ===
			entries, err := os.ReadDir(inputPath)
			if err != nil {
				return err
			}

			var files []string
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				name := entry.Name()
				ext := strings.ToLower(filepath.Ext(name))
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
					files = append(files, filepath.Join(inputPath, name))
				}
			}

			if len(files) == 0 {
				return fmt.Errorf("no supported image files found in directory %s", inputPath)
			}

			fmt.Printf("Processing directory: found %d files.\n", len(files))

			// Output Handling for Directory Mode
			// If --output is provided, treat it as the target directory.
			// If NOT provided, output to same directory with _translated suffix for files.
			var targetDir string
			if outputFile != "" {
				targetDir = outputFile
				if err := os.MkdirAll(targetDir, 0755); err != nil {
					return fmt.Errorf("failed to create output directory: %v", err)
				}
				fmt.Printf("Output directory set to: %s\n", targetDir)
			}

			for _, file := range files {
				// Determine output filename
				baseName := filepath.Base(file)
				ext := filepath.Ext(file)
				nameNoExt := strings.TrimSuffix(baseName, ext)

				var outPath string

				if targetDir != "" {
					// Use original filename (but force .png) inside target directory
					outPath = filepath.Join(targetDir, nameNoExt+".png")
				} else {
					// Suffix in same directory (force .png)
					outPath = filepath.Join(filepath.Dir(file), nameNoExt+"_translated.png")
				}

				if err := processSingleFile(file, outPath, trans); err != nil {
					fmt.Printf("Failed to process %s: %v\n", file, err)
				}
			}

		} else {
			// === Single File Mode ===
			var outPath string
			if outputFile != "" {
				outPath = outputFile
				// If user manually specified output but didn't end with .png,
				// warn them or just let it happen (imager saves as PNG regardless of ext).
				// For best practice, we can check.
				if !strings.HasSuffix(strings.ToLower(outPath), ".png") {
					fmt.Printf("Warning: output file %s does not end in .png, but will be saved as PNG format.\n", outPath)
				}
			} else {
				ext := filepath.Ext(inputPath)
				nameNoExt := strings.TrimSuffix(inputPath, ext)
				outPath = nameNoExt + "_translated.png"
			}

			if err := processSingleFile(inputPath, outPath, trans); err != nil {
				return err
			}
		}

		fmt.Println("All done!")
		return nil
	},
}

// processSingleFile encapsulates the logic for OCR -> Translate -> Image Process for one file.
func processSingleFile(inputPath string, outputPath string, trans translator.Translator) error {
	fmt.Printf("-> Processing %s...\n", inputPath)

	// 2. OCR Extraction
	blocks, err := ocr.ExtractText(inputPath, sourceLang)
	if err != nil {
		return fmt.Errorf("OCR error: %w", err)
	}

	// 3. Translation
	var translatedBlocks []ocr.TextBlock
	for _, block := range blocks {
		translatedText, err := trans.Translate(block.Text, sourceLang, targetLang)
		if err != nil {
			fmt.Printf("   Warning: failed to translate block '%s': %v\n", block.Text, err)
			translatedText = block.Text
		}

		newBlock := block
		newBlock.Text = translatedText
		translatedBlocks = append(translatedBlocks, newBlock)
	}

	// 4. Image Processing
	if err := imager.ProcessImage(inputPath, outputPath, translatedBlocks, fontPath); err != nil {
		return fmt.Errorf("image processing error: %w", err)
	}

	fmt.Printf("   Saved to %s\n", outputPath)
	return nil
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
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (single mode) or directory (batch mode)")
	rootCmd.Flags().StringVarP(&engine, "engine", "e", "mock", "Translation engine (mock, openai)")
	rootCmd.Flags().StringVar(&apiKey, "api-key", "", "API Key for translation service")
	rootCmd.Flags().StringVar(&fontPath, "font", "", "Path to font file (optional)")
}

package imager

import (
	"fmt"
	"image/color"
	"os"
	"github.com/fogleman/gg"
	"comic-translator/pkg/ocr"
)

// ProcessImage loads an image, draws over the specified blocks with translated text, and saves the result.
func ProcessImage(inputPath string, outputPath string, blocks []ocr.TextBlock, fontPath string) error {
	im, err := gg.LoadImage(inputPath)
	if err != nil {
		return err
	}

	dc := gg.NewContextForImage(im)

	if fontPath == "" {
		// Try to use Noto Sans CJK which is installed via fonts-noto-cjk
		// Check common paths
		candidates := []string{
			"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf",
		}

		for _, p := range candidates {
			if _, err := os.Stat(p); err == nil {
				fontPath = p
				break
			}
		}

		if fontPath == "" {
			return fmt.Errorf("no CJK font found. Please install fonts-noto-cjk or specify --font")
		}
	} else {
		if _, err := os.Stat(fontPath); err != nil {
			return fmt.Errorf("font file not found: %s", fontPath)
		}
	}

	for _, block := range blocks {
		x := float64(block.Box.Min.X)
		y := float64(block.Box.Min.Y)
		w := float64(block.Box.Max.X - block.Box.Min.X)
		h := float64(block.Box.Max.Y - block.Box.Min.Y)

		// Draw white background (simple erasing)
		dc.SetColor(color.White)
		dc.DrawRectangle(x, y, w, h)
		dc.Fill()

		// Determine font size based on box height
		// Heuristic: 60% of box height, minimum 12pt
		fontSize := h * 0.6
		if fontSize < 12 {
			fontSize = 12
		}

		// Load font
		if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
			return fmt.Errorf("failed to load font %s: %v", fontPath, err)
		}

		// Draw text centered in the box
		dc.SetColor(color.Black)
		dc.DrawStringWrapped(block.Text, x+w/2, y+h/2, 0.5, 0.5, w, 1.2, gg.AlignCenter)
	}

	return dc.SavePNG(outputPath)
}

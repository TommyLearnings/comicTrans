package ocr

import (
	"image"
	"github.com/otiai10/gosseract/v2"
)

// TextBlock represents a detected text block with its content and bounding box.
type TextBlock struct {
	Text string
	Box  image.Rectangle
}

// ExtractText extracts text from the image at the given path using the specified language.
// It returns a slice of TextBlock, where each block represents a line of text.
func ExtractText(imagePath string, lang string) ([]TextBlock, error) {
	client := gosseract.NewClient()
	defer client.Close()

	if err := client.SetImage(imagePath); err != nil {
		return nil, err
	}

	// Set language (e.g., "jpn", "chi_tra", "eng")
	if err := client.SetLanguage(lang); err != nil {
		return nil, err
	}

	// Get bounding boxes at the textline level
	// RIL_TEXTLINE provides lines of text which is suitable for bubbles.
	boxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
	if err != nil {
		return nil, err
	}

	var results []TextBlock
	for _, box := range boxes {
		// Clean up text if needed. Empty string means no text detected in that box.
		if box.Word == "" {
			continue
		}

		results = append(results, TextBlock{
			Text: box.Word,
			Box:  box.Box,
		})
	}

	return results, nil
}

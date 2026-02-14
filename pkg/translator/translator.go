package translator

import (
	"errors"
	"fmt"
)

// Translator defines the interface for translation services.
type Translator interface {
	// Translate translates the given text from source language to target language.
	// source and target should be ISO 639-1 codes (e.g., "ja", "zh", "en").
	Translate(text string, source string, target string) (string, error)
}

// NewTranslator creates a new translator instance based on the engine name.
func NewTranslator(engine string, apiKey string) (Translator, error) {
	switch engine {
	case "mock":
		return &MockTranslator{}, nil
	case "openai":
		if apiKey == "" {
			return nil, errors.New("api key is required for openai translator")
		}
		return &OpenAITranslator{APIKey: apiKey}, nil
	default:
		return nil, fmt.Errorf("unknown translator engine: %s", engine)
	}
}

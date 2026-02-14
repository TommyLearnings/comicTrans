package translator

// MockTranslator is a dummy translator for testing purposes.
type MockTranslator struct{}

// Translate implements the Translator interface.
func (m *MockTranslator) Translate(text, source, target string) (string, error) {
	return "[Translated] " + text, nil
}

package lapi

// Translator helps translation
type Translator interface {
	// Translate returns a translated string
	Translate(word string, args ...interface{}) string
}
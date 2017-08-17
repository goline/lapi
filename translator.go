package lapi

// Translator helps translation
type Translator interface {
	// T returns a translated string
	T(word string, args ...interface{}) string
}
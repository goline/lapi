package lapi

// Error represents for a common error
type Error interface {
	// Code returns error's code
	Code() string

	// Message returns error's message
	Message() string

	// Root returns original system error
	Root() error
}
package lapi

// Handler is a request's handler
type Handler interface {
	// Handle performs logic for solving request
	Handle(req Request, res Response)
}
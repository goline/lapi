package lapi

type Connection interface {
	// Request returns an instance of request
	Request() Request

	// WithRequest sets request
	WithRequest(request Request) Connection

	// Response returns an instance of response
	Response() Response

	// WithResponse sets response
	WithResponse(response Response) Connection
}

func NewConnection(request Request, response Response) Connection {
	return &FactoryConnection{request, response}
}

type FactoryConnection struct {
	request  Request
	response Response
}

func (c *FactoryConnection) Request() Request {
	return c.request
}

func (c *FactoryConnection) WithRequest(request Request) Connection {
	c.request = request
	return c
}

func (c *FactoryConnection) Response() Response {
	return c.response
}

func (c *FactoryConnection) WithResponse(response Response) Connection {
	c.response = response
	return c
}

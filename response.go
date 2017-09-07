package lapi

import (
	"fmt"
	"net/http"
)

// Response is a application's response
type Response interface {
	ResponseDescriber
	ResponseInformer
	ResponseCookies
	ResponseHeader
	ResponseSender
	ParserManager
}

type ResponseInformer interface {
	// Status gets HTTP status code
	Status() int

	// WithStatus sets HTTP status code
	WithStatus(status int) Response

	// Message returns HTTP status message
	Message() string

	// WithMessage sets HTTP status message
	WithMessage(message string) Response
}

// ResponseDescriber handles content
type ResponseDescriber interface {
	// Content gets response's content
	Content() interface{}

	// WithContent sets response's content
	WithContent(content interface{}) Response
}

// ResponseHeader helps to manage header
type ResponseHeader interface {
	// Header returns an instance of response's header
	Header() Header

	// WithHeader sets response's header
	WithHeader(header Header) Response
}

// ResponseCookies helps to manage cookies
type ResponseCookies interface {
	// Cookie returns an instance of response's cookies
	Cookies() []*http.Cookie

	// WithCookie sets response's cookies
	WithCookies(cookies []*http.Cookie) Response
}

// ResponseSender sends response to client
type ResponseSender interface {
	// Send flushes response out
	Send() error

	// IsSent returns true if response has already been sent
	// false if otherwise
	IsSent() bool
}

func NewResponse(w http.ResponseWriter) (Response, error) {
	return &FactoryResponse{
		w:             w,
		status:        http.StatusOK,
		header:        NewHeader(),
		isSent:        false,
		ParserManager: NewParserManager(),
	}, nil
}

type FactoryResponse struct {
	w       http.ResponseWriter
	status  int
	message string
	content interface{}
	header  Header
	cookies []*http.Cookie
	isSent  bool
	ParserManager
}

func (r *FactoryResponse) Status() int {
	return r.status
}

func (r *FactoryResponse) WithStatus(status int) Response {
	r.status = status
	return r
}

func (r *FactoryResponse) Message() string {
	return r.message
}

func (r *FactoryResponse) WithMessage(message string) Response {
	r.message = message
	return r
}

func (r *FactoryResponse) Content() interface{} {
	return r.content
}

func (r *FactoryResponse) WithContent(content interface{}) Response {
	r.content = content
	return r
}

func (r *FactoryResponse) Header() Header {
	return r.header
}

func (r *FactoryResponse) WithHeader(header Header) Response {
	r.header = header
	return r
}

func (r *FactoryResponse) Cookies() []*http.Cookie {
	return r.cookies
}

func (r *FactoryResponse) WithCookies(cookies []*http.Cookie) Response {
	r.cookies = cookies
	return r
}

func (r *FactoryResponse) Send() error {
	if r.isSent {
		return NewSystemError(ERROR_RESPONSE_ALREADY_SENT, "Response is already sent")
	}

	for _, cookie := range r.cookies {
		http.SetCookie(r.w, cookie)
	}

	for k, v := range r.header.All() {
		r.w.Header().Set(k, v)
	}

	ct, _ := r.header.Get("Content-Type")
	p, ok := r.Parser(ct)
	if ok == false {
		return NewSystemError(ERROR_NO_PARSER_FOUND, fmt.Sprintf("Unable to find an appropriate parser for %s", ct))
	}

	b, err := p.Encode(r.content)
	if err != nil {
		return err
	}

	if r.message != "" {
		http.Error(r.w, r.message, r.status)
	} else {
		r.w.WriteHeader(r.status)
	}
	r.w.Write(b)

	r.isSent = true
	return nil
}

func (r *FactoryResponse) IsSent() bool {
	return r.isSent
}

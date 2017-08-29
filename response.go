package lapi

import (
	"encoding/json"
	"net/http"
)

// Response is a application's response
type Response interface {
	ResponseDescriber
	ResponseInformer
	ResponseCookies
	ResponseHeader
	ResponseSender
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
}

func NewResponse(w http.ResponseWriter) Response {
	return &FactoryResponse{
		w:      w,
		status: http.StatusOK,
		header: NewHeader(),
	}
}

type FactoryResponse struct {
	w       http.ResponseWriter
	status  int
	message string
	content interface{}
	header  Header
	cookies []*http.Cookie
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
	for _, cookie := range r.cookies {
		http.SetCookie(r.w, cookie)
	}

	for k, v := range r.header.Lines() {
		r.w.Header().Set(k, v)
	}

	if r.message != "" {
		http.Error(r.w, r.message, r.status)
	} else {
		r.w.WriteHeader(r.status)
	}

	b, err := json.Marshal(r.content)
	if err != nil {
		return err
	}
	r.w.Write(b)

	return nil
}

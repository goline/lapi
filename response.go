package lapi

import (
	"fmt"
	"net/http"
)

// Response is a application's response
type Response interface {
	ResponseBody
	ResponseHeader
	ResponseSender
	ResponseCookies
	ResponseInformer
}

type ResponseBody Body

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
	// It requires Content-Type is set before changing its content,
	// as an error will be returned when putting content while content-type is empty
	Send() error

	// IsSent returns true if response has already been sent
	// false if otherwise
	IsSent() bool
}

func NewResponse(w http.ResponseWriter) Response {
	return &FactoryResponse{
		writer: w,
		status: http.StatusOK,
		header: NewHeader(),
		isSent: false,
		Body:   NewBody(),
	}
}

type FactoryResponse struct {
	writer  http.ResponseWriter
	status  int
	message string
	header  Header
	cookies []*http.Cookie
	isSent  bool
	Body
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
	if r.writer == nil {
		return NewSystemError(ERROR_NO_WRITER_FOUND, "No writer found")
	}

	if r.isSent {
		return NewSystemError(ERROR_RESPONSE_ALREADY_SENT, "Response is already sent")
	}

	contentType := r.ContentType()
	if contentType == "" {
		return NewSystemError(ERROR_CONTENT_TYPE_EMPTY, "Content-Type is required")
	}
	charset := r.Charset()
	if charset != "" {
		r.header.Set("content-type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	} else {
		r.header.Set("content-type", contentType)
	}

	for _, cookie := range r.cookies {
		http.SetCookie(r.writer, cookie)
	}

	for k, v := range r.header.All() {
		r.writer.Header().Set(k, v)
	}

	content, err := r.ContentBytes()
	if err != nil {
		return err
	}

	if r.message != "" {
		http.Error(r.writer, r.message, r.status)
	} else {
		r.writer.WriteHeader(r.status)
	}
	r.writer.Write(content)

	r.isSent = true
	return nil
}

func (r *FactoryResponse) IsSent() bool {
	return r.isSent
}

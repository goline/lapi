package lapi

import (
	"fmt"
	"net/http"

	"github.com/goline/errors"
)

// Response is a application's response
type Response interface {
	ResponseBody
	ResponseHeader
	ResponseSender
	ResponseCookies
	ResponseInformer
	ResponseAncestor
}

type ResponseBody interface {
	// Body returns an instance of Body
	Body() Body

	// WithBody sets body's instance
	WithBody(body Body) Response
}

// ResponseAncestor keeps original http.Request
type ResponseAncestor interface {
	Ancestor() http.ResponseWriter
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
		ancestor: w,
		status:   http.StatusOK,
		header:   NewHeader(),
		isSent:   false,
		body:     NewBody(nil, w),
	}
}

type FactoryResponse struct {
	ancestor  http.ResponseWriter
	status    int
	message   string
	header    Header
	cookies   []*http.Cookie
	body      Body
	isSent    bool
	isSending bool
}

func (r *FactoryResponse) Ancestor() http.ResponseWriter {
	return r.ancestor
}

func (r *FactoryResponse) Body() Body {
	return r.body
}

func (r *FactoryResponse) WithBody(body Body) Response {
	r.body = body
	return r
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
	if r.lock() == true {
		return errors.New(ERR_RESPONSE_IS_SENDING, "Sending response is in progress")
	}
	defer r.unlock()

	if r.ancestor == nil {
		return errors.New(ERR_NO_WRITER_FOUND, "No writer found")
	}

	if r.isSent {
		return errors.New(ERR_RESPONSE_ALREADY_SENT, "Response is already sent")
	}

	contentType := r.body.ContentType()
	if contentType != "" {
		charset := r.body.Charset()
		if charset != "" {
			r.header.Set("content-type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		} else {
			r.header.Set("content-type", contentType)
		}
	}

	for _, cookie := range r.cookies {
		http.SetCookie(r.ancestor, cookie)
	}

	for k, v := range r.header.All() {
		r.ancestor.Header().Set(k, v)
	}

	if r.message != "" {
		http.Error(r.ancestor, r.message, r.status)
	} else {
		r.ancestor.WriteHeader(r.status)
	}

	if err := r.body.Flush(); err != nil {
		return err
	}

	r.isSent = true
	return nil
}

func (r *FactoryResponse) IsSent() bool {
	return r.isSent
}

func (r *FactoryResponse) lock() bool {
	if r.isSending == true {
		return true
	}
	r.isSending = true
	return false
}

func (r *FactoryResponse) unlock() {
	r.isSending = false
}

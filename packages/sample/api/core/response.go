package core

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"unicode/utf8"
)

const (
	Content_TYPE_HEADER_KEY = "Content-Type"
)

// This assumes we are using a web action with raw=true,
type MainResponseArgs struct {
	StatusCode int                 `json:"statusCode"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
}

// MainResponseWriter implements http.ResponseWriter and adds the method
// necessary to return an MainResponseArgs object
type MainResponseWriter struct {
	headers   http.Header
	body      bytes.Buffer
	status    int
	observers []chan<- bool
}

// NewMainResponseWriter returns a new MainResponseWriter object.
// The object is initialized with an empty map of headers and a
func NewMainResponseWriter() *MainResponseWriter {
	return &MainResponseWriter{
		headers:   make(http.Header),
		status:    200,
		observers: make([]chan<- bool, 0),
	}

}

func (r *MainResponseWriter) CloseNotify() <-chan bool {
	ch := make(chan bool, 1)

	r.observers = append(r.observers, ch)

	return ch
}

func (r *MainResponseWriter) notifyClosed() {
	for _, v := range r.observers {
		v <- true
	}
}

// Header implementation from the http.ResponseWriter interface.
func (r *MainResponseWriter) Header() http.Header {
	return r.headers
}

// Write sets the response body in the object. If no status code
// was set before with the WriteHeader method it sets the status
// for the response to 200 OK.
func (r *MainResponseWriter) Write(body []byte) (int, error) {
	if r.status == 200 {
		r.status = http.StatusOK
	}

	// if the content type header is not set when we write the body we try to
	// detect one and set it by default. If the content type cannot be detected
	// it is automatically set to "application/octet-stream" by the
	// DetectContentType method
	if r.Header().Get(Content_TYPE_HEADER_KEY) == "" {
		r.Header().Add(Content_TYPE_HEADER_KEY, http.DetectContentType(body))
	}

	return (&r.body).Write(body)
}

// WriteHeader sets a status code for the response. This method is used
// for error responses.
func (r *MainResponseWriter) WriteHeader(status int) {
	r.status = status
}

// GetMainResponse converts the data passed to the response writer into
// an MainResponseArgs object.
// Returns a populated response object.
func (r *MainResponseWriter) GetMainResponse() (MainResponseArgs, error) {
	r.notifyClosed()

	var output string

	bb := (&r.body).Bytes()

	if utf8.Valid(bb) {
		output = string(bb)
	} else {
		output = base64.StdEncoding.EncodeToString(bb)
	}

	return MainResponseArgs{
		StatusCode: r.status,
		Headers:    http.Header(r.headers),
		Body:       output,
	}, nil
}

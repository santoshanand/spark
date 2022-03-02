package spark

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

type (
	// Response - response struct
	Response struct {
		writer    http.ResponseWriter
		status    int
		size      int64
		committed bool
	}
)

// NewResponse - new response function will take http.ResponseWriter and returns *Response
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{writer: w}
}

// SetWriter - SetWriter method
func (r *Response) SetWriter(w http.ResponseWriter) {
	r.writer = w
}

// Header - get header
func (r *Response) Header() http.Header {
	return r.writer.Header()
}

// Writer - writer method
func (r *Response) Writer() http.ResponseWriter {
	return r.writer
}

// WriteHeader - writter header
func (r *Response) WriteHeader(code int) {
	if r.committed {
		// TODO: Warning
		log.Printf("spark => %s", "response already committed")
		return
	}
	r.status = code
	r.writer.WriteHeader(code)
	r.committed = true
}

// Write - write
func (r *Response) Write(b []byte) (n int, err error) {
	n, err = r.writer.Write(b)
	r.size += int64(n)
	return n, err
}

// Flush wraps response writer's Flush function.
func (r *Response) Flush() {
	r.writer.(http.Flusher).Flush()
}

// Hijack wraps response writer's Hijack function.
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.writer.(http.Hijacker).Hijack()
}

// CloseNotify wraps response writer's CloseNotify function.
func (r *Response) CloseNotify() <-chan bool {
	return r.writer.(http.CloseNotifier).CloseNotify()
}

// Status - get status code
func (r *Response) Status() int {
	return r.status
}

// Size - response size
func (r *Response) Size() int64 {
	return r.size
}

func (r *Response) reset(w http.ResponseWriter) {
	r.writer = w
	r.size = 0
	r.status = http.StatusOK
	r.committed = false
}

func (r *Response) clear() {
	r.Header().Del(ContentType)
	r.committed = false
}

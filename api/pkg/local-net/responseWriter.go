package localnet

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	header http.Header
	Buf    bytes.Buffer
	StCode int
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		header: make(http.Header),
		Buf:    *bytes.NewBuffer(nil),
	}
}

func (r *ResponseWriter) Header() http.Header {
	return r.header
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	if r.StCode == 0 {
		r.StCode = http.StatusOK
	}
	return r.Buf.Write(b)
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.StCode = statusCode
}

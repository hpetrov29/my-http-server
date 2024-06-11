package http

import (
	"fmt"
	"io"
	"net"
)

type Header map[string][]string

func (h Header) Add(key, value string) {
	values := h[key]

	values = append(values, value)
	h[key] = values
}

func (h Header) Write(w io.Writer) error {
	for key, values := range h {
		for _, value := range values {
			_, err := fmt.Fprintf(w, "%s: %s\r\n", key, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type ResponseWriter interface {
	WriteHeader(statusCode int)
	Write([]byte) (int, error)
	Header() Header
}

type responseWriter struct {
	conn net.Conn
	header Header
	wroteHeader bool
}

func NewResponseWriter(conn net.Conn) ResponseWriter {
    return &responseWriter{conn: conn, header: make(map[string][]string), wroteHeader: false}
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if rw.wroteHeader {
		return
	}

    statusText := "OK"
    if statusCode == 404 {
        statusText = "Not Found"
    }

	// Write the status line
    rw.conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)))

	// Write the response headers
	rw.Header().Write(rw.conn)

	// Write a blank line to separate headers and body
	rw.conn.Write([]byte("\r\n"))

	rw.wroteHeader = true
}

func (rw *responseWriter) Write(data []byte) (int, error) {
    return rw.conn.Write(data)
}

func (rw *responseWriter) Header() Header {
    return rw.header
}
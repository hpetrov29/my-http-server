package http

import (
	"fmt"
	"net"
)

type ResponseWriter interface {
	WriteHeader(statusCode int)
	Write([]byte) (int, error)
}

type responseWriter struct {
	conn net.Conn
}

func NewResponseWriter(conn net.Conn) ResponseWriter {
    return &responseWriter{conn: conn}
}

func (rw *responseWriter) WriteHeader(statusCode int) {
    statusText := "OK"
    if statusCode == 404 {
        statusText = "Not Found"
    }
    rw.conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)))
}

func (rw *responseWriter) Write(data []byte) (int, error) {
    return rw.conn.Write(data)
}
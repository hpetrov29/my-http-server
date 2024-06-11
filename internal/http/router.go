package http

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/internal/utils"
)

type Handler func(ResponseWriter, Request)

type Router struct {
	routes map[string]Handler
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]Handler)}
}

func (r *Router) Handle(pattern string, h Handler) {
	r.routes[pattern] = h
}

func (r *Router) Serve(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
    requestLine, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading request: ", err.Error())
        return
    }

	// requestLine format
	// GET /echo/abc HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n
	parts := strings.Fields(requestLine)
	if len(parts) < 3 {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	method := parts[0]
	path := parts[1]

	// Only GET requests are allowed for now
	if method != "GET" {
        conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n\r\n"))
        return
    }

    // Read and parse headers
    headers := make(map[string]string)
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading header line: ", err.Error())
            return
        }

        line = strings.TrimSpace(line)
        if line == "" {
            break
        }

        headerParts := strings.SplitN(line, ": ", 2)
        if len(headerParts) == 2 {
            headers[headerParts[0]] = headerParts[1]
        }
    }

	req := Request{Method: method, Path: path, Headers: headers }
	res := NewResponseWriter(conn)

	for pattern, handler := range r.routes {
		if matches, params := utils.Match(pattern, path); matches {
			req.Params = params
			handler(res, req)
			return
		}
	}

	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}
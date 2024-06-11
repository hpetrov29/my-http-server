package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/internal/utils"
)

const (
	GET = "GET"
	POST = "POST"
)

type Handler func(ResponseWriter, Request)

type Router struct {
	routes map[string]map[string] Handler
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]map[string]Handler)}
}

func (r *Router) Handle(method string, pattern string, h Handler) {
	if r.routes[pattern] == nil {
		r.routes[pattern] = make(map[string]Handler)
	}

	r.routes[pattern][method] = h
}

func (r *Router) Serve(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// -------------- Request line --------------
	// POST /files/number HTTP/1.1
	// \r\n

	// -------------- Headers --------------
	// Host: localhost:4221\r\n
	// User-Agent: curl/7.64.1\r\n
	// Accept: */*\r\n
	// Content-Type: application/octet-stream  // Header that specifies the format of the request body
	// Content-Length: 5\r\n                   // Header that specifies the size of the request body, in bytes
	// \r\n

	// -------------- Request Body --------------
	// 12345

    requestLine, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading request: ", err.Error())
        return
    }

	parts := strings.Fields(requestLine)
	if len(parts) < 3 {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	reqMethod := parts[0]
	reqPath := parts[1]

	// Only GET and POST requests are allowed for now
	if reqMethod != "GET" && reqMethod != "POST" {
        conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n\r\n"))
        return
    }

    // Read and parse request headers
    reqHeaders := make(map[string]string)
    for {
        headerLine, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading header line: ", err.Error())
            return
        }

        headerLine = strings.TrimSpace(headerLine)
        if headerLine == "" {
            break
        }

        headerParts := strings.SplitN(headerLine, ": ", 2)
        if len(headerParts) == 2 {
            reqHeaders[headerParts[0]] = headerParts[1]
        }
    }

	// Read and parse request body
	reqBody := []byte{}
    if contentLengthStr, ok := reqHeaders["Content-Length"]; ok {
        contentLength, err := strconv.Atoi(contentLengthStr)
        if err != nil {
            fmt.Println("Invalid Content-Length: ", err.Error())
            return
        }

        reqBody = make([]byte, contentLength)
        _, err = io.ReadFull(reader, reqBody)
        if err != nil {
            fmt.Println("Error reading body: ", err.Error())
            return
        }
    }

	req := Request{Method: reqMethod, Path: reqPath, Headers: reqHeaders, Body: reqBody}
	res := NewResponseWriter(conn)

	for route, methods := range r.routes {
		if matches, params := utils.Match(route, reqPath); matches {
			for method, handler := range methods {
				if reqMethod == method {
					req.Params = params	
					handler(res, req)
					return
				}
			}
		}
	}

	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}
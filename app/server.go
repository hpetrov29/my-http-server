package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type handler func(net.Conn, []string)

type router struct{
	routes map[string]handler
}

func newRouter() *router {
	return &router{routes: make(map[string]handler)}
}

func (r *router) handle(pattern string, h handler) {
	r.routes[pattern] = h
}

func (r *router) serve(conn net.Conn) {
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

	for pattern, handler := range r.routes {
		if matches, params := match(pattern, path); matches {
			handler(conn, params)
			return
		}
	}

	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}

func match(pattern, path string) (bool, []string) {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if(len(patternParts) != len(pathParts)) {
		return false, nil
	}

	params := make([]string, 0, len(patternParts))
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			params = append(params, pathParts[i])
		} else if part != pathParts[i] {
			return false, nil
		}
	}

	return true, params
}

func homeHandler(conn net.Conn, params []string) {
	resposne := "HTTP/1.1 200 OK\r\n\r\n"
	_, err := conn.Write([]byte(resposne))
    if err != nil {
        fmt.Println("Error writing response: ", err.Error())
    }
}

func echoHandler(conn net.Conn, params []string) {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(params[0]), params[0])
	_, err := conn.Write([]byte(response))
    if err != nil {
        fmt.Println("Error writing response: ", err.Error())
    }
}

func main() {
		r := newRouter()
		r.handle("/", homeHandler)
		r.handle("/echo/:param", echoHandler)

		l, err := net.Listen("tcp", "0.0.0.0:4221")
		if err != nil {
			fmt.Println("Failed to bind to port 4221")
			os.Exit(1)
		}
		defer l.Close()

		fmt.Println("Listening on port 4221...")
		
		for {
			connection, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}

			go r.serve(connection)
		}
	}

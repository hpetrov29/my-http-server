package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/internal/http"
)


func homeHandler(res http.ResponseWriter, req http.Request) {
	resposne := "HTTP/1.1 200 OK\r\n\r\n"
	_, err := res.Write([]byte(resposne))
    if err != nil {
        fmt.Println("Error writing response: ", err.Error())
    }
}

func echoHandler(res http.ResponseWriter, req http.Request) {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(req.Params[0]), req.Params[0])
	_, err := res.Write([]byte(response))
    if err != nil {
        fmt.Println("Error writing response: ", err.Error())
    }
}

func userAgentHandler(res http.ResponseWriter, req http.Request) {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(req.Headers["User-Agent"]), req.Headers["User-Agent"])
	_, err := res.Write([]byte(response))
    if err != nil {
        fmt.Println("Error writing response: ", err.Error())
    }
}

func main() {
		r := http.NewRouter()
		r.Handle("/", homeHandler)
		r.Handle("/echo/:param", echoHandler)
		r.Handle("/user-agent", userAgentHandler)

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

			go r.Serve(connection)
		}
	}

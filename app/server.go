package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
    request, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading request: ", err.Error())
        return
    }

	target := strings.Split(request, " ")[1]

	switch target {
		case "/":
			_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
    		if err != nil {
        		fmt.Println("Error writing response: ", err.Error())
    		}
		default:
			_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
    		if err != nil {
        		fmt.Println("Error writing response: ", err.Error())
    		}
	}
}

func main() {
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

			go handleConnection(connection)
		}
	}

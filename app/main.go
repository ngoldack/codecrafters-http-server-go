package main

import (
	"fmt"
	"net"
	"os"
)

const httpVersion = "HTTP/1.1"

var statusCodes = map[int]string{
	200: "OK",
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	status, err := getStatus(200)
	if err != nil {
		fmt.Println("error creating status", err)
		os.Exit(1)
	}

	_, err = conn.Write([]byte(status))
	if err != nil {
		fmt.Println("error writing status")
	}
}

func getStatus(statusCode int) (string, error) {
	reason, ok := statusCodes[statusCode]
	if !ok {
		return "", fmt.Errorf("invalid status code: %d", statusCode)
	}

	statusLine := fmt.Sprintf("%s %d %s\r\n", httpVersion, statusCode, reason)
	return statusLine, nil
}

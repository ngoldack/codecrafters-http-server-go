package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strings"
)

const httpVersion = "HTTP/1.1"

var statusCodes = map[int]string{
	200: "OK",
	404: "Not Found",
}

var paths = []string{
	"/",
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

	req, err := parseRequest(conn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	status := 200
	if !slices.Contains(paths, req.target) {
		status = 404
	}

	resp, err := createStatus(status)
	if err != nil {
		fmt.Println("error creating status", err)
		os.Exit(1)
	}

	resp += "\r\n" // body

	_, err = conn.Write([]byte(resp))
	if err != nil {
		fmt.Println("error writing status")
	}
}

func createStatus(statusCode int) (string, error) {
	reason, ok := statusCodes[statusCode]
	if !ok {
		return "", fmt.Errorf("invalid status code: %d", statusCode)
	}

	statusLine := fmt.Sprintf("%s %d %s\r\n", httpVersion, statusCode, reason)
	return statusLine, nil
}

type request struct {
	httpMethod  string
	target      string
	httpVersion string

	headers map[string]string
	body    string
}

func parseRequestLine(line string) (method, target, version string, err error) {
	s := strings.Split(line, " ")

	if len(s) < 3 {
		err = fmt.Errorf("invalid request line length: %d", len(s))
		return
	}

	method = s[0]
	target = s[1]
	version = s[2]

	return
}

func parseRequest(r io.Reader) (request, error) {
	buf := bytes.NewBuffer(nil)

	_, err := r.Read(buf.Bytes())
	if err != nil {
		return request{}, err
	}

	split := strings.Split(buf.String(), "\r\n")
	if len(split) < 4 {
		return request{}, fmt.Errorf("invalid request, should consist of 3 segments. found &d", len(split))
	}
	requestLine := split[0]

	// todo: maybe pointer func
	method, target, version, err := parseRequestLine(requestLine)
	if err != nil {
		return request{}, err
	}

	// headers
	headers := make(map[string]string)
	if len(split) > 4 {
		for _, line := range split[1 : len(split)-3] {
			fmt.Println("debug header", line)
		}
	}

	// body
	body := split[len(split)-2]

	return request{
		httpMethod:  method,
		target:      target,
		httpVersion: version,
		body:        body,
		headers:     headers,
	}, nil
}

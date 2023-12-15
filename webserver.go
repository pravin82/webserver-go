package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	listenAddr := "127.0.0.1:7878"
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("Error Creating Listener")
		return
	}
	defer listener.Close()
	fmt.Println("Listening on ", listenAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting connection", err)
			continue
		}
		handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted Connection from:", conn.RemoteAddr())
	reader := bufio.NewReader(conn)

	request, err := readHTTPRequest(reader)
	if err != nil {
		fmt.Println("Error reading from Connection", err)
		return
	}
	statusLine := "HTTP/1.1 200 OK"
	filePath := "index.html"
	path := extractPath(request)
	if path != "GET / HTTP/1.1\r" {
		statusLine = "HTTP/1.1 404 NOT FOUND"
		filePath = "404.html"
	}
	fmt.Println("Path of request ", path)
	fileContent, err := readFile(filePath)
	if err != nil {
		fmt.Println("Error opening file")
	}
	fileContentLength := len(fileContent)
	response := fmt.Sprintf("%s\r\nContent-Length: %d\r\n\r\n%s", statusLine, fileContentLength, fileContent)
	conn.Write([]byte(response))

}
func extractPath(request string) string {
	lines := strings.Split(request, "\n")

	firstLine := lines[0]

	//fields := strings.Fields(firstLine)

	return firstLine
}

func readHTTPRequest(reader *bufio.Reader) (string, error) {
	var requestLines []string

	// Read lines until an empty line is encountered (end of headers)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		if line == "\r\n" || line == "\n" {
			break
		}
		requestLines = append(requestLines, line)
	}

	// Combine the request lines into a single string
	return strings.Join(requestLines, ""), nil
}

func sleep(writer http.ResponseWriter, request *http.Request) {
	duration := 15 * time.Second
	time.Sleep(duration)
	io.WriteString(writer, "Response after sleep")
}

func readFile(urlPath string) (string, error) {
	if urlPath == "" {
		urlPath = "index.html"
	}
	filePath := "www/" + urlPath
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(file)
	var content string
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	defer file.Close()
	return content, nil

}

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
)

// this function returns only User-Agent for now
func parseHeaders(requestArray []string) string {
	headers := make(map[string]string)
	for _, a := range requestArray {
		aSplit := strings.Split(a, ":")
		var a1, a2 string
		if len(aSplit) == 2 {
			a1, a2 = aSplit[0], aSplit[1]
		}
		headers[a1] = a2
	}
	return headers["User-Agent"]
}

// parse the request of the form "GET / HTTP/1.1"
func parseRequest(request string) (string, string, string) {
	parsedRequest := strings.Split(request, " ")
	var method string
	var requestTarget string
	var protocolVersion string

	if len(parsedRequest) > 2 {
		method, requestTarget, protocolVersion = parsedRequest[0], parsedRequest[1], parsedRequest[2]
	}
	return method, requestTarget, protocolVersion
}

// Handler to generate the responses based on the routes
func generateResponse(query []string) []byte {
	response := make([]byte, 1024)
	_, endPoint, _ := parseRequest(query[0])
	userAgent := strings.TrimSpace(parseHeaders(query[1:]))
	if endPoint == "/" {
		responseValue := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: \r\n\r\n"
		response = []byte(responseValue)
	} else if endPoint == "/user-agent" {
		responseValue := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
		response = []byte(responseValue)
	} else if strings.Contains(endPoint, "/echo/") {
		temp := strings.Split(endPoint, "/")
		responseBody := temp[len(temp)-1]
		responseValue := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(responseBody), responseBody)
		response = []byte(responseValue)
	} else if strings.Contains(endPoint, "/files/") {
		fmt.Println(query)

		isFilePresent, fileContents := fileHandlerRoute(endPoint)
		responseValue := ""
		if !isFilePresent {
			responseValue = "HTTP/1.1 404 Not Found\r\n\r\n"
		} else {
			responseValue = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(fileContents), fileContents)
		}
		response = []byte(responseValue)
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}
	return response
}

func getFileDetails(fileName string, directory string) (string, bool) {
	filePath := path.Join(directory, fileName)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", false
	}
	fmt.Println(fileInfo)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", false
	}
	fmt.Println(string(data))
	return string(data), true
}

func fileHandlerRoute(endPoint string) (bool, string) {
	args := os.Args
	fmt.Println(args)
	if len(args) < 3 {
		panic("args not sufficient")
	}
	directoryPath := args[2]
	temp := strings.Split(endPoint, "/")
	fileName := temp[len(temp)-1]

	fmt.Printf("searching for the file %s in directory location := %s\n", fileName, directoryPath)
	fmt.Println(endPoint)
	fileContents, present := getFileDetails(fileName, directoryPath)
	if !present {
		return false, ""
	}
	return true, fileContents

}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		localConnection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		prefix := make([]byte, 400)
		reader := bufio.NewReader(localConnection)
		_, err = reader.Read(prefix)
		lines := strings.Split(strings.TrimSpace(string(prefix)), "\r\n")
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		response := generateResponse(lines)
		_, err = localConnection.Write(response)
		localConnection.Close()
	}
}

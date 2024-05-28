package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
)

func handlePostForFiles(fileContents string, endPoint string) []byte {
	var responseValue string

	destinationFileName := strings.Split(endPoint, "/")[len(strings.Split(endPoint, "/"))-1]
	destinationDirectory := handleArgs()

	_, err := os.Create(path.Join(destinationDirectory, destinationFileName))
	if err != nil {
		// panic("unable to create file")
		responseValue = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
		return []byte(responseValue)
	}

	os.Remove(path.Join(destinationDirectory, destinationFileName))

	fileContentBytes := []byte(fileContents)
	// trim the null characters
	fileContentBytes = bytes.Trim(fileContentBytes, "\x00")

	err = os.WriteFile(path.Join(destinationDirectory, destinationFileName), fileContentBytes, 777)
	_, err = os.ReadFile(path.Join(destinationDirectory, destinationFileName))
	if err != nil {
		responseValue = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
		return []byte(responseValue)
	}
	responseValue = "HTTP/1.1 201 Created\r\n\r\n"
	return []byte(responseValue)
}

func handleGetForFiles(endPoint string) []byte {
	response := make([]byte, 1024)
	isFilePresent, fileContents := fileHandlerRoute(endPoint)
	responseValue := ""
	if !isFilePresent {
		responseValue = "HTTP/1.1 404 Not Found\r\n\r\n"
	} else {
		responseValue = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(fileContents), fileContents)
	}
	response = []byte(responseValue)
	return response
}

func isGzipPresent(headerValues string) bool {
	a := strings.Split(headerValues, ",")
	for _, h := range a {
		if strings.TrimSpace(h) == "gzip" {
			return true
		}
	}
	return false
}

// this function returns only User-Agent for now
func parseHeaders(requestArray []string) (string, map[string]string) {
	headers := make(map[string]string)
	for _, a := range requestArray {
		aSplit := strings.Split(a, ":")
		//fmt.Println(aSplit, len(aSplit))
		if aSplit[0] == "Accept-Encoding" {
			isGzipPresent(aSplit[1])
		}
		var a1, a2 string
		if len(aSplit) == 2 {
			a1, a2 = aSplit[0], aSplit[1]
		}
		headers[a1] = a2
	}
	//fmt.Println(headers)

	return headers["User-Agent"], headers
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
	method, endPoint, _ := parseRequest(query[0])
	header, _ := parseHeaders(query[1:])
	userAgent := strings.TrimSpace(header)
	if endPoint == "/" {
		responseValue := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: \r\n\r\n"
		response = []byte(responseValue)
	} else if endPoint == "/user-agent" {
		responseValue := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
		response = []byte(responseValue)
	} else if strings.Contains(endPoint, "/echo/") {
		_, headers := parseHeaders(query[1:])
		headerValue, _ := headers["Accept-Encoding"]

		var responseBody string
		var responseValue string

		temp := strings.Split(endPoint, "/")
		gzipPresent := isGzipPresent(headerValue)

		if gzipPresent {
			responseBody = temp[len(temp)-1]
			responseValue = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(responseBody), responseBody)
		} else {
			responseBody = temp[len(temp)-1]
			responseValue = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(responseBody), responseBody)
		}
		response = []byte(responseValue)
	} else if strings.Contains(endPoint, "/files/") {
		if method == "GET" {
			response = handleGetForFiles(endPoint)
		} else if method == "POST" {
			fileContents := query[len(query)-1]
			response = handlePostForFiles(fileContents, endPoint)
		} else {
		}
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}
	return response
}

func getFileDetails(fileName string, directory string) (string, bool) {
	filePath := path.Join(directory, fileName)
	_, err := os.Stat(filePath)
	if err != nil {
		return "", false
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", false
	}
	return string(data), true
}

// handleArgs returns the second argument from the args
func handleArgs() string {
	args := os.Args
	if len(args) < 3 {
		panic("args not sufficient")
	}
	return args[2]
}

func fileHandlerRoute(endPoint string) (bool, string) {
	directoryPath := handleArgs()
	temp := strings.Split(endPoint, "/")
	fileName := temp[len(temp)-1]

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

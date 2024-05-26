package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func generateResponse(query string) []byte {
	response := make([]byte, 1024)
	queryParts := strings.Split(query, " ")[1]

	if len(queryParts) > 2 {
		arr := strings.Split(queryParts, "/")
		if len(arr) > 2 {
			responseBody := arr[len(arr)-1]
			responseValue := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(responseBody), responseBody)
			response = []byte(responseValue)
		} else {
			response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
		}
	} else {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	}
	return response
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

		reader := bufio.NewReader(localConnection)
		requestLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		response := generateResponse(requestLine)
		// fmt.Println(response)
		// if len(strings.Split(requestLine, " ")[1]) > 1 {
		// 	response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
		// } else {
		// 	response = []byte("HTTP/1.1 200 OK\r\n\r\n")
		// }
		_, err = localConnection.Write(response)
		localConnection.Close()
	}

	// localConnectionTwo, err := l.Accept()
	// if err != nil {
	// 	fmt.Println("Error accepting connection: ", err.Error())
	// 	os.Exit(1)
	// }
	//
	// responseTwo := []byte("HTTP/1.1 200 OK\r\n\r\n")
	// _, err = localConnectionTwo.Write(responseTwo)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// localConnectionTwo.Close()
}

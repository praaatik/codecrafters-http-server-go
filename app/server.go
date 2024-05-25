package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

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

		response := make([]byte, 1024)
		if len(strings.Split(requestLine, " ")[1]) > 1 {
			response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
		} else {
			response = []byte("HTTP/1.1 200 OK\r\n\r\n")
		}
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

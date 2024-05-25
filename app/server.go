package main

import (
	"fmt"
	"net"
	"os"
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

	localconnection, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	response := []byte("HTTP/1.1 200 OK\r\n\r\n")
	_, err = localconnection.Write(response)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(localconnection.Write([]{'ok'}))

}

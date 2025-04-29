package main

import (
	"fmt"
	"net"
)

func handleConnection(c net.Conn) {
	defer c.Close()
	// Upon establishing connection, server will send the following message
	// to the client:
	// 0000 | version (1 byte) | player (1 bit) | grid rows (2 bytes) | grid cols (2 bytes) | current state (1 byte)

	// Inputs from the client will come in the format
	// 0000 | version (1 byte) | player (1 bit) | direction request (2 bits)
	// direction request
	// 00 -> up
	// 01 -> right
	// 10 -> left
	// 11 -> down
	for {
		packet := make([]byte, 4096)
		n, err := c.Read(packet)
		if err != nil {
			return
		}
		fmt.Printf("%s\n", packet[:n])
	}
}

func startServer() {
	l, err := net.Listen("tcp", ":6969")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(c)
	}
}

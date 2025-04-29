package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/term"
)

func cleanupFunc(oldState *term.State) func() {
	return func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Printf("\033[?25h") // restore cursor
	}
}

func startClient() {
	if !term.IsTerminal(0) || !term.IsTerminal(1) {
		log.Fatal(fmt.Errorf("stdin/stdout should be terminal"))
	}
	fmt.Printf("\033[?25l") // remove cursor

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	cleanup := cleanupFunc(oldState)
	if err != nil {
		cleanup()
		log.Fatal(err)
	}

	s, err := net.ResolveTCPAddr("tcp", "localhost:6969")
	if err != nil {
		cleanup()
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp", nil, s)
	if err != nil {
		cleanup()
		log.Fatal(err)
	}

	go func() {
		for {
			buf := make([]byte, 1)
			n, err := os.Stdin.Read(buf)
			if err != nil {
				cleanup()
				log.Fatal(err)
			}
			if n > 0 {
				switch buf[0] {
				case 3: // ^C
					cleanup()
					os.Exit(0)

				case 'a':
					_, err := conn.Write([]byte("a"))
					if err != nil {
						cleanup()
						log.Fatal(err)
					}
				}
			}
		}
	}()
	for {
	}
}

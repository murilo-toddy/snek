package main

import "flag"

var mode = flag.String("mode", "server", "server/client")

func main() {
	flag.Parse()
	switch *mode {
	case "server":
		startServer()
	case "client":
		startClient()
	}
}

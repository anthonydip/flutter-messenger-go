package main

import (
	"fmt"
	"os"

	"github.com/anthonydip/flutter-messenger-go/app/messenger-api/webserver"
)

func main() {
	port := 3333
	srv, err := webserver.New(port)

	if err != nil {
		fmt.Printf("Invalid configuration: %s\n", err)
		os.Exit(1)
	}

	srv.Start(BuildPipeline)
}

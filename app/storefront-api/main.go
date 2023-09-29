package main

import (
	"fmt"
	"os"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"
)

func main() {
	hydratedConfig := webserver.Config{
		Port: 3333,
	}

	srv, err := webserver.New(hydratedConfig)

	if err != nil {
		fmt.Printf("Invalid configuration: %s\n", err)
		os.Exit(1)
	}

	srv.Start(BuildPipeline)
}

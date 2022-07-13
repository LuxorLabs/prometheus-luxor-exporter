package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	parseCliArgs()
	if enableDebug {
		fmt.Printf("[DEBUG] Debug mode enabled.\n")
	}

	if err := runServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
}

func parseCliArgs() {
	flag.BoolVar(&enableDebug, "debug", defaultDebug, "Show debug messages.")
	flag.StringVar(&endpoint, "endpoint", defaultEndpoint, "The address-port endpoint to bind to.")

	// Exits on error
	flag.Parse()
}

package main

import (
	"fmt"
	"os"

	"github.com/minfaatong/mft-micro-cockpit/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "micro-cockpit failed: %v\n", err)
		os.Exit(1)
	}
}

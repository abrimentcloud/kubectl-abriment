package main

import (
	"fmt"
	"os"

	"github.com/abrimentcloud/kubectl-abriment/cmd"
)

var version = "dev" // This will be set during build via ldflags

func main() {
	// Check for version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("kubectl-abriment version %s\n", version)
		return
	}

	cmd.Execute()
}

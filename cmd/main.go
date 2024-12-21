package main

import (
	"fmt"

	"github.com/debug-ing/revergo/config"
)

func main() {
	fmt.Println("Starting the application...")
	// Load the configuration
	config := config.LoadConfig()
	fmt.Println("Configuration loaded:", len(config.Projects))
	//
}

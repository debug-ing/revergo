package main

import (
	"flag"
	"fmt"

	"github.com/debug-ing/revergo/config"
)

func main() {
	// Parse the flags and get address config file
	configPath := flag.String("config", "config.toml", "config file")
	flag.Parse()
	fmt.Println("Starting the application...")
	// Load the configuration
	config := config.LoadConfig(*configPath)
	fmt.Println("Configuration loaded:", len(config.Projects))
	//
}

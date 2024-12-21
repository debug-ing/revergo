package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/debug-ing/revergo/config"
	"github.com/debug-ing/revergo/internal"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Parse the flags and get address config file
	configPath := flag.String("config", "config.toml", "config file")
	flag.Parse()
	fmt.Println("Starting the application...")
	// Load the configuration
	config := config.LoadConfig(*configPath)
	// Start the reverse proxy
	reverse := internal.NewReverse(config)
	go reverse.Reverse()
	// Start gin server for monitoring...
	r := gin.Default()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	if err := r.Run(":8090"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

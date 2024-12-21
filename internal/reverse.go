package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/debug-ing/revergo/config"
)

type Reverse struct {
	config *config.AppConfig
}

func NewReverse(config *config.AppConfig) *Reverse {
	return &Reverse{config: config}
}

func (r *Reverse) Reverse() {
	for _, project := range r.config.Projects {
		//
		fmt.Println(project.Name)
		//
		listener, err := net.Listen("tcp", project.Port)
		if err != nil {
			log.Fatalf("Failed to start TCP proxy: %v", err)
		}
		defer listener.Close()
		for {

			clientConn, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept TCP connection: %v", err)
				continue
			}
			addr := clientConn.RemoteAddr()
			fmt.Println("Client connected from", addr)
			//
			// go handleTCPConnection(clientConn, project.Proxy)
			// go handleTCPConnection2(clientConn, project.Proxy)
			go handleTCPConnection(clientConn, project.Proxy)
		}
	}
}

// this function just forward the request to the proxy server
func handleTCPConnectionWithIoCopy(clientConn net.Conn, proxy string) {
	defer clientConn.Close()
	targetConn, err := net.Dial("tcp", proxy)
	if err != nil {
		log.Printf("Failed to connect to target server: %v", err)
		return
	}
	defer targetConn.Close()
	go io.Copy(targetConn, clientConn)
	io.Copy(clientConn, targetConn)
}

// this function parse the request and response
func handleTCPConnection2(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()
	buffer := make([]byte, 4096)
	n, err := clientConn.Read(buffer)
	if err != nil {
		log.Printf("Error reading client data: %v", err)
		return
	}
	requestLine := string(buffer[:n])
	method, url, _ := parseHTTPRequestLine(requestLine)
	fmt.Printf("HTTP Method: %s, URL: %s\n", method, url)
	httpRequest, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buffer[:n])))
	if err != nil {
		log.Printf("Failed to parse HTTP request: %v", err)
		return
	}
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Failed to connect to target server: %v", err)
		return
	}
	defer targetConn.Close()
	err = httpRequest.Write(targetConn)
	if err != nil {
		log.Printf("Failed to forward request to target server: %v", err)
		return
	}
	responseBuffer := make([]byte, 4096)
	n, err = targetConn.Read(responseBuffer)
	if err != nil {
		log.Printf("Failed to read response from target server: %v", err)
		return
	}
	statusLine := string(responseBuffer[:n])
	statusCode := extractHTTPStatusCode(statusLine)
	fmt.Printf("Response Status Code: %d\n", statusCode)
	clientConn.Write(responseBuffer[:n])
}

func handleTCPConnection(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	// Read client request
	buffer := make([]byte, 4096)
	n, err := clientConn.Read(buffer)
	if err != nil {
		log.Printf("Error reading client data: %v", err)
		return
	}

	// Parse HTTP request
	httpRequest, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buffer[:n])))
	if err != nil {
		log.Printf("Failed to parse HTTP request: %v", err)
		return
	}

	// Set the target address as the Host
	httpRequest.Host = targetAddr
	httpRequest.RequestURI = "" // Clear the RequestURI for forwarding

	// Connect to the target server
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Failed to connect to target server: %v", err)
		return
	}
	defer targetConn.Close()

	// Forward the HTTP request to the target server
	err = httpRequest.Write(targetConn)
	if err != nil {
		log.Printf("Failed to forward request to target server: %v", err)
		return
	}

	// Read the response from the target server
	targetReader := bufio.NewReader(targetConn)
	httpResponse, err := http.ReadResponse(targetReader, httpRequest)
	if err != nil {
		log.Printf("Failed to read response from target server: %v", err)
		return
	}
	defer httpResponse.Body.Close()

	// Write the response back to the client
	err = httpResponse.Write(clientConn)
	if err != nil {
		log.Printf("Failed to forward response to client: %v", err)
		return
	}

	fmt.Printf("Forwarded request to %s, response status: %d\n", targetAddr, httpResponse.StatusCode)
}

func parseHTTPRequestLine(requestLine string) (method, url, version string) {
	parts := strings.Fields(requestLine)
	if len(parts) >= 3 {
		return parts[0], parts[1], parts[2]
	}
	return "", "", ""
}

func extractHTTPStatusCode(statusLine string) int {
	parts := strings.Fields(statusLine)
	if len(parts) >= 2 {
		// Assuming the status code is the second element
		var statusCode int
		_, err := fmt.Sscanf(parts[1], "%d", &statusCode)
		if err == nil {
			return statusCode
		}
	}
	return 0
}

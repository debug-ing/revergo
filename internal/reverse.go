package internal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
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
			go handleConnection(clientConn, project.Proxy)
		}
	}
}
func handleConnection(clientConn net.Conn, port string) {
	defer clientConn.Close()
	backendConn, err := net.Dial("tcp", port)
	if err != nil {
		return
	}
	defer backendConn.Close()
	clientReader := bufio.NewReader(clientConn)
	backendWriter := bufio.NewWriter(backendConn)
	for {
		line, err := clientReader.ReadString('\n')
		if err != nil {
			return
		}
		if line == "\r\n" {
			backendWriter.WriteString(line)
			backendWriter.Flush()
			break
		}
		if strings.HasPrefix(line, "Host:") {
			line = fmt.Sprintf("Host: %s\r\n", port)
		}

		backendWriter.WriteString(line)
		backendWriter.Flush()
	}
	go io.Copy(backendConn, clientConn)
	io.Copy(clientConn, backendConn)
}

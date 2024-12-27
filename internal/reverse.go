package internal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/debug-ing/revergo/config"
	"github.com/debug-ing/revergo/pkg/logger"
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
			go r.handleConnectionDetail(clientConn, project.Proxy, project.Domain[0])
		}
	}
}

// handleConnection this function reverse proxy with out log
func (r *Reverse) handleConnection(clientConn net.Conn, port string) {
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

// handleConnectionDetail this function reverse proxy with detail
func (r *Reverse) handleConnectionDetail(clientConn net.Conn, port, allowedDomain string) {
	defer clientConn.Close()
	backendConn, err := net.Dial("tcp", port)
	if err != nil {
		log.Printf("Failed to connect to backend: %v", err)
		return
	}
	defer backendConn.Close()
	clientReader := bufio.NewReader(clientConn)
	clientWriter := bufio.NewWriter(clientConn)
	backendReader := bufio.NewReader(backendConn)
	backendWriter := bufio.NewWriter(backendConn)
	req, err := http.ReadRequest(clientReader)
	if err != nil {
		log.Printf("Failed to read HTTP request: %v", err)
		return
	}
	///move to function
	// if !strings.HasSuffix(req.Host, allowedDomain) {
	// 	log.Printf("Connection rejected: Host %s is not allowed", req.Host)
	// 	clientWriter.WriteString("HTTP/1.1 403 Forbidden\r\n")
	// 	clientWriter.WriteString("Content-Type: text/plain\r\n")
	// 	clientWriter.WriteString("Connection: close\r\n")
	// 	clientWriter.WriteString("\r\n")
	// 	clientWriter.WriteString("Access denied: invalid domain\r\n")
	// 	clientWriter.Flush()
	// 	return
	// }
	//if !r.checkHost(clientWriter, req.Host, allowedDomain) {
	//	return
	//}
	log.Printf("Request: Method=%s, URL=%s", req.Method, req.URL)
	err = req.Write(backendWriter)
	if err != nil {
		log.Printf("Failed to forward HTTP request: %v", err)
		return
	}
	backendWriter.Flush()
	resp, err := http.ReadResponse(backendReader, req)
	if err != nil {
		log.Printf("Failed to read HTTP response: %v", err)
		return
	}
	log.Printf("Response: StatusCode=%d", resp.StatusCode)
	r.addLog(*resp, *req)
	err = resp.Write(clientWriter)
	if err != nil {
		log.Printf("Failed to forward HTTP response: %v", err)
		return
	}
	clientWriter.Flush()
}

// checkHost this function check host
func (r *Reverse) checkHost(clientWriter *bufio.Writer, host string, allowedDomain string) bool {
	if !strings.HasSuffix(host, allowedDomain) {
		log.Printf("Connection rejected: Host %s is not allowed", host)
		clientWriter.WriteString("HTTP/1.1 403 Forbidden\r\n")
		clientWriter.WriteString("Content-Type: text/plain\r\n")
		clientWriter.WriteString("Connection: close\r\n")
		clientWriter.WriteString("\r\n")
		clientWriter.WriteString("Access denied: invalid domain\r\n")
		clientWriter.Flush()
		return false
	}
	return true
}

// addLog this function add log file
func (r *Reverse) addLog(resp http.Response, req http.Request) {
	itemLog := map[string]interface{}{}
	itemLog["status"] = resp.StatusCode
	itemLog["url"] = req.URL.Path
	itemLog["method"] = req.Method
	logger.Info(itemLog)
}

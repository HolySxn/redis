package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	port = flag.Int("port", 8080, "Port number")
	help = flag.Bool("help", false, "Show help screen")
)

type Server struct {
	listenAddr string
	quitch     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
	}
}

func (s *Server) StartServer() error {
	addr, err := net.ResolveUDPAddr("udp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("error resolving address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("error listening on UDP: %w", err)
	}
	defer conn.Close()

	fmt.Println("connection: ", conn.LocalAddr().String())

	go s.readLoop(conn)

	<-s.quitch

	return nil
}

func (s *Server) readLoop(conn *net.UDPConn) {
	buf := make([]byte, 2048)
	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		fmt.Print(string(buf[:n]))
	}
}

func main() {
	flag.Parse()
	if !isValidPort(*port) {
		log.Fatal("Invalid port", *port)
	}

	server := NewServer(fmt.Sprintf(":%d", *port))
	server.StartServer()

}

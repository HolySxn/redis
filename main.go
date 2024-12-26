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
	ln         net.PacketConn
	quitch     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
	}
}

func (s *Server) StartServer() error {
	ln, err := net.ListenPacket("udp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	<-s.quitch

	return nil
}

func main() {
	flag.Parse()
	if !isValidPort(*port) {
		log.Fatal("Invalid port", *port)
	}

	server := NewServer(fmt.Sprintf(":%d", *port))
	server.StartServer()

}

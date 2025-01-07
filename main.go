package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	port = flag.Int("port", 8080, "Port number")
	help = flag.Bool("help", false, "Show help screen")
)

const maxMessageSize = 1024

type Storage struct {
	data map[string]string
	px   map[string]time.Time
	mu   sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
		px:   make(map[string]time.Time),
	}
}

func StartServer(listenAddr string, storage *Storage) error {
	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return fmt.Errorf("error resolving address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("error listening on UDP: %w", err)
	}
	defer conn.Close()

	fmt.Println("connection: ", conn.LocalAddr().String())

	readLoop(conn, storage)

	return nil

}

func readLoop(conn *net.UDPConn, store *Storage) {
	buf := make([]byte, maxMessageSize)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		if n == maxMessageSize && buf[maxMessageSize-1] != '\n' {
			conn.WriteToUDP([]byte("(error) message too large(>1024 bytes)"), remoteAddr)
			continue
		}

		go Handle(buf[:n], conn, remoteAddr, store)
	}
}

func main() {
	flag.Parse()
	if *help {
		hlp()
		return
	}

	if !isValidPort(*port) {
		log.Fatal("Invalid port", *port)
	}

	store := NewStorage()

	err := StartServer(fmt.Sprintf(":%d", *port), store)
	if err != nil {
		log.Fatal("error to start server: ", err)
	}

}

func hlp() {
	msg := `Own Redis

Usage:
  own-redis [--port <N>]
  own-redis --help

Options:
  --help       Show this screen.
  --port N     Port number.

Commands:
  PING                Check if the server is running.
  SET [key] [value]   Store a value with a key.
  SET [key] [value] px [ms]
                      Store a value with a key and expiration in milliseconds.
  GET [key]           Retrieve the value by key.
`
	fmt.Println(msg)
}

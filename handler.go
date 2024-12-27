package main

import (
	"fmt"
	"net"
	"strings"
)

func Handle(input []byte, conn *net.UDPConn, remoteAddr *net.UDPAddr, store *Storage) {
	s := PrepareInput(input)
	fmt.Println(s)

	switch strings.ToLower(s[0]) {
	case "ping":
		conn.WriteToUDP([]byte("PONG\n"), remoteAddr)
	case "set":

	}
}

func set() {

}

func get() {

}

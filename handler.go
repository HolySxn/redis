package main

import (
	"net"
	"strconv"
	"strings"
	"time"
)

func Handle(input []byte, conn *net.UDPConn, remoteAddr *net.UDPAddr, store *Storage) {
	arguments := PrepareInput(input)

	if len(arguments) == 0 {
		conn.WriteToUDP([]byte("(error) ERR invalid input\n"), remoteAddr)
		return
	}

	switch strings.ToLower(arguments[0]) {
	case "ping":
		conn.WriteToUDP([]byte("PONG\n"), remoteAddr)
	case "set": // SET foo bar px 10000
		if len(arguments) < 3 {
			conn.WriteToUDP([]byte("(error) ERR wrong number of arguments for 'SET' command\n"), remoteAddr)
			return
		}

		duration := time.Duration(0)
		key := arguments[1]
		values := strings.Join(arguments[2:], " ")
		if len(arguments) > 4 && strings.ToLower(arguments[len(arguments)-2]) == "px" {
			ms, err := strconv.Atoi(arguments[len(arguments)-1])
			if err != nil || ms < 0 {
				conn.WriteToUDP([]byte("(error) ERR invalid px argument\n"), remoteAddr)
				return
			}
			duration = time.Duration(ms) * time.Millisecond
			values = strings.Join(arguments[2:len(arguments)-2], " ")
		}
		store.set(key, values, duration)
		conn.WriteToUDP([]byte("OK\n"), remoteAddr)
	case "get":
		if len(arguments) != 2 {
			conn.WriteToUDP([]byte("(error) ERR wrong number of arguments for 'GET' command\n"), remoteAddr)
			return
		}

		key := arguments[1]
		value, exists := store.get(key)
		if exists {
			conn.WriteToUDP([]byte(value+"\n"), remoteAddr)
		} else {
			conn.WriteToUDP([]byte("(nil)\n"), remoteAddr)
		}

	default:
		conn.WriteToUDP([]byte("(error) ERR wrong command\n"), remoteAddr)
	}
}

func (s *Storage) set(key, value string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	if duration > 0 {
		s.px[key] = time.Now().Add(duration)
	} else {
		delete(s.px, key)
	}
}

func (s *Storage) get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if expiry, exists := s.px[key]; exists {
		if time.Now().After(expiry) {
			delete(s.data, key)
			delete(s.px, key)
		}
	}

	value, exists := s.data[key]
	return value, exists
}

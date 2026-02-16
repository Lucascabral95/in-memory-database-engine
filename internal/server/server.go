package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/lucas-dev/in-memory-db/internal/storage"
)

type Server struct {
	store *storage.MemoryStore
	addr  string
}

func NewServer(store *storage.MemoryStore, addr string) *Server {
	if addr == "" {
		addr = ":6379"
	}

	return &Server{
		store: store,
		addr:  addr,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("no se pudo iniciar Redis TCP en %s: %w", s.addr, err)
	}

	fmt.Printf("Redis TCP escuchando en %s\n", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSuffix(line, "\r\n")

		if strings.HasPrefix(line, "*") {
			count, err := strconv.Atoi(line[1:])
			if err != nil || count <= 0 {
				if _, writeErr := conn.Write([]byte("-ERR invalid RESP array\r\n")); writeErr != nil {
					return
				}
				continue
			}

			args := make([]string, 0, count)

			for i := 0; i < count; i++ {
				lenLine, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				lenLine = strings.TrimSuffix(lenLine, "\r\n")

				length, err := strconv.Atoi(strings.TrimPrefix(lenLine, "$"))
				if err != nil || length < 0 {
					if _, writeErr := conn.Write([]byte("-ERR invalid bulk length\r\n")); writeErr != nil {
						return
					}
					return
				}

				val := make([]byte, length+2)
				if _, err = io.ReadFull(reader, val); err != nil {
					return
				}
				args = append(args, string(val[:length]))
			}

			response := s.processCommand(args)
			if _, err := conn.Write([]byte(response)); err != nil {
				return
			}
		}
	}
}

func (s *Server) processCommand(args []string) string {
	if len(args) == 0 {
		return "-ERR empty command\r\n"
	}

	cmd := strings.ToUpper(args[0])

	switch cmd {
	case "PING":
		return "+PONG\r\n"

	case "SET":
		if len(args) != 3 && len(args) != 5 {
			return "-ERR wrong number of arguments\r\n"
		}

		key := args[1]
		val := args[2]
		ttl := 0

		if len(args) == 5 {
			if strings.ToUpper(args[3]) != "EX" {
				return "-ERR syntax error\r\n"
			}

			parsedTTL, err := strconv.Atoi(args[4])
			if err != nil || parsedTTL <= 0 {
				return "-ERR value is not an integer or out of range\r\n"
			}

			ttl = parsedTTL
		}

		s.store.Set(key, []byte(val), ttl)
		return "+OK\r\n"

	case "GET":
		if len(args) < 2 {
			return "-ERR wrong number of arguments\r\n"
		}

		val, ok := s.store.Get(args[1])
		if !ok {
			return "$-1\r\n"
		}

		return fmt.Sprintf("$%d\r\n%s\r\n", len(val), string(val))

	case "DEL":
		if len(args) < 2 {
			return "-ERR wrong number of arguments\r\n"
		}

		s.store.Del(args[1])
		return ":1\r\n"

	default:
		return "-ERR unknown command '" + cmd + "'\r\n"
	}
}

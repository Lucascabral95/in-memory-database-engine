package server

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lucas-dev/in-memory-db/internal/storage"
)

func TestProcessCommand(t *testing.T) {
	store := storage.NewMemoryStore()
	s := NewServer(store, ":0")

	tests := []struct {
		name     string
		args     []string
		expected string
		setup    func()
	}{
		{"PING command", []string{"PING"}, "+PONG\r\n", nil},

		{"SET basic", []string{"SET", "k1", "v1"}, "+OK\r\n", nil},

		{"GET existing", []string{"GET", "k1"}, "$2\r\nv1\r\n", nil},

		{"GET non-existing", []string{"GET", "nope"}, "$-1\r\n", nil},

		{"SET with TTL", []string{"SET", "k_ttl", "v_ttl", "EX", "10"}, "+OK\r\n", nil},

		{"DEL existing", []string{"DEL", "k1"}, ":1\r\n", nil},

		{"SET wrong args", []string{"SET", "key"}, "-ERR wrong number of arguments\r\n", nil},
		{"SET invalid TTL", []string{"SET", "k", "v", "EX", "-5"}, "-ERR value is not an integer or out of range\r\n", nil},
		{"Unknown command", []string{"FOO"}, "-ERR unknown command 'FOO'\r\n", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := s.processCommand(tt.args)
			if resp != tt.expected {
				t.Errorf("Expected: %q, Got: %q", tt.expected, resp)
			}
		})
	}
}

func readFullResponse(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(line, "$") {
		if strings.HasPrefix(line, "$-1") {
			return line, nil
		}

		lengthStr := strings.TrimSuffix(strings.TrimPrefix(line, "$"), "\r\n")
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}

		body := make([]byte, length+2)
		if _, err := io.ReadFull(reader, body); err != nil {
			return "", err
		}
		return line + string(body), nil
	}

	return line, nil
}

func TestHandleConnection_Integration(t *testing.T) {
	store := storage.NewMemoryStore()
	s := NewServer(store, ":0")

	client, server := net.Pipe()
	defer client.Close()

	go s.handleConnection(server)

	reader := bufio.NewReader(client)

	t.Run("PING", func(t *testing.T) {
		_, err := client.Write([]byte("*1\r\n$4\r\nPING\r\n"))
		if err != nil {
			t.Fatal(err)
		}

		resp, err := readFullResponse(reader)
		if err != nil {
			t.Fatal(err)
		}
		if resp != "+PONG\r\n" {
			t.Errorf("Expected '+PONG\\r\\n', got %q", resp)
		}
	})

	t.Run("SET and GET", func(t *testing.T) {
		_, err := client.Write([]byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"))
		if err != nil {
			t.Fatal(err)
		}
		resp, _ := readFullResponse(reader)
		if !strings.Contains(resp, "+OK") {
			t.Errorf("Expected OK, got %q", resp)
		}

		_, err = client.Write([]byte("*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"))
		if err != nil {
			t.Fatal(err)
		}
		resp, _ = readFullResponse(reader)
		expected := "$5\r\nvalue\r\n"
		if resp != expected {
			t.Errorf("Expected %q, got %q", expected, resp)
		}
	})
}

func TestTTL_Integration(t *testing.T) {
	store := storage.NewMemoryStore()
	s := NewServer(store, ":0")

	client, server := net.Pipe()
	defer client.Close()

	go s.handleConnection(server)
	reader := bufio.NewReader(client)

	req := "*5\r\n$3\r\nSET\r\n$4\r\nkey2\r\n$4\r\ndata\r\n$2\r\nEX\r\n$1\r\n1\r\n"
	_, err := client.Write([]byte(req))
	if err != nil {
		t.Fatal("Error writing SET:", err)
	}

	resp, err := readFullResponse(reader)
	if err != nil || !strings.Contains(resp, "+OK") {
		t.Fatalf("Expected OK after SET, got %q, err %v", resp, err)
	}

	_, err = client.Write([]byte("*2\r\n$3\r\nGET\r\n$4\r\nkey2\r\n"))
	if err != nil {
		t.Fatal("Error writing GET 1:", err)
	}

	resp, err = readFullResponse(reader)
	if err != nil {
		t.Fatal("Error reading GET 1:", err)
	}

	if !strings.Contains(resp, "data") {
		t.Fatalf("Expected 'data' immediately after SET, got %q", resp)
	}

	time.Sleep(1100 * time.Millisecond)

	_, err = client.Write([]byte("*2\r\n$3\r\nGET\r\n$4\r\nkey2\r\n"))
	if err != nil {
		t.Fatal("Error writing GET 2:", err)
	}

	resp, err = readFullResponse(reader)
	if err != nil {
		t.Fatal("Error reading GET 2:", err)
	}

	if resp != "$-1\r\n" {
		t.Errorf("Expected expiration ($-1), got %q", resp)
	}
}

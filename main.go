package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("kewl-db listening on :4000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	log.Printf("client connected: %s", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	conn.Write([]byte("welcome to kewl-db\n> "))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			conn.Write([]byte("> "))
			continue
		}

		response := execute(line)
		conn.Write([]byte(response + "\n> "))
	}

	log.Printf("client disconnected: %s", conn.RemoteAddr())
}

func execute(input string) string {
	upper := strings.ToUpper(input)

	switch {
	case upper == "PING":
		return "PONG"
	case upper == "QUIT":
		return "bye"
	case strings.HasPrefix(upper, "SELECT"):
		return fmt.Sprintf("received query: %s (not implemented yet)", input)
	case strings.HasPrefix(upper, "CREATE"):
		return fmt.Sprintf("received DDL: %s (not implemented yet)", input)
	case strings.HasPrefix(upper, "INSERT"):
		return fmt.Sprintf("received insert: %s (not implemented yet)", input)
	default:
		return fmt.Sprintf("unknown command: %s", input)
	}
}

package cli

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func startTCP() {
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
	log.Printf("client connected :%s", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	conn.Write([]byte("welcome to kewl-db\n> "))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			conn.Write([]byte("> "))
			continue
		} else if line == "exit" { // disconnect on exit.
			conn.Write([]byte("exiting..."))
			break
		}
		conn.Write([]byte("got: "))
		conn.Write([]byte(line + "\n> "))

	}
	log.Printf("client disconnected: %s", conn.RemoteAddr())
}

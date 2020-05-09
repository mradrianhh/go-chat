package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/mradrianhh/go-chat/internal/pkg/info"
	"github.com/mradrianhh/go-chat/internal/pkg/models"
)

var conns map[string]net.Conn

func main() {
	conns = make(map[string]net.Conn)

	service := "0.0.0.0:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		if err := handleConnection(conn); err == nil {
			go listen(conn)
		} else {
			continue
		}
	}
}

func handleConnection(conn net.Conn) error {
	var buffer [512]byte
	n, err := conn.Read(buffer[0:])
	if err != nil {
		conn.Write([]byte(info.UNAUTHENTICATED))
		return err
	}

	username := string(buffer[0:n])
	conns[username] = conn
	fmt.Printf("Added new conn(%s) with username %s\n", conn.RemoteAddr().String(), username)
	conn.Write([]byte(info.AUTHENTICATED))
	return nil
}

func listen(conn net.Conn) {
	defer conn.Close()

	for {
		decoder := gob.NewDecoder(conn)

		var message models.Message
		decoder.Decode(&message)

		if message.Text == "exit" {
			conns[message.User].Close()
			delete(conns, message.User)
		}

		broadcast(message.User + ": " + message.Text)
	}
}

func broadcast(message string) {
	for _, conn := range conns {
		conn.Write([]byte(message))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/mradrianhh/go-chat/internal/pkg/info"
	"github.com/mradrianhh/go-chat/internal/pkg/models"
)

var username string
var state string

func main() {
	state = info.UNAUTHENTICATED

	conn, err := net.Dial("tcp", "0.0.0.0:1200")
	checkError(err)

	for {
		if result := authenticate(conn); result == info.AUTHENTICATED {
			go listen(conn)
			write(conn)
		} else {
			fmt.Println("Result: " + result)
		}
	}
}

func authenticate(conn net.Conn) string {
	fmt.Print("Enter username: ")
	if _, err := fmt.Scan(&username); err == nil {
		conn.Write([]byte(username))
	}
	var buffer [512]byte
	n, err := conn.Read(buffer[0:])
	if err != nil {
		return info.UNAUTHENTICATED
	}

	result := string(buffer[0:n])
	fmt.Println(result)
	return result
}

func listen(conn net.Conn) {
	var buffer [512]byte
	for {
		n, err := conn.Read(buffer[0:])
		if err != nil {
			continue
		}

		text := string(buffer[0:n])
		fmt.Println(text)
	}
}

func write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimRight(line, "\t\r\n")
		if err != nil {
			break
		}

		encoder := gob.NewEncoder(conn)

		if s := strings.ToLower(line); s == "exit" {
			fmt.Println("Closing connection...")
			message := models.NewMessage(username, s)
			encoder.Encode(message)
			conn.Close()
			time.Sleep(1 * time.Second)
			fmt.Println("Exiting...")
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}

		if s := strings.TrimSpace(line); s != "" {
			message := models.NewMessage(username, line)
			encoder.Encode(message)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	clients = make(map[net.Conn]string)
	mutex   sync.Mutex
)

// 向所有客户端广播消息
func broadcastMessage(msg string, sender net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	senderName, ok := clients[sender]
	if !ok {
		fmt.Println("Sender not found in client list")
		return
	}

	for client := range clients {
		if client != sender {
			message := fmt.Sprintf("[%s]: %s", senderName, msg)
			client.Write([]byte(message))
		}
	}
}

// 处理每个客户端的连接
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 首先获取用户的名字
	conn.Write([]byte("Please enter your name:\n"))
	nameReader := bufio.NewReader(conn)
	name, err := nameReader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read name:", err)
		return
	}
	name = strings.TrimSpace(name)

	mutex.Lock()
	clients[conn] = name
	mutex.Unlock()

	fmt.Printf("%s has joined the chat\n", name)
	broadcastMessage("has joined the chat\n", conn)

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		broadcastMessage(msg, conn)
	}

	mutex.Lock()
	fmt.Printf("%s has left the chat\n", name)
	broadcastMessage("has left the chat\n", conn)
	delete(clients, conn)
	mutex.Unlock()
	
}

func main() {
	fmt.Println("Chat server is starting...")

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Chat server started. Waiting for clients...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

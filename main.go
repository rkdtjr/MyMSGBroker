package main

import (
	"embed"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

//go:embed web/index.html
var content embed.FS

type Broker struct {
	subscribers map[string]map[net.Conn]bool
	mu          sync.RWMutex
	monitor     chan string 
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string]map[net.Conn]bool),
		monitor:     make(chan string, 100),
	}
}

func (b *Broker) Subscribe(topic string, conn net.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.subscribers[topic] == nil {
		b.subscribers[topic] = make(map[net.Conn]bool)
	}
	b.subscribers[topic][conn] = true
	fmt.Printf("[SUB] %s -> %s\n", conn.RemoteAddr(), topic)
}

func (b *Broker) Unsubscribe(conn net.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for topic := range b.subscribers {
		delete(b.subscribers[topic], conn)
	}
}

func (b *Broker) Publish(topic, payload string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	msg := fmt.Sprintf("[%s] %s", topic, payload)
	
	select {
	case b.monitor <- msg:
	default:
	}

	if clients, ok := b.subscribers[topic]; ok {
		for conn := range clients {
			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil { continue }
		}
	}
}

func (b *Broker) handleTCPClient(conn net.Conn) {
	defer func() {
		b.Unsubscribe(conn)
		conn.Close()
	}()

	header := make([]byte, 6)
	for {
		if _, err := io.ReadFull(conn, header); err != nil { return }

		action := header[0]
		tLen := int(header[1])
		pLen := binary.BigEndian.Uint32(header[2:6])

		body := make([]byte, tLen+int(pLen))
		if _, err := io.ReadFull(conn, body); err != nil { return }

		topic := string(body[:tLen])
		payload := string(body[tLen:])

		switch action {
		case 0x01: b.Subscribe(topic, conn)
		case 0x02: b.Publish(topic, payload)
		}
	}
}

func main() {
	broker := NewBroker()

	go func() {
		ln, _ := net.Listen("tcp", ":8080")
		fmt.Println("✅ TCP Broker: 8080")
		for {
			conn, _ := ln.Accept()
			go broker.handleTCPClient(conn)
		}
	}()

	http.Handle("/", http.FileServer(http.FS(content)))
	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		for msg := range broker.monitor {
			if err := websocket.Message.Send(ws, msg); err != nil { break }
		}
	}))

	fmt.Println("🚀 Dashboard: http://localhost")
	http.ListenAndServe(":80", nil)
}
package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type client chan<- string

var (
	entering   = make(chan client)
	leaving    = make(chan client)
	messages   = make(chan string)
	inMessages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

//func handleConn(c net.Conn) {
//	defer c.Close()
//
//	for {
//		_, err := io.WriteString(c, time.Now().Format("15:04:05\n\r"))
//		if err != nil {
//			return
//		}
//		time.Sleep(1 * time.Second)
//	}
//}

func handleConn(conn net.Conn) {
	var mu sync.Mutex
	ch := make(chan string)
	go clientWriter(conn, ch)
	for {
		go func() {
			var text string
			_, err := fmt.Scanf("%s\n", &text)
			if err != nil {
				return
			}
			inMessages <- text
		}()
		select {
		case msg := <-inMessages:
			mu.Lock()
			ch <- msg
			mu.Unlock()
		default:
			mu.Lock()
			ch <- time.Now().Format("15:04:05\r")
			mu.Unlock()
			time.Sleep(1 * time.Second)
		}
	}
	//leaving <- ch
	//messages <- who + " has left"
	//conn.Close()
}
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			continue
		}
	}
}

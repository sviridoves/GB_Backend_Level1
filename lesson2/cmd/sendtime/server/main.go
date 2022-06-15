package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
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

func handleConn(conn net.Conn) {
	var exitGo = false
	ch := make(chan string)
	go clientWriter(conn, ch)
	who := conn.RemoteAddr().String()
	fmt.Println("open: ", who)
	ch <- "Start session\n"
	entering <- ch
	go func() {
		input := bufio.NewScanner(conn)
		for !input.Scan() {
			leaving <- ch
			fmt.Println("close: ", who)
			conn.Close()
			exitGo = true
			return
		}
	}()
	for {
		go func() {
			var text string
			srvInput := bufio.NewScanner(os.Stdin)
			srvInput.Scan()
			text = srvInput.Text()
			messages <- fmt.Sprintln(text)
		}()
		ch <- time.Now().Format("15:04:05\n\r")
		time.Sleep(2 * time.Second)
		if exitGo {
			break
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprint(conn, msg)
		if err != nil {
			continue
		}
	}
}

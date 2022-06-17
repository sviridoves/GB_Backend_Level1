package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			log.Fatal(err)
		}
	}()
	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Fatal(err)
	}
}

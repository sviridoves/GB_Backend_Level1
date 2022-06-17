package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	myUserName := flag.String("name", "", "Имя пользователя")
	flag.Parse()
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	if *myUserName != "" {
		if _, err := io.WriteString(conn, fmt.Sprintf("myName:%s\n", *myUserName)); err != nil {
			log.Fatal(err)
		}
	} else {
		if _, err := io.WriteString(conn, fmt.Sprintf("myName:%s\n", conn.LocalAddr())); err != nil {
			log.Fatal(err)
		}
	}
	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			log.Fatal(err)
		}
	}()
	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Fatal(err)
	}
	if *myUserName != "" {
		fmt.Printf("%s: exit", *myUserName)
	} else {
		fmt.Printf("%s: exit", conn.LocalAddr())
	}
}

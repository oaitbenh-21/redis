package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	wg := sync.WaitGroup{}
	for {
		connection, err := listener.Accept()
		if err == nil {
			wg.Add(1)
			go HandleConnection(connection, wg)
		}

	}
}

func HandleConnection(conn net.Conn, wg sync.WaitGroup) {
	defer wg.Done()
	for {
		buff := make([]byte, 1024)
		count, _ := conn.Read(buff)
		fmt.Print(string(buff[:count]))

		result := string(buff[:count])
		fmt.Printf("%#v\n", result)
	}
}

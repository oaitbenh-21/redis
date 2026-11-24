package main

import (
	"fmt"
	"net"
	"redis/internal/server"
	"redis/internal/storage"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	wg := sync.WaitGroup{}
	s := storage.NewStore()
	for {
		connection, err := listener.Accept()
		if err == nil {
			wg.Add(1)
			go server.HandleConnection(connection, &wg, s)
		}

	}
}

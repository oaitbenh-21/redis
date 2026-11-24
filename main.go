package main

import (
	"fmt"
	"net"
	"redis/internal/engine"
	"redis/internal/storage"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
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
			go HandleConnection(connection, &wg, s)
		}

	}
}

func HandleConnection(conn net.Conn, wg *sync.WaitGroup, s *storage.Store) {
	defer wg.Done()
	for {
		buff := make([]byte, 1024)
		count, _ := conn.Read(buff)
		result := string(buff[:count])
		cmd, err := engine.Parse(result)
		if err != nil {
			fmt.Printf("%#v\n", result)
			fmt.Println(err)
			return
		}
		fmt.Println("command : ", cmd)
		switch cmd.Name {
		case "SET":
			s.Set(cmd.Args[0], cmd.Args[1])
			conn.Write([]byte("ok\n"))
		case "GET":
			value, _ := s.Get(cmd.Args[0])
			conn.Write([]byte(value + "\n"))
		}

	}
}

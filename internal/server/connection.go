package server

import (
	"net"
	"redis/internal/engine"
	"redis/internal/storage"
	"redis/resp"
	"strings"
	"sync"
)

func HandleConnection(conn net.Conn, wg *sync.WaitGroup, s *storage.Store) {
	defer wg.Done()
	for {
		buff := make([]byte, 1024)
		count, err := conn.Read(buff)
		if err != nil {
			return
		}
		data := string(buff[:count])
		result, _ := resp.Parse(strings.NewReader(data))
		cmd, _ := engine.Parse(result)
		switch cmd.Name {
		case "PING":
			conn.Write([]byte(resp.SimpleString("PONG")))
		case "SET":
			s.Set(cmd.Args[0], cmd.Args[1])
			conn.Write([]byte(resp.SimpleString("OK")))
		case "GET":
			value, _ := s.Get(cmd.Args[0])
			conn.Write([]byte(resp.SimpleString(value)))
		case "DEL":
			ok := s.Delete(cmd.Args[0])
			if ok {
				conn.Write([]byte(resp.SimpleString("true")))
			} else {
				conn.Write([]byte(resp.SimpleString("false")))
			}
		default:
			conn.Write([]byte(resp.SimpleString("Err-- Wrong command \"" + cmd.Name + "\"")))
		}
	}
}

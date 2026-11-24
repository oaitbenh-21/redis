package resp

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Parse(conn io.Reader) (string, error) {
	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if len(line) == 0 || line[0] != '*' {
		return "", fmt.Errorf("expected RESP array")
	}

	var count int
	fmt.Sscanf(strings.TrimSpace(line[1:]), "%d", &count)

	args := make([]string, 0, count)

	for i := 0; i < count; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		if len(line) == 0 || line[0] != '$' {
			return "", fmt.Errorf("expected bulk string")
		}

		var length int
		fmt.Sscanf(strings.TrimSpace(line[1:]), "%d", &length)

		if length < 0 {
			return "", fmt.Errorf("null bulk string")
		}

		data := make([]byte, length+2)

		if _, err := io.ReadFull(reader, data); err != nil {
			return "", err
		}

		if string(data[length:]) != "\r\n" {
			return "", fmt.Errorf("invalid terminator")
		}

		args = append(args, string(data[:length]))
	}

	var result strings.Builder

	for i, arg := range args {
		if i > 0 {
			result.WriteByte(' ')
		}

		if i == 0 || i == 1 {
			result.WriteString(strings.ToLower(arg))
		} else {
			fmt.Fprintf(&result, "%q", arg)
		}
	}

	return result.String(), nil
}

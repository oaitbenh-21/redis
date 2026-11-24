package engine

import (
	"fmt"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func Parse(input string) (Command, error) {
	args, err := tokenize(input)
	if err != nil {
		fmt.Println(err)
		return Command{}, err
	}

	if len(args) == 0 {
		return Command{}, fmt.Errorf("empty command")
	}

	name := strings.ToUpper(args[0])
	args = args[1:]

	switch name {
	case "GET", "DEL", "EXISTS":
		if len(args) != 1 {
			return Command{}, fmt.Errorf("%s expects 1 argument", name)
		}

	case "SET":
		if len(args) != 2 {
			return Command{}, fmt.Errorf("SET expects 2 arguments")
		}

	case "KEYS":
		if len(args) != 1 {
			return Command{}, fmt.Errorf("KEYS expects 1 argument")
		}

	case "PING", "FLUSHDB", "DBSIZE":
		if len(args) != 0 {
			return Command{}, fmt.Errorf("%s expects no arguments", name)
		}

	default:
		return Command{}, fmt.Errorf("unknown command: %s", name)
	}

	return Command{
		Name: name,
		Args: args,
	}, nil
}

func tokenize(input string) ([]string, error) {
	var args []string
	var current strings.Builder

	var quote rune
	escaped := false
	inToken := false

	for _, ch := range strings.TrimSpace(input) {
		if escaped {
			current.WriteRune(ch)
			escaped = false
			continue
		}

		if ch == '\\' && quote != 0 {
			escaped = true
			continue
		}

		if quote != 0 {
			if ch == quote {
				quote = 0
			} else {
				current.WriteRune(ch)
			}
			continue
		}

		if ch == '"' || ch == '\'' {
			quote = ch
			inToken = true
			continue
		}

		if ch == ' ' || ch == '\t' {
			if inToken {
				args = append(args, current.String())
				current.Reset()
				inToken = false
			}
			continue
		}

		current.WriteRune(ch)
		inToken = true
	}

	if escaped {
		return nil, fmt.Errorf("unfinished escape sequence")
	}

	if quote != 0 {
		return nil, fmt.Errorf("unterminated quote")
	}

	if inToken {
		args = append(args, current.String())
	}

	return args, nil
}

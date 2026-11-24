package resp

import (
	"fmt"
	"strconv"
)

func SimpleString(value string) []byte {
	return []byte("+" + value + "\r\n")
}

func Error(message string) []byte {
	return []byte("-" + message + "\r\n")
}

func Integer(value int64) []byte {
	return []byte(":" + strconv.FormatInt(value, 10) + "\r\n")
}

func BulkString(value string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len([]byte(value)), value))
}

func Null() []byte {
	return []byte("$-1\r\n")
}

func Array(values [][]byte) []byte {
	result := fmt.Sprintf("*%d\r\n", len(values))

	for _, value := range values {
		result += string(value)
	}

	return []byte(result)
}

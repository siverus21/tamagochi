package main

import (
	"errors"
	"fmt"
)

func main() {
	message := make([]string, 5)

	message[0] = "1"

	fmt.Println(message)
}

func printMessages(msg []string) error{
	if len(msg) == 0 {
		return errors.New("empty array")
	}
	msg[2] = "6"
	fmt.Println(msg)

	return nil
}
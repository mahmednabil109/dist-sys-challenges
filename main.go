package main

import (
	"bufio"
	"os"

	"github.com/mahmednabil109/dist-sys-challenges/node"
)

func main() {
	node := node.New(os.Stdout)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		node.HandleMessage(scanner.Bytes())
	}
}

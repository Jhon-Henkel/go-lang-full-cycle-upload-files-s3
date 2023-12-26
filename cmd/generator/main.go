package main

import (
	"fmt"
	"os"
)

func main() {
	index := 0
	for {
		file, err := os.Create(fmt.Sprintf("./tmp/file%d.txt", index))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		file.WriteString("Hello World!")
		index++
	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	f, err := os.Open("/var/log/system.log")
	if err != nil {
		panic(err.Error())
	}
	// f.Seek(0, 2)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			fmt.Println(line)
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(err.Error())
		}
	}
}

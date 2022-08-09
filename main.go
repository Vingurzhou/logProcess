package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

type Reader interface {
	Read(rc chan string)
}
type Writer interface {
	Write(wc chan string)
}

type LogProcess struct {
	rc    chan string
	wc    chan string
	read  Reader
	write Writer
}

type ReadFromFile struct {
	path string
}
type WriteToInfluxDB struct {
	influxDbDsn string
}

func (r *ReadFromFile) Read(rc chan string) {
	f, err := os.Open(r.path)
	if err != nil {
		panic(err.Error())
	}
	f.Seek(0, 2)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			continue
		} else if err != nil {
			panic(err.Error())
		}
		rc <- string(strings.TrimRight(string(line), "\n"))
	}
	// wg.Done()
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	for v := range wc {
		fmt.Println(v)
	}
	wg.Done()
}
func (l *LogProcess) Process() {
	for v := range l.rc {
		l.wc <- strings.ToUpper(v)
	}
	wg.Done()
}

func main() {
	r := &ReadFromFile{
		path: "./access.log",
	}
	w := &WriteToInfluxDB{
		influxDbDsn: "username&password..",
	}
	lp := &LogProcess{
		rc:    make(chan string),
		wc:    make(chan string),
		read:  r,
		write: w,
	}
	wg.Add(3)
	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)
	wg.Wait()
}

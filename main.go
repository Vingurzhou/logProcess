package main

import (
	"fmt"
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
	line := "message"
	rc <- line
	wg.Done()
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	fmt.Println(<-wc)
	wg.Done()
}
func (l *LogProcess) Process() {
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
	wg.Done()
}

func main() {
	r := &ReadFromFile{
		path: "",
	}
	w := &WriteToInfluxDB{
		influxDbDsn: "",
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

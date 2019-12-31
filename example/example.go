package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	//golog.InitLogger("", 0, false)
	start := time.Now()
	for i:= 0; i < 1000000; i++ {
		f, _ := os.OpenFile("aaa.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		s := strconv.Itoa(i)
		f.WriteString(s + "\n")
		f.Close()
	}
	end := time.Now()
	fmt.Println(time.Since(start).Seconds())
	f1, _ := os.OpenFile("aaa.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

	for i:= 0; i < 1000000; i++ {
		s := strconv.Itoa(i)
		f1.WriteString(s + "\n")
		f1.Sync()
	}

	f1.Close()
	fmt.Println(time.Since(end).Seconds())
}
package golog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var now = time.Now().Format("2006-01-02 15:04:05")

func control(name string, format string, args ...interface{}) {
	if stdOut {
		printLine(name, format, args...)
		return
	}
	if _, ok := logName[name]; !ok {
		path := filepath.Join(logPath, name+".log")

		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			// 如果失败，切换到控制台输出
			//panic(err)
			log.Println("Permission denied,  auto change to Stdout")
			stdOut = true
			return
		}

		logName[name] = &file{
			Mu:       &sync.Mutex{},
			Filebyte: f,
		}

	}
	logName[name].Mu.Lock()
	defer logName[name].Mu.Unlock()
	getLine(name, format, args...)
}

type file struct {
	Filebyte *os.File    // 文件
	Mu       *sync.Mutex // 文件锁
}

func Close(name string) {
	logName[name].Filebyte.Close()
}

func printLine(name string, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	msg := fmt.Sprintf("%s\t[%s]\t%s", now, name, line)
	log.Println(msg)

}

func getLine(name string, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)

	out := fmt.Sprintf("%s\t%s\n", now, line )
	write(name, out)

}

func write(name string, message string) {
	//fmt.Println("write")

	logName[name].Filebyte.WriteString(message)
	info, err := logName[name].Filebyte.Stat()
	if err != nil {
		fmt.Printf("not found %v \n", name)
	}
	localtime := time.Now()
	//每天一次来分割
	if everyDay && localtime.Day() != info.ModTime().Day() {
		prefix := fmt.Sprintf("%d-%d-%d", localtime.Year(), localtime.Month(), localtime.Day())

		os.Rename(info.Name(), prefix+info.Name())
		logName[name].Filebyte.Close()
		delete(logName, name)

	} else if fileSize > 0 && uint64(info.Size()) >= fileSize {
		// 根据文件大小来分割
		prefix := fmt.Sprintf("%d", time.Now().UnixNano())
		os.Rename(info.Name(), prefix+info.Name())
		logName[name].Filebyte.Close()
		delete(logName, name)
	}

}

func printfileline() string {

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	return fmt.Sprintln(file, line)
}

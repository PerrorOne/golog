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
	// 判断是输出控制台 还是写入文件
	if stdOut {
		printLine(name, format, args...)
		return
	} else {
		if everyDay {
			// 如果每天备份的话， 文件名需要更新
			localtime := time.Now()
			prefix := fmt.Sprintf("%d-%d-%d", localtime.Year(), localtime.Month(), localtime.Day())
			name = prefix + name
			// 删掉上一天的句柄
			localtime.AddDate(0,0,-1)
			lastlog := fmt.Sprintf("%d-%d-%d", localtime.Year(), localtime.Month(), localtime.Day())
			delname := filepath.Join(logPath, lastlog + name + ".log")
			if v, ok := logName[delname]; ok {
				v.Filebyte.Close()
				delete(logName, delname)
			}

			writeToFile(name , format , args...)
			return
		}

		// 如果按照文件大小判断的话，名字不变
		name = name + ".log"
		name := filepath.Join(logPath, name)
		writeToFile(name , format , args...)
	}

}

func writeToFile(name string, format string, args ...interface{}) {
	//
	if _, ok := logName[name]; !ok {
		//不存在就新建
		f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			// 如果失败，切换到控制台输出
			log.Println("Permission denied,  auto change to Stdout")
			stdOut = true
			printLine(name, format, args...)
			return
		}

		logName[name] = &file{
			Mu:       &sync.Mutex{},
			Filebyte: f,
		}

	}
	// 写入的时候加锁， 防止串行
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
	//line := fmt.Sprintf(format, args...)
	if len(args) > 0 {
		tmp := make([]interface{},0)
		tmp = append(tmp, now, name)
		tmp = append(tmp, args...)
		log.Printf("%s\t[%s]\t" + format + "\n", tmp...)
	} else {
		log.Printf("%s\t[%s]\t%s\n",now, name,format )
	}


}

// 写入文件
func getLine(name string, format string, args ...interface{}) {
	//合成字符串
	out := format
	if len(args) > 0 {
		tmp := make([]interface{},0)
		tmp = append(tmp, now)
		tmp = append(tmp, args...)
		out = fmt.Sprintf("%s\t" + format, tmp...)
	}
	write(name, out)

}

func write(name string, message string) {
	//fmt.Println("write")
	//
	if fileSize > 0 {
		//如果是按照大小切割
		info, err := logName[name].Filebyte.Stat()
		if err != nil {
			fmt.Printf("not found %v \n", name)
		}
		if info.Size() >= fileSize {

			prefix := fmt.Sprintf("%d", time.Now().UnixNano())
			// 删除句柄， 删除map
			logName[name].Filebyte.Close()
			delete(logName, name)
			// 移动文件，时间戳前缀
			os.Rename(info.Name(), prefix+info.Name())
			// 打开新文件
			f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				// 如果失败，切换到控制台输出
				log.Println("Permission denied,  auto change to Stdout")
				stdOut = true
				printLine(name, message)
				return
			}

			logName[name] = &file{
				Mu:       &sync.Mutex{},
				Filebyte: f,
			}


		}
	}
	logName[name].Filebyte.WriteString(message)


}

func printFileline() string {

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	return fmt.Sprintln(file, line)
}

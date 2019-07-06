package golog

import (
	"log"
	"os"
	"runtime"
)

var (
	Logpath  string // 文件路径
	FileSize int64  // 切割的文件大小
	EveryDay bool   // 每天一个来切割文件 （这个比上面个优先级高）
)

var LogName map[string]*file

func InitLogger(logpath string, filesize int64, everyday bool) {

	LogName = make(map[string]*file, 0)

	Logpath = addXieGang(logpath)
	//filepath.
	err := os.MkdirAll(Logpath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	FileSize = filesize * (1 << 20) // 默认单位M
	EveryDay = everyday

}

//  需要一个日志，多一条就好
// 最大缓存多少条日志
const MAXCACHELOG = 10000

// open file，  所有日志默认前面加了时间，
func Info(format string, args ...interface{}) {
	// info,
	name := "info" // 文件名 生成的文件为 info.log
	control(name, format, args...)
}

// open file，  所有日志默认前面加了时间，
func Access(format string, args ...interface{}) {
	// Access,
	name := "access" // 文件名 生成的文件为 info.log
	control(name, format, args...)
}

// open file，  所有日志默认前面加了时间，
func Debug(format string, args ...interface{}) {
	// debug,
	name := "debug" // 文件名 生成的文件为 debug.log
	control(name, format, args...)
}


func Email(format string, args ...interface{}) {
	// Email,
	name := "email" // 文件名 生成的文件为 email.log
	control(name, format, args...)
}

// 可以根据下面格式一样，在format 后加上更详细的输出值
func Error(format string, args ...interface{}) {
	// error日志，添加了错误函数，
	name := "error"
	format = printfileline() + format // printfileline()打印出错误的文件和行数
	control(name, format, args...)
}

func addXieGang(path string) string {
	l := len(path)
	// 如果是windows
	if runtime.GOOS == "windows" {
		if path[l-1:] != "\\" {
			return path + "\\"

		}

	} else {
		if path[l-1:] != "/" {
			return path + "/"

		}
	}
	return path
}

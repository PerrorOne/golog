# golog
### 安装
```
 go get github.com/hyahm/golog
```

### 使用
> 在main 函数开始直接调用InitLogger或 InitLoggerInConf 方法  
```
package main

import (
	"app"
	"github.com/hyahm/golog"
)

func main() {
	golog.InitLogger("/data/log", 0, false)
	//goconfig.InitConf("showgo.conf")
	app.ShowApp()
}
```
方法 1 , 直接指定配置(日志文件夹路径, 分割日志的大小, 单位Mb, 是否每天分割日志)   
```
golog.InitLogger("/data/log", 0, false)  
```
方法 2, 从配置文件读取  
```
golog.InitLoggerInConf()
```
已经带了, 后面任何地方直接调用方法即可写入, 每个方法都会生成一个文件   
```
golog.Info(format string, args ...interface{})
golog.Access(format string, args ...interface{})
golog.Email(format string, args ...interface{})
golog.Error(format string, args ...interface{})
```

除了Error方法里面的加了报错行外,其他的一样,  想要多余的文件记录, 可以复制格式改下名称就好  
```
func Error(format string, args ...interface{}) {
	// error日志，添加了错误函数，
	name := "error"
	format = printfileline() + format // printfileline()打印出错误的文件和行数
	control(name, format, args...)
}
func Email(format string, args ...interface{}) {
	// Email,
	name := "email" // 文件名 生成的文件为 info.log
	control(name, format, args...)  
}
```



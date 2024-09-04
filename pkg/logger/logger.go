package logger

import (
	"log"
	"os"
)

var file *os.File

func init() {
	var err error
	file, err = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(file)
}

func Log(logInfo string) {

	log.Output(2, logInfo) // 第二个参数是日志信息，第二个参数表示距离调用栈的深度
}

func Close() {
	file.Close()
}

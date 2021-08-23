package log

import (
	"io"
	"log"
	"os"
	"qiniu/setting"
)

var (
	AppLog *log.Logger // app运行日志
	HttpAccess *log.Logger // http运行日志
	HttpError *log.Logger // http错误日志
)

func InitLogger() {
	applogFile, err := os.OpenFile(setting.Setting().LogPath + "appLog.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln("Failed to open app log file:", err)
	}
	httpAccessFile, err := os.OpenFile(setting.Setting().LogPath + "httpAccess.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln("Failed to open http access log file:", err)
	}
	httpErrorFile, err := os.OpenFile(setting.Setting().LogPath + "httpError.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln("Failed to open http err log file:", err)
	}

	AppLog = log.New(io.MultiWriter(applogFile, os.Stderr), "", log.Ldate | log.Ltime | log.Llongfile)
	HttpAccess = log.New(io.MultiWriter(httpAccessFile, os.Stdout), "", log.Ldate | log.Ltime | log.Llongfile)
	HttpError = log.New(io.MultiWriter(httpErrorFile, os.Stderr), "", log.Ldate | log.Ltime | log.Llongfile)
}
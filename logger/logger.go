package logger

import (
	"log"
	"os"
)

var (
	Mode  = "dev"
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
)

func init() {
	var output *os.File
	if Mode == "prod" {
		file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		output = file
	} else {
		output = os.Stdout
	}

	Info = log.New(output, "INFO: ", log.Ldate|log.Ltime)
	Error = log.New(output, "ERROR: ", log.Ldate|log.Ltime)
	Warn = log.New(output, "WARN: ", log.Ldate|log.Ltime)

}

func Infof(format string, v ...interface{}) {
	Info.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	Error.Printf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	Warn.Printf(format, v...)
}

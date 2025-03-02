package logger

import (
	"log"
	"os"
)

var (
	Mode  = "dev"
	Info  *log.Logger
	Error *log.Logger
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

	Info = log.New(output, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(output, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Infof(format string, v ...interface{}) {
	Info.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	Error.Printf(format, v...)
}

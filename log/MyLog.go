package log

import (
	"fmt"
	"log"
	"os"
)

const pathLogFile = "log.txt"

type myLog struct {
	CommonLog *log.Logger
	ErrorLog  *log.Logger
}

func NewLog() *myLog {
	_, err := os.Stat(pathLogFile)

	if err != nil {
		fileCreated, err := os.Create(pathLogFile)

		if err != nil {
			fmt.Println("doesnt crate file...")
		}

		fileCreated.Close()

	}

	openLogfile, err := os.OpenFile(pathLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		if err != nil {
			fmt.Println("doesnt open file...")
		}
	}

	l := new(myLog)

	l.CommonLog = log.New(openLogfile, "Common Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)

	l.ErrorLog = log.New(openLogfile, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)

	return l

}

func (l *myLog) PrintCommon(text string) {
	l.CommonLog.Println(text)
}

func (l *myLog) PrintError(err error) {
	l.CommonLog.Println(err)
}

func (l *myLog) Fatal(err error) {
	log.Fatal(err)
}

package log

import (
	"FucknGO/config"
	"fmt"
	"log"
	"os"
)

type myLog struct {
	commonLog *log.Logger
	errorLog  *log.Logger
}

var debugResume = false

func NewLog() *myLog {
	pathLogFile := config.Config{}.JsonStr.Log.Path
	_, err := os.Stat(pathLogFile)

	if err != nil {
		fileCreated, err := os.Create(pathLogFile)

		if err != nil {
			fmt.Println("doesnt create file...")
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
	l.commonLog = log.New(openLogfile, "INFO:\t", log.Ldate|log.Ltime|log.Lshortfile)
	l.errorLog = log.New(openLogfile, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	config, err := config.GetConfig(config.Path)

	if err != nil {
		fmt.Println(err)
	}

	debugResume = config.JsonStr.Resume.IsDebug

	fmt.Println("port", config.JsonStr.ServerConfig.Port)

	return l
}

func (l *myLog) PrintCommon(text string) {
	if debugResume {
		l.commonLog.Println(text)
		fmt.Println(text)
	} else {
		l.commonLog.Println(text)
	}
}

func (l *myLog) PrintError(err error) {
	if debugResume {
		l.errorLog.Println(err)
		fmt.Println(err)
	} else {
		l.errorLog.Println(err)
	}
}

func (l *myLog) Fatal(err error) {
	log.Fatal(err)
}

package log

import (
	"FucknGO/config"
	"fmt"
	"log"
	"os"
	"strconv"
)

type myLog struct {
	commonLog *log.Logger
	errorLog  *log.Logger
}

var debugMode bool

func NewLog() *myLog {
	config, err := config.GetConfig(config.Path)

	pathLogFile := config.JsonStr.Log.Path

	_, err = os.Stat(pathLogFile)

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

	debug, errBool := os.LookupEnv("debug")

	if errBool {
		debug, err := strconv.ParseBool(debug)

		if err == nil {
			debugMode = debug
		}
	}

	return l
}

func (l *myLog) PrintCommon(text string) {
	if debugMode {
		l.commonLog.Println(text)
		fmt.Println(text)
	} else {
		l.commonLog.Println(text)
	}
}

func (l *myLog) PrintError(err error) {
	if debugMode {
		l.errorLog.Println(err)
		fmt.Println(err)
	} else {
		l.errorLog.Println(err)
	}
}

func (l *myLog) Fatal(err error) {
	log.Fatal(err)
}

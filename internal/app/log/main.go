package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

func SetOutputFile() {
	now := time.Now()
	name := fmt.Sprintf("log%s.log", now.Format("2006-01-02"))

	f, err := os.OpenFile(fmt.Sprintf("logs/%s", name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkError(err)
	/*defer func(){
		err := f.Close()
		checkError(err)
	}()*/

	logger.SetOutput(f)
}

func GetLogger() *log.Logger {
	return logger
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

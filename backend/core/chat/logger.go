package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

const LOG_FILE = "./log/logs"

type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	fileLog  *os.File
}

func NewLogger() *Logger {
	logFile, err := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}

	infoLog := log.New(logFile, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(logFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		infoLog:  infoLog,
		fileLog:  logFile,
		errorLog: errorLog,
	}
}

func (l *Logger) Fatal(err error) {
	fatalLog := log.New(l.fileLog, "[FATAL]\t", log.Ldate|log.Ltime)
	fatalLog.Println(err)
	defer l.fileLog.Close()

	fileLog, err := os.Open("log/logs")
	if err != nil {
		log.Println(err)
	}
	defer fileLog.Close()

	_, err = io.WriteString(l.fileLog, "")
	if err != nil {
		log.Println(err)
	}

	scan := bufio.NewScanner(fileLog)
	for scan.Scan() {
		line := scan.Text()
		log.Println(line)
	}
}

//infoLog := log.New(logFile, "INFO\t", log.Ldate|log.Ltime)
//infoLog.Println("Start servers")
//log.SetOutput(logFile)
//log.Println("Sending an entry to log!")

//logInfo := log.New(file, "INFO\t", log.Ldate|log.Ltime) //os.Stdout
//logInfo.Println("Start servers")

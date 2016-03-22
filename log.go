package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const p string = string(os.PathSeparator)

func getLogPath(level string) string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)

	paths := strings.Split(path, p)
	fileName := p + paths[len(paths)-1]
	logFile := p + "ddns-" + level + ".log"
	path = strings.Replace(path, fileName, logFile, 1)
	return path
}

func Exist(filename string) (os.FileInfo, bool) {
	fi, err := os.Stat(filename)
	return fi, err == nil || os.IsExist(err)
}

func log(message string, logfile string) {
	fi, exist := Exist(logfile)
	if exist {
		if fi.Size() > 1024*1024 {
			os.Rename(logfile, logfile+"."+time.Now().Format("20060102150405"))
		}
	}

	m := time.Now().Format("[01-02 15:04:05]") + " " + message + "\n"
	fil, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0420)
	if err != nil {
		return
	}
	defer fil.Close()
	fil.WriteString(m)
}

func Debug(msg string) {
	if config.LogConfig.Debug {
		log(msg, getLogPath("debug"))
	} else {
		fmt.Println(msg)
	}
}

func DebugR(msg string, record *Record) {
	msg = msg + " [" + record.RecordName + "." + record.DomainName + "]"
	log(msg, getLogPath("debug"))
}

func Info(msg string) {
	if config.LogConfig.Info {
		log(msg, getLogPath("info"))
	} else {
		fmt.Println(msg)
	}
}

func InfoR(msg string, record *Record) {
	msg = msg + " [" + record.RecordName + "." + record.DomainName + "]"
	log(msg, getLogPath("info"))
}

func Error(msg string) {
	if config.LogConfig.Error {
		log(msg, getLogPath("error"))
	} else {
		fmt.Println(msg)
	}
}

func ErrorR(msg string, record *Record) {
	msg = msg + " [" + record.RecordName + "." + record.DomainName + "]"
	log(msg, getLogPath("error"))
}

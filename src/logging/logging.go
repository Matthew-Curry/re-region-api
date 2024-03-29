/* Package exposing method to retrieving an aggregated logging to log on all required levels */

package logging

import (
	"log"
	"os"
	"errors"
	"io"
	"runtime"
	"fmt"
	"strings"
)

type Logger struct {
	infoLogger *log.Logger
	warnLogger *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

func (l *Logger) Info(s string, a ...interface{}) {
	s = l.addCaller(s)
	l.infoLogger.Printf(s, a...)
}

func (l *Logger) Warn(s string, a ...interface{}) {
	s = l.addCaller(s)
	l.warnLogger.Printf(s, a...)
}

func (l *Logger) Error(s string, a ...interface{}) {
	s = l.addCaller(s)
	l.errorLogger.Printf(s, a...)
}

func (l *Logger) Fatal(s string, a ...interface{}) {
	s = l.addCaller(s)
	l.errorLogger.Panicf(s, a...)
}

func (l *Logger) messageSetup(s string, a []interface{}) (string, []interface{}) {
	return l.addCaller(s), l.inputAsString(a)
}

func (l *Logger) inputAsString(a ...interface{}) []interface{} {
	n := make([]interface{}, len(a))
	for _, i := range a {
		n = append(n, i.(string))
	}

	return n
}

func (l *Logger) addCaller(s string) string{
	_, file, line, _ := runtime.Caller(2)
	fileParts := strings.Split(file, "/")
	source := fileParts[len(fileParts)-1]
	return source + " " + fmt.Sprint(line) + " " + s
}

// Public method to instantiate base loggers, return aggregated logger object
func GetLogger(logPath string) (Logger, *os.File) {
	var file *os.File
	if exists(logPath){
		file, _ = os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		file, _ = os.Create(logPath)
	}
	mw := io.MultiWriter(os.Stdout, file) 
	// common logging flags
	flags := log.LstdFlags
	infoLogger := log.New(mw, "INFO ", flags)
	warnLogger := log.New(mw, "WARN ", flags)
	errorLogger := log.New(mw, "ERROR ", flags)
	fatalLogger := log.New(mw, "FATAL ", flags)

	return Logger{infoLogger: infoLogger, warnLogger: warnLogger, errorLogger: errorLogger, fatalLogger: fatalLogger}, file
}

// helper method to check if given log file already exists
func exists(name string) bool {
    _, err := os.Stat(name)
    if err == nil {
        return true
    }
    if errors.Is(err, os.ErrNotExist) {
        return false
    }
    return false
}

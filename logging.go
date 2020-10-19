package gologger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	LOG_DEBUG = iota
	LOG_INFO
	LOG_WARN
	LOG_ERROR
)

type Logger struct {
	simple     *bool
	file       *os.File
	writeLevel *int
	tagAlign   *int
	name       string
	lock       sync.Mutex
}

func NewLogger() *Logger {
	simple := false
	writeLevel := LOG_DEBUG
	tagAlign := 0

	return &Logger{
		simple:     &simple,
		file:       nil,
		writeLevel: &writeLevel,
		tagAlign:   &tagAlign,
		name:       "main",
	}
}

func (l *Logger) Module(name string) *Logger {
	return &Logger{
		simple:     l.simple,
		file:       l.file,
		writeLevel: l.writeLevel,
		tagAlign:   l.tagAlign,
		name:       name,
	}
}

func (l *Logger) SetLevel(level int) *Logger {
	*l.writeLevel = level
	return l
}

func (l *Logger) SetSimple(simple bool) *Logger {
	*l.simple = simple
	return l
}

func (l *Logger) SetFile(path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	l.file = file
	return nil
}

func (l *Logger) SetTagAlign(n int) *Logger {
	*l.tagAlign = n
	return l
}

func formatLevel(level int) string {
	switch level {
	case LOG_DEBUG:
		return "\033[92m[Debug]\033[0m"
	case LOG_INFO:
		return "\033[96m[Info]\033[0m"
	case LOG_WARN:
		return "\033[93m[Warn]\033[0m"
	case LOG_ERROR:
		return "\033[91m[Error]\033[0m"
	default:
		return "UNKNOWN"
	}
}

func (l *Logger) formatMessage(level int, message []interface{}) string {
	var msgs []string
	for _, d := range message {
		msgs = append(msgs, fmt.Sprintf("%v", d))
	}
	now := time.Now()
	var tag string
	if *l.tagAlign != 0 && *l.tagAlign > len(l.name) {
		tag = l.name + strings.Repeat(" ", *l.tagAlign-len(l.name))
	} else {
		tag = l.name
	}
	if *l.simple {
		return fmt.Sprintf("[%s] [%s] %s %s", now.Format("15:04:05"), tag, formatLevel(level), strings.Join(msgs, " "))
	} else {
		pc, _, line, _ := runtime.Caller(3)
		return fmt.Sprintf("[%s] [%s] %s [%s:%d] %s", now.Format("2006/01/02 15:04:05.0000"), tag, formatLevel(level), runtime.FuncForPC(pc).Name(), line, strings.Join(msgs, " "))
	}
}

func (l *Logger) Log(level int, message ...interface{}) {
	if level < *l.writeLevel {
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	mess := l.formatMessage(level, message)
	fmt.Println(mess)
	if l.file != nil {
		l.file.WriteString(mess + "\n")
	}
}

func (l *Logger) Debug(message ...interface{}) {
	l.Log(LOG_DEBUG, message...)
}

func (l *Logger) Info(message ...interface{}) {
	l.Log(LOG_INFO, message...)
}

func (l *Logger) Warn(message ...interface{}) {
	l.Log(LOG_WARN, message...)
}

func (l *Logger) Error(message ...interface{}) {
	l.Log(LOG_ERROR, message...)
}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

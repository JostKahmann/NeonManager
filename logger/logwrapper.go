package logger

import l "log"

type LogLevel int

const (
	FATAL LogLevel = iota - 2
	PANIC
	ERROR
	WARN
	INFO
	DEBUG
)

func Log(level LogLevel, format string, args ...interface{}) {
	var levelString string
	switch level {
	case FATAL:
		l.Fatalf("[FATAL]: "+format, args)
	case PANIC:
		l.Panicf("[PANIC]: "+format, args)
	case ERROR:
		levelString = "[ERROR]"
	case WARN:
		levelString = "[WARN]"
	case INFO:
		levelString = "[INFO]"
	case DEBUG:
		levelString = "[DEBUG]"
	}
	if len(levelString) > 0 {
		l.Printf(levelString+": "+format, args...)
	}
}

func Debug(format string, args ...interface{}) {
	Log(DEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
	Log(INFO, format, args...)
}

func Warning(format string, args ...interface{}) {
	Log(WARN, format, args...)
}

func Warn(format string, args ...interface{}) {
	Log(WARN, format, args...)
}

func Error(format string, args ...interface{}) {
	Log(ERROR, format, args...)
}

func Panic(format string, args ...interface{}) {
	Log(PANIC, format, args...)
}

func Fatal(format string, args ...interface{}) {
	Log(FATAL, format, args...)
}

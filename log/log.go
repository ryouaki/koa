package log

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// LevelDebug
const (
	LevelDebug int = 1
	LevelTrace int = 2
	LevelInfo  int = 3
	LevelWarn  int = 4
	LevelError int = 5
)

// LogStd
const (
	LogStd  int = 0
	LogFile int = 1
)

// Config struct
type Config struct {
	Level    int
	Mode     int
	FileName string
	LogPath  string
	MaxDays  int
}

var logger *Logger = nil

// Logger struct
type Logger struct {
	level   int
	Writers []logWriter
}

// logWriter interface
type logWriter interface {
	Write(level int, msg string)
	Close()
}

// Data struct
type Data struct {
	level   int
	message string
}

type stdLogger struct {
	recv chan *Data
}

type fileLogger struct {
	recv     chan *Data
	file     *os.File
	FileName string
	LogPath  string
	currTime string
	level    int
	maxDays  int
}

// New func
func New(conf *Config) error {
	if logger != nil {
		return errors.New("Logger was created")
	}

	logger = &Logger{
		level: conf.Level,
	}

	switch conf.Mode {
	case LogStd:
		logger.Writers = make([]logWriter, 1)
		logger.Writers[0] = createStdLogger(conf)
	case LogFile:
		logger.Writers = make([]logWriter, 5)
		logger.Writers[0] = createFileLogger(LevelInfo, conf)
		logger.Writers[1] = createFileLogger(LevelWarn, conf)
		logger.Writers[2] = createFileLogger(LevelError, conf)
		logger.Writers[3] = createFileLogger(LevelTrace, conf)
		logger.Writers[4] = createFileLogger(LevelDebug, conf)
	}
	return nil
}

// Debug func
func Debug(msg string) {
	if logger.level <= LevelDebug {
		logger.write(LevelDebug, msg)
	}
}

// Info func
func Info(msg string) {
	if logger.level <= LevelInfo {
		logger.write(LevelInfo, msg)
	}
}

// Warn func
func Warn(msg string) {
	if logger.level <= LevelWarn {
		logger.write(LevelWarn, msg)
	}
}

// Error func
func Error(msg string) {
	if logger.level <= LevelError {
		logger.write(LevelError, msg)
	}
}

// Trace func
func Trace(msg string) {
	if logger.level <= LevelTrace {
		logger.write(LevelTrace, msg)
	}
}

func (log *Logger) write(level int, msg string) {
	for _, lw := range log.Writers {
		lw.Write(level, msg)
	}
}

func (std *stdLogger) Write(level int, msg string) {
	std.recv <- &Data{
		level:   level,
		message: msg,
	}
}

func (std *stdLogger) Close() {
	close(std.recv)
	return
}

func (std *fileLogger) Write(level int, msg string) {
	if level == std.level {
		std.recv <- &Data{
			level:   level,
			message: msg,
		}
	}
}

func (std *fileLogger) Close() {
	close(std.recv)
	std.file.Close()
	return
}

func createFileLogger(level int, conf *Config) logWriter {
	std := fileLogger{
		recv:    make(chan *Data),
		file:    nil,
		level:   level,
		maxDays: conf.MaxDays,
	}

	if conf.LogPath == "" {
		std.LogPath = "."
	} else {
		std.LogPath = conf.LogPath
	}

	if conf.FileName == "" {
		std.FileName = "log"
	} else {
		std.FileName = conf.FileName
	}

	go func() {
		defer func() {
			if std.file != nil {
				std.file.Close()
			}
		}()
		for {
			select {
			case res, _ := <-std.recv:
				now := time.Now()
				nowStr := fmt.Sprintf("%04d-%02d-%02d.%02d", now.Year(), now.Month(), now.Day(), now.Hour())
				if nowStr > std.currTime && std.file != nil {
					std.file.Close()
					std.file = nil
				}
				if std.file == nil {
					if std.currTime == "" || std.currTime != nowStr {
						std.currTime = nowStr
					}
					if std.maxDays > 0 {
						oldTime := now.AddDate(0, 0, 0-std.maxDays)
						oldTimeStr := fmt.Sprintf("%04d-%02d-%02d.%02d", oldTime.Year(), oldTime.Month(), oldTime.Day(), oldTime.Hour())
						oldFile := fmt.Sprintf("%s/%s.%s-%s.log", std.LogPath, getType(res.level), std.FileName, oldTimeStr)
						os.Remove(oldFile)
					}
					fd, err := os.OpenFile(fmt.Sprintf("%s/%s.%s-%s.log", std.LogPath, getType(res.level), std.FileName, std.currTime), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
					if err != nil {
						return
					}
					std.file = fd
				}
				_, err := fmt.Fprintln(std.file, fmt.Sprintf("[%s][%s] %s", getTime(), getType(res.level), res.message))
				if err != nil {
					std.file.Close()
					std.file = nil
				}
			}
		}
	}()

	return &std
}

func createStdLogger(conf *Config) logWriter {
	std := stdLogger{
		recv: make(chan *Data),
	}

	go func() {
		for {
			select {
			case res, _ := <-std.recv:
				fmt.Println(fmt.Sprintf("[%s][%s] %s", getTime(), getType(res.level), res.message))
			}
		}
	}()

	return &std
}

func getTime() string {
	now := time.Now()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
	)
}

func getType(mode int) string {
	switch mode {
	case LevelInfo:
		return "Info"
	case LevelWarn:
		return "Warn"
	case LevelError:
		return "Error"
	case LevelTrace:
		return "Trace"
	case LevelDebug:
		return "Debug"
	default:
		return "Logger"
	}
}

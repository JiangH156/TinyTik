package logger

import (
	"TinyTik/utils/files"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Settings struct {
	Path       string `yaml:"path"`
	Name       string `yaml:"name"`
	Ext        string `yaml:"ext"`
	TimeFormat string `yaml:"time-format"`
}

var (
	F                  *os.File
	DefaultPrefix      = ""
	DefaultCallerDepth = 2
	logger             *log.Logger
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup(settings *Settings) {
	var err error
	dir := settings.Path
	fileName := fmt.Sprintf("%s-%s.%s",
		settings.Name,
		time.Now().Format(settings.TimeFormat),
		settings.Ext)
	//fmt.Println(settings.TimeFormat, fileName)
	logFile, err := files.MustOpen(fileName, dir)
	if err != nil {
		log.Fatalf("logging.Setup err: %s", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(mw, DefaultPrefix, log.LstdFlags)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d] ", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s] ", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}

func Debug(v ...any) {
	setPrefix(DEBUG)
	logger.Println(v)
}
func Info(v ...any) {
	setPrefix(INFO)
	logger.Println(v)
}
func Warn(v ...any) {
	setPrefix(WARNING)
	logger.Println(v)
}
func Error(v ...any) {
	setPrefix(ERROR)
	logger.Println(v)
}
func Fatal(v ...any) {
	setPrefix(FATAL)
	logger.Println(v)
}

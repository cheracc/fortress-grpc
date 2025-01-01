package fortress

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	errorLogger   *log.Logger
	infoLogger    *log.Logger
	consoleLogger *log.Logger
	file          *os.File
}

const (
	errorPrefix      string = "[ERROR] "
	fatalPrefix      string = "[FATAL] "
	warnPrefix       string = "[WARN] "
	infoPrefix       string = "[INFO] "
	logInfoToConsole bool   = true
)

func NewLogger() *Logger {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	errorLogger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger := log.New(file, "", log.LstdFlags)
	consoleLogger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	logger := &Logger{errorLogger, infoLogger, consoleLogger, file}

	logger.Logf("New Logger started at %s", time.Now().String())
	return logger
}

// sends basic (debug) info to the logfile (and console if configured)
func (l *Logger) Logf(message string, args ...any) {
	logMessage := infoPrefix + message
	l.infoLogger.Printf(logMessage, args...)

	if logInfoToConsole {
		l.logToConsolef(logMessage, args...)
	}
}

func (l *Logger) Log(message string) {
	logMessage := infoPrefix + message
	l.infoLogger.Println(logMessage)

	if logInfoToConsole {
		l.logToConsole(logMessage)
	}
}

func (l *Logger) Warn(message string) {
	logMessage := warnPrefix + message
	l.errorLogger.Println(logMessage)
	l.logToConsole(logMessage)
}

func (l *Logger) Warnf(message string, args ...any) {
	logMessage := warnPrefix + message
	l.errorLogger.Printf(logMessage, args...)
	l.logToConsolef(logMessage, args...)
}

func (l *Logger) Error(message string) error {
	logMessage := errorPrefix + message
	l.errorLogger.Println(logMessage)
	l.logToConsole(logMessage)
	return fmt.Errorf("%s", logMessage)
}

func (l *Logger) Errorf(message string, args ...any) error {
	logMessage := errorPrefix + message
	l.errorLogger.Printf(logMessage, args...)
	l.logToConsolef(logMessage, args...)
	return fmt.Errorf(logMessage, args...)
}

func (l *Logger) Fatal(message string) {
	logMessage := fatalPrefix + message
	l.errorLogger.Println(logMessage)
	l.logToConsole(logMessage)
	defer os.Exit(1)
}

func (l *Logger) Fatalf(message string, args ...any) {
	logMessage := fatalPrefix + message
	l.errorLogger.Printf(logMessage, args...)
	l.logToConsolef(logMessage, args...)
	defer os.Exit(1)
}

func (l *Logger) logToConsole(message string) {
	l.consoleLogger.Output(3, message)
}

func (l *Logger) logToConsolef(message string, args ...any) {
	l.consoleLogger.Output(3, fmt.Sprintf(message, args...))
}

// sends a message to the console without the logger's prefix
func (l *Logger) ToConsole(message string) {
	fmt.Printf("\033[2K\r%s\n", message)
	l.infoLogger.Output(3, infoPrefix+message)
	fmt.Print("> ")
}

func (l *Logger) ToConsolef(message string, args ...any) {
	fmt.Printf("\033[2K\r%s\n", fmt.Sprintf(message, args...))
	l.infoLogger.Output(3, fmt.Sprintf(infoPrefix+message, args...))
	fmt.Print("> ")
}

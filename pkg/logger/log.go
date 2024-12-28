package logger

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Log struct
type Log struct {
	Method     string
	URL        string
	StatusCode string
	Time       string
}

// All variable for log
var (
	once        sync.Once
	infoLog     *os.File
	errorLog    *os.File
	infoLogger  zerolog.Logger
	errorLogger zerolog.Logger
)

// InitLogger initializes the loggers with given file paths for info and error logs.
func InitLogger(infoPath, errorPath string) error {
	var err error
	once.Do(func() {
		// Open the log files
		infoLog, err = os.OpenFile(infoPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}

		errorLog, err = os.OpenFile(errorPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}

		// Configure zerolog
		zerolog.TimeFieldFormat = time.RFC3339

		infoWriter := zerolog.MultiLevelWriter(os.Stdout, infoLog)
		errorWriter := zerolog.MultiLevelWriter(os.Stderr, errorLog)

		infoLogger = zerolog.New(infoWriter).With().Timestamp().Logger()
		errorLogger = zerolog.New(errorWriter).With().Timestamp().Logger()
	})

	return err
}

// Info logs an informational message.
// 5.160.41.78 - - [10/Oct/2024:12:18:29 +0330] "GET /DownloadCenter/index4/logo.svg HTTP/1.1" 200 2228 "https://hyperonline.shop/product/ep-7547/%DA%A9%D8%A7%D9%BE%D9%88%DA%86%DB%8C%D9%86%D9%88%20%D8%A7%DB%8C%D8%AA%D8%A7%D9%84%DB%8C%D8%A7%DB%8C%DB%8C%2025%20%DA%AF%D8%B1%D9%85%20%D9%81%D9%88%D9%84%20%DA%A9%D8%A7%D9%81%D9%87" "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Mobile Safari/537.36"
func Info(fields map[string]interface{}) {
	event := infoLogger.Info()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Send()
}

// Error logs an error message.
// // 5.160.41.78 - - [10/Oct/2024:12:18:29 +0330] "GET /DownloadCenter/index4/logo.svg HTTP/1.1" 200 2228 "https://hyperonline.shop/product/ep-7547/%DA%A9%D8%A7%D9%BE%D9%88%DA%86%DB%8C%D9%86%D9%88%20%D8%A7%DB%8C%D8%AA%D8%A7%D9%84%DB%8C%D8%A7%DB%8C%DB%8C%2025%20%DA%AF%D8%B1%D9%85%20%D9%81%D9%88%D9%84%20%DA%A9%D8%A7%D9%81%D9%87" "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Mobile Safari/537.36"
func Error(msg string, fields map[string]interface{}) {
	event := errorLogger.Error().Str("level", "error").Timestamp()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

// CloseLogger closes the log files.
func CloseLogger() {
	if infoLog != nil {
		infoLog.Close()
	}
	if errorLog != nil {
		errorLog.Close()
	}
}

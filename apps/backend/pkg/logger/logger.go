package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	level      LogLevel
	output     *os.File
	jsonFormat bool
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	File      string                 `json:"file,omitempty"`
	Line      int                    `json:"line,omitempty"`
}

func NewLogger() *Logger {
	level := INFO
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		switch envLevel {
		case "DEBUG":
			level = DEBUG
		case "INFO":
			level = INFO
		case "WARN":
			level = WARN
		case "ERROR":
			level = ERROR
		case "FATAL":
			level = FATAL
		}
	}

	jsonFormat := os.Getenv("LOG_FORMAT") == "json"

	return &Logger{
		level:      level,
		output:     os.Stdout,
		jsonFormat: jsonFormat,
	}
}

func (l *Logger) logWithLevel(level LogLevel, message string, data map[string]interface{}) {
	if level < l.level {
		return
	}

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level.String(),
		Message:   message,
		Data:      data,
		File:      file,
		Line:      line,
	}

	if l.jsonFormat {
		jsonData, _ := json.Marshal(entry)
		fmt.Fprintln(l.output, string(jsonData))
	} else {
		var dataStr string
		if len(data) > 0 {
			jsonData, _ := json.Marshal(data)
			dataStr = fmt.Sprintf(" | %s", string(jsonData))
		}

		fmt.Fprintf(l.output, "[%s] %s: %s%s\n",
			entry.Timestamp,
			entry.Level,
			entry.Message,
			dataStr)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) Debug(message string, data ...map[string]interface{}) {
	var logData map[string]interface{}
	if len(data) > 0 {
		logData = data[0]
	}
	l.logWithLevel(DEBUG, message, logData)
}

func (l *Logger) Info(message string, data ...map[string]interface{}) {
	var logData map[string]interface{}
	if len(data) > 0 {
		logData = data[0]
	}
	l.logWithLevel(INFO, message, logData)
}

func (l *Logger) Warn(message string, data ...map[string]interface{}) {
	var logData map[string]interface{}
	if len(data) > 0 {
		logData = data[0]
	}
	l.logWithLevel(WARN, message, logData)
}

func (l *Logger) Error(message string, data ...map[string]interface{}) {
	var logData map[string]interface{}
	if len(data) > 0 {
		logData = data[0]
	}
	l.logWithLevel(ERROR, message, logData)
}

func (l *Logger) Fatal(message string, data ...map[string]interface{}) {
	var logData map[string]interface{}
	if len(data) > 0 {
		logData = data[0]
	}
	l.logWithLevel(FATAL, message, logData)
}

func (l *Logger) LogError(err error, context map[string]interface{}) {
	if context == nil {
		context = make(map[string]interface{})
	}
	context["error"] = err.Error()
	l.logWithLevel(ERROR, "Error occurred", context)
}

func (l *Logger) LogRequest(method, path, ip, userAgent string, status int, duration time.Duration) {
	l.logWithLevel(INFO, "HTTP Request", map[string]interface{}{
		"method":     method,
		"path":       path,
		"ip":         ip,
		"user_agent": userAgent,
		"status":     status,
		"duration":   duration.String(),
		"duration_ms": duration.Milliseconds(),
	})
}

func (l *Logger) LogDatabase(operation, table string, duration time.Duration, rowsAffected int64) {
	l.logWithLevel(DEBUG, "Database Operation", map[string]interface{}{
		"operation":     operation,
		"table":         table,
		"duration":      duration.String(),
		"duration_ms":   duration.Milliseconds(),
		"rows_affected": rowsAffected,
	})
}

func (l *Logger) LogCache(operation, key string, hit bool, duration time.Duration) {
	status := "miss"
	if hit {
		status = "hit"
	}

	l.logWithLevel(DEBUG, "Cache Operation", map[string]interface{}{
		"operation":   operation,
		"key":         key,
		"status":      status,
		"duration":    duration.String(),
		"duration_ms": duration.Milliseconds(),
	})
}

func (l *Logger) LogAuth(action, username, ip string, success bool) {
	status := "failed"
	if success {
		status = "success"
	}

	l.logWithLevel(INFO, "Authentication Event", map[string]interface{}{
		"action":   action,
		"username": username,
		"ip":       ip,
		"status":   status,
	})
}

func (l *Logger) LogBookOperation(operation, username string, bookID uint, title string) {
	l.logWithLevel(INFO, "Book Operation", map[string]interface{}{
		"operation": operation,
		"username":  username,
		"book_id":   bookID,
		"title":     title,
	})
}

func (l *Logger) LogStartup(version, env string, config map[string]interface{}) {
	l.logWithLevel(INFO, "Application Starting", map[string]interface{}{
		"version": version,
		"env":     env,
		"config":  config,
	})
}

func (l *Logger) LogShutdown(reason string) {
	l.logWithLevel(INFO, "Application Shutting Down", map[string]interface{}{
		"reason": reason,
	})
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) SetOutput(output *os.File) {
	l.output = output
}

func (l *Logger) SetJSONFormat(enabled bool) {
	l.jsonFormat = enabled
}

// WithFields returns a logger with preset fields
func (l *Logger) WithFields(fields map[string]interface{}) *FieldLogger {
	return &FieldLogger{
		logger: l,
		fields: fields,
	}
}

// FieldLogger is a logger with preset fields
type FieldLogger struct {
	logger *Logger
	fields map[string]interface{}
}

// mergeFields merges preset fields with additional fields
func (fl *FieldLogger) mergeFields(additional map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	// Copy preset fields
	for k, v := range fl.fields {
		merged[k] = v
	}

	// Copy additional fields (they can override preset fields)
	for k, v := range additional {
		merged[k] = v
	}

	return merged
}

// Debug logs a debug message with preset fields
func (fl *FieldLogger) Debug(message string, data ...map[string]interface{}) {
	var logData map[string]interface{}
	if len(data) > 0 {
		logData = fl.mergeFields(data[0])
	} else {
		logData = fl.fields
	}
	fl.logger.logWithLevel(DEBUG, message, logData)
}

// Info logs an info message with preset fields
func (fl *FieldLogger) Info(message string, data ...map[string]interface{}) {
	var logData map[string]interface{}
	if len(data) > 0 {
		logData = fl.mergeFields(data[0])
	} else {
		logData = fl.fields
	}
	fl.logger.logWithLevel(INFO, message, logData)
}

// Warn logs a warning message with preset fields
func (fl *FieldLogger) Warn(message string, data ...map[string]interface{}) {
	var logData map[string]interface{}
	if len(data) > 0 {
		logData = fl.mergeFields(data[0])
	} else {
		logData = fl.fields
	}
	fl.logger.logWithLevel(WARN, message, logData)
}

// Error logs an error message with preset fields
func (fl *FieldLogger) Error(message string, data ...map[string]interface{}) {
	var logData map[string]interface{}
	if len(data) > 0 {
		logData = fl.mergeFields(data[0])
	} else {
		logData = fl.fields
	}
	fl.logger.logWithLevel(ERROR, message, logData)
}

// GetStandardLogger returns a standard library logger for compatibility
func (l *Logger) GetStandardLogger() *log.Logger {
	return log.New(l.output, "", 0)
}

// Global logger instance
var globalLogger *Logger

// Init initializes the global logger
func Init(level, format string) {
	os.Setenv("LOG_LEVEL", level)
	os.Setenv("LOG_FORMAT", format)
	globalLogger = NewLogger()
}

// Global logging functions
func Debug(message string, data ...map[string]interface{}) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.Debug(message, data...)
}

func Info(message string, data ...map[string]interface{}) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.Info(message, data...)
}

func Warn(message string, data ...map[string]interface{}) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.Warn(message, data...)
}

func Error(message string, data ...map[string]interface{}) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.Error(message, data...)
}

func Fatal(message string, data ...map[string]interface{}) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.Fatal(message, data...)
}

// InfoWithData logs an info message with structured data
func InfoWithData(message string, data map[string]interface{}) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.Info(message, data)
}

// ErrorWithData logs an error message with structured data
func ErrorWithData(message string, data map[string]interface{}) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.Error(message, data)
}

// Global specialized logging functions
func LogRequest(method, path, ip, userAgent string, status int, duration time.Duration) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.LogRequest(method, path, ip, userAgent, status, duration)
}

func LogDatabase(operation, table string, duration time.Duration, rowsAffected int64) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.LogDatabase(operation, table, duration, rowsAffected)
}

func LogCache(operation, key string, hit bool, duration time.Duration) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.LogCache(operation, key, hit, duration)
}

func LogAuth(action, username, ip string, success bool) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.LogAuth(action, username, ip, success)
}

func LogBookOperation(operation, username string, bookID uint, bookTitle string) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.LogBookOperation(operation, username, bookID, bookTitle)
}

func LogError(err error, context map[string]interface{}) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.LogError(err, context)
}

func LogStartup(version, environment string, config map[string]interface{}) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.LogStartup(version, environment, config)
}

func LogShutdown(reason string) {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	globalLogger.LogShutdown(reason)
}

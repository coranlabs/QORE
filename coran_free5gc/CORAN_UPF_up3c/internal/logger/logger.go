package logger

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ColorFormatter struct{}

var (
	Log      *log.Logger
	NfLog    *log.Entry
	MainLog  *log.Entry
	CfgLog   *log.Entry
	Pfcplog  *log.Entry
	BuffLog  *log.Entry
	PerioLog *log.Entry
	FwderLog *log.Entry
	AppLog   *log.Entry
	InitLog  *log.Entry
	ServerLog *log.Entry
	PDRLog    *log.Entry
	QERLog    *log.Entry
)

// ANSI color codes for different log levels
var (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorBlue    = "\033[38;2;94;205;193m"
	colorYellow  = "\033[33m"
	colorMagenta = "\033[35m"
)

// Format implements the Formatter interface
func (f *ColorFormatter) Format(entry *log.Entry) ([]byte, error) {
	// Define color based on log level
	var color string
	switch entry.Level {
	case log.DebugLevel:
		color = colorMagenta
	case log.InfoLevel:
		color = colorBlue
	case log.WarnLevel:
		color = colorYellow
	case log.ErrorLevel:
		color = colorRed
	case log.FatalLevel:
		color = colorRed
	default:
		color = colorReset
	}

	// Format the timestamp
	timestamp := entry.Time.Format("2006-01-02 | 15:04:05")

	// Format log level with 1-character padding between the log level and brackets
	coloredLevel := fmt.Sprintf("%s[ %-5s ]%s", color, strings.ToUpper(entry.Level.String()), colorReset)

	// Get the component field from the log entry, default to "general" if not provided
	component, exists := entry.Data["component"]
	if !exists {
		component = "general"
	}

	// Format component with a fixed width of 4 characters, with 1-character padding
	coloredComponent := fmt.Sprintf("%s[ %-4s ]%s", color, strings.ToUpper(component.(string)), colorReset)

	// Format the log message
	logLine := fmt.Sprintf(
		"%s %s %s | %s\n",
		timestamp, coloredLevel, coloredComponent, entry.Message,
	)

	return []byte(logLine), nil
}

// InitializeLogger sets the custom logger with ColorFormatter
func InitializeLogger(level log.Level) {
	log.SetFormatter(&ColorFormatter{})
	log.SetLevel(level) // Set default log level
	NfLog = CreateLoggerWithComponent("NF")
	MainLog = CreateLoggerWithComponent("MAIN")
	CfgLog = CreateLoggerWithComponent("CFG")
	Pfcplog = CreateLoggerWithComponent("PFCP")
	BuffLog = CreateLoggerWithComponent("BUFF")
	PerioLog = CreateLoggerWithComponent("PERIO")
	FwderLog = CreateLoggerWithComponent("FWD")
	AppLog = CreateLoggerWithComponent("APP")
	InitLog = CreateLoggerWithComponent("INIT")
	ServerLog = CreateLoggerWithComponent("SERVER")
	PDRLog = CreateLoggerWithComponent("PDR")
	QERLog = CreateLoggerWithComponent("QER")
	Log = log.StandardLogger()
}

// GetLogger allows custom log level configuration or any other setup before using the logger
func GetLogger() *log.Logger {
	return log.StandardLogger()
}

// CreateLoggerWithComponent returns a logger instance with predefined fields like component
func CreateLoggerWithComponent(component string) *log.Entry {
	return log.WithFields(log.Fields{
		"component": strings.ToUpper(component),
	})
}

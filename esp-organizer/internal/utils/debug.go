package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// DebugConfig holds the configuration for debug operations
type DebugConfig struct {
	Enabled      bool
	LogToFile    bool
	LogFilePath  string
	LogRequests  bool
	LogResponses bool
	Verbose      bool
}

var (
	// DefaultDebugConfig is the default debug configuration
	DefaultDebugConfig = DebugConfig{
		Enabled:      getEnvBool("DEBUG_MODE", "DEBUG", "true"),
		LogToFile:    getEnvBool("DEBUG_LOG_TO_FILE", "", "true"),
		LogFilePath:  getEnvWithDefault("DEBUG_LOG_PATH", "./debug.log"),
		LogRequests:  true,
		LogResponses: true,
		Verbose:      getEnvBool("DEBUG_VERBOSE", "", "true"),
	}

	// Track the current log file path that's being used
	currentLogPath string
)

// Debugger provides debugging utilities
type Debugger struct {
	config DebugConfig
	logger *log.Logger
}

// NewDebugger creates a new debugger with the given configuration
func NewDebugger(config DebugConfig) *Debugger {
	// Start with a simple stdout logger
	logger := log.New(os.Stdout, "[DEBUG] ", log.LstdFlags)

	// Debug the actual environment variable values to help troubleshoot
	debugEnvValue := os.Getenv("DEBUG_LOG_TO_FILE")
	log.Printf("📊 DEBUG_LOG_TO_FILE environment value: '%s'", debugEnvValue)

	// Force true for LogToFile if environment variable is "true" (ignoring case)
	if strings.EqualFold(debugEnvValue, "true") {
		config.LogToFile = true
		log.Printf("✅ Forcing LogToFile=true based on environment variable")
	}

	var logFile *os.File

	if config.LogToFile {
		// Get the main application log file first (for combined logging)
		mainLogPath := "api-server.log" // This is your main app log

		// Try to open the main log file first
		var err error
		logFile, err = os.OpenFile(mainLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			absPath, _ := filepath.Abs(mainLogPath)
			log.Printf("✅ Using main application log file: %s", absPath)
			currentLogPath = absPath

			// Create a multi-writer to write to both stdout and the log file
			multiWriter := io.MultiWriter(os.Stdout, logFile)
			logger = log.New(multiWriter, "[DEBUG] ", log.LstdFlags)

			// Redirect standard log package to also write to our file
			log.SetOutput(multiWriter)
		} else {
			log.Printf("⚠️ Could not open main log file, falling back to separate debug log")

			// Try the custom debug log path
			if config.LogFilePath != "" {
				// Create logs directory if needed
				logDir := filepath.Dir(config.LogFilePath)
				if logDir != "" && logDir != "." {
					os.MkdirAll(logDir, 0755)
				}

				// Try to open the custom debug log
				logFile, err = os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err == nil {
					absPath, _ := filepath.Abs(config.LogFilePath)
					log.Printf("✅ Using separate debug log file: %s", absPath)
					currentLogPath = absPath

					// Create a multi-writer for stdout and log file
					multiWriter := io.MultiWriter(os.Stdout, logFile)
					logger = log.New(multiWriter, "[DEBUG] ", log.LstdFlags)
				} else {
					log.Printf("❌ Failed to open debug log file: %v", err)
					currentLogPath = "CONSOLE_ONLY"
				}
			}
		}
	} else {
		log.Printf("ℹ️ File logging is disabled. Set DEBUG_LOG_TO_FILE=true to enable.")
		currentLogPath = "CONSOLE_ONLY"
	}

	// Log the full debug configuration
	if config.Enabled {
		log.Printf("🔍 Debug mode enabled with config: LogToFile=%v, Path=%s, Verbose=%v",
			config.LogToFile, config.LogFilePath, config.Verbose)
	}

	return &Debugger{
		config: config,
		logger: logger,
	}
}

// GetCurrentLogPath returns the current log file path being used
func GetCurrentLogPath() string {
	return currentLogPath
}

// InitializeDebugger ensures the debug system is properly set up
func InitializeDebugger() *Debugger {
	// Check environment variables again in case they were set after program start
	debugEnvValue := os.Getenv("DEBUG_LOG_TO_FILE")
	logToFile := strings.EqualFold(debugEnvValue, "true")

	config := DebugConfig{
		Enabled:      getEnvBool("DEBUG_MODE", "DEBUG", "true"),
		LogToFile:    logToFile, // Use direct check instead of getEnvBool
		LogFilePath:  getEnvWithDefault("DEBUG_LOG_PATH", "./debug.log"),
		LogRequests:  true,
		LogResponses: true,
		Verbose:      getEnvBool("DEBUG_VERBOSE", "", "true"),
	}

	log.Printf("🔍 Re-initializing debug configuration from environment. DEBUG_LOG_TO_FILE='%s', LogToFile=%v",
		debugEnvValue, config.LogToFile)

	return NewDebugger(config)
}

// Helper function to get environment variable with fallbacks
func getEnvWithDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// Helper function to get boolean environment variable with fallbacks
func getEnvBool(key, altKey, trueVal string) bool {
	val := os.Getenv(key)
	if val == "" && altKey != "" {
		val = os.Getenv(altKey)
	}

	// Add extra debugging to troubleshoot environment variable issues
	if key == "DEBUG_LOG_TO_FILE" {
		log.Printf("📊 getEnvBool for %s: value='%s', comparing with '%s', result=%v",
			key, val, trueVal, strings.EqualFold(val, trueVal))
	}

	return strings.EqualFold(val, trueVal)
}

// Log logs a debug message if debugging is enabled
func (d *Debugger) Log(requestID, format string, args ...interface{}) {
	if !d.config.Enabled {
		return
	}

	message := fmt.Sprintf(format, args...)
	d.logger.Printf("[%s] %s", requestID, message)
}

// LogObject logs an object as JSON if debugging is enabled
func (d *Debugger) LogObject(requestID, label string, obj interface{}) {
	if !d.config.Enabled {
		return
	}

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		d.logger.Printf("[%s] Error marshaling %s: %v", requestID, label, err)
		return
	}

	d.logger.Printf("[%s] %s: %s", requestID, label, string(data))
}

// TraceFunction logs entry and exit of a function with timing information
func (d *Debugger) TraceFunction(requestID string) func() {
	if !d.config.Enabled || !d.config.Verbose {
		return func() {}
	}

	pc, _, _, ok := runtime.Caller(1)
	funcName := "unknown"
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}

	d.logger.Printf("[%s] ENTER: %s", requestID, funcName)
	startTime := time.Now()

	return func() {
		duration := time.Since(startTime)
		d.logger.Printf("[%s] EXIT: %s (took %v)", requestID, funcName, duration)
	}
}

// Default global debugger
var DefaultDebugger = InitializeDebugger()

// Log logs a message to the default debugger
func Log(requestID, format string, args ...interface{}) {
	DefaultDebugger.Log(requestID, format, args...)
}

// LogObject logs an object to the default debugger
func LogObject(requestID, label string, obj interface{}) {
	DefaultDebugger.LogObject(requestID, label, obj)
}

// TraceFunction traces a function with the default debugger
func TraceFunction(requestID string) func() {
	return DefaultDebugger.TraceFunction(requestID)
}

// ReloadDebugger reinitializes the debug configuration from environment variables
func ReloadDebugger() {
	DefaultDebugger = InitializeDebugger()
	log.Printf("✅ Debug configuration reloaded from environment variables")
}

// Add this function to force debug file logging on regardless of environment variables
func ForceFileLogging(logPath string) {
	if logPath == "" {
		logPath = "./debug.log"
	}

	config := DebugConfig{
		Enabled:      true,
		LogToFile:    true,
		LogFilePath:  logPath,
		LogRequests:  true,
		LogResponses: true,
		Verbose:      true,
	}

	DefaultDebugger = NewDebugger(config)
	log.Printf("✅ Forced file logging enabled to: %s", logPath)
}

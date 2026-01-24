package utils

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// ScanLogger collects warnings and errors during scanning and writes them to a log file
type ScanLogger struct {
	mu       sync.Mutex
	entries  []string
	logFile  string
	hasError bool
}

var (
	globalLogger *ScanLogger
	loggerOnce   sync.Once
)

// GetLogger returns the global scan logger instance
func GetLogger() *ScanLogger {
	loggerOnce.Do(func() {
		timestamp := time.Now().Format("20060102-150405")
		globalLogger = &ScanLogger{
			logFile: fmt.Sprintf("cf-nuke-scan-%s.log", timestamp),
			entries: make([]string, 0),
		}
	})
	return globalLogger
}

// LogWarning logs a warning message
func (l *ScanLogger) LogWarning(format string, args ...interface{}) {
	l.log("WARNING", format, args...)
}

// LogError logs an error message
func (l *ScanLogger) LogError(format string, args ...interface{}) {
	l.mu.Lock()
	l.hasError = true
	l.mu.Unlock()
	l.log("ERROR", format, args...)
}

func (l *ScanLogger) log(level string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	entry := fmt.Sprintf("[%s] %s: %s", timestamp, level, message)

	l.mu.Lock()
	l.entries = append(l.entries, entry)
	l.mu.Unlock()
}

// HasEntries returns true if there are any logged entries
func (l *ScanLogger) HasEntries() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries) > 0
}

// EntryCount returns the number of logged entries
func (l *ScanLogger) EntryCount() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries)
}

// LogFilePath returns the path to the log file
func (l *ScanLogger) LogFilePath() string {
	return l.logFile
}

// Flush writes all entries to the log file
func (l *ScanLogger) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.entries) == 0 {
		return nil
	}

	file, err := os.Create(l.logFile)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer file.Close()

	for _, entry := range l.entries {
		if _, err := fmt.Fprintln(file, entry); err != nil {
			return fmt.Errorf("failed to write to log file: %w", err)
		}
	}

	return nil
}

// PrintSummary prints a red warning message if there were errors/warnings
func (l *ScanLogger) PrintSummary() {
	if !l.HasEntries() {
		return
	}

	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("\n%s %d warnings/errors occurred during scan. See %s for details.\n",
		red("[SCAN ISSUES]"), l.EntryCount(), l.logFile)
}

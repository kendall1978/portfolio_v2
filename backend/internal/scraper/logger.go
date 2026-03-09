package scraper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// ScrapeLogger writes to both a log file and stdout.
type ScrapeLogger struct {
	file   *os.File
	logger *log.Logger
	Path   string
}

// NewScrapeLogger creates a timestamped log file in the given directory.
func NewScrapeLogger(logDir string) (*ScrapeLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	filename := fmt.Sprintf("scrape_%s.log", time.Now().Format("2006-01-02_150405"))
	path := filepath.Join(logDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("create log file: %w", err)
	}

	logger := log.New(f, "", log.LstdFlags)

	return &ScrapeLogger{file: f, logger: logger, Path: path}, nil
}

// Log writes a formatted message to both the log file and stdout.
func (sl *ScrapeLogger) Log(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	sl.logger.Println(msg)
	log.Println(msg)
}

// Close closes the log file.
func (sl *ScrapeLogger) Close() {
	sl.file.Close()
}
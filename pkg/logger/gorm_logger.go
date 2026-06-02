package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/gorm/logger"
)

// FileLogger handles daily log rotation and cleanup
type FileLogger struct {
	mu          sync.Mutex
	baseDir     string
	filePrefix  string
	currentDay  string
	file        *os.File
	maxFileDays int

	fallbackLogFile *os.File
}

func NewFileLogger(baseDir, filePrefix string, maxFileDays int, fallbackLogFile *os.File) *FileLogger {
	fl := &FileLogger{
		baseDir:         baseDir,
		filePrefix:      filePrefix,
		maxFileDays:     maxFileDays,
		fallbackLogFile: fallbackLogFile,
	}
	fl.rotateFile()
	return fl
}

func (fl *FileLogger) Write(p []byte) (n int, err error) {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	currentDay := time.Now().Format("2006-01-02")
	if currentDay != fl.currentDay {
		fl.rotateFile()
	}

	return fl.file.Write(p)
}

func (fl *FileLogger) rotateFile() {
	if fl.file != nil {
		fl.file.Close()
	}

	_ = os.MkdirAll(fl.baseDir, 0755)

	fl.currentDay = time.Now().Format("2006-01-02")
	filename := filepath.Join(
		fl.baseDir,
		fmt.Sprintf("%s_%s.log", fl.filePrefix, fl.currentDay),
	)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("cannot open log file: %v", err)
		fl.file = fl.fallbackLogFile
	} else {
		fl.file = f
	}

	// Bersihkan file log lama
	if fl.maxFileDays > 0 {
		fl.cleanupOldFiles()
	}
}

func (fl *FileLogger) cleanupOldFiles() {
	files, err := os.ReadDir(fl.baseDir)
	if err != nil {
		log.Printf("failed to read log dir: %v", err)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -fl.maxFileDays)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if !filepath.HasPrefix(name, fl.filePrefix+"_") || filepath.Ext(name) != ".log" {
			continue
		}

		path := filepath.Join(fl.baseDir, name)
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			_ = os.Remove(path)
			log.Printf("🧹 removed old log: %s", path)
		}
	}
}

// StartDailyRotation rotates the log file automatically at midnight
func (fl *FileLogger) StartDailyRotation() {
	go func() {
		for {
			now := time.Now()
			next := now.Add(24 * time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			time.Sleep(time.Until(next))

			fl.mu.Lock()
			fl.rotateFile()
			fl.mu.Unlock()
		}
	}()
}

// GetGormLogger creates a GORM logger that writes to rotating files
func GetGormLogger(baseDir, filePrefix string, maxFileDays int, level logger.LogLevel, fallbackLogFile *os.File) logger.Interface {
	fileLogger := NewFileLogger(baseDir, filePrefix, maxFileDays, fallbackLogFile)
	fileLogger.StartDailyRotation()

	newLogger := logger.New(
		log.New(fileLogger, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      level, // bisa diatur: logger.Info, logger.Warn, logger.Error, dll
			Colorful:      false,
		},
	)
	return newLogger
}

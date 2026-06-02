package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"github.com/Ajulll22/payment-ai-assistant/pkg/generator"
)

func NewLogger(cfgLog constant.LogConfig) *Logger {
	randNum := generator.RandomNumber(1, 99999)

	layoutFormat := "2006-01-02"
	t := time.Now()

	if err := os.MkdirAll(cfgLog.Path, os.ModePerm); err != nil {
		log.Fatalf("create folder log error: %v", err)
	}

	appDebugFile, err := os.OpenFile(filepath.Join(cfgLog.Path, cfgLog.DebugFilename+"_"+fmt.Sprint(t.Format(layoutFormat))+".log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	appErrorFile, err := os.OpenFile(filepath.Join(cfgLog.Path, cfgLog.ErrorFilename+"_"+fmt.Sprint(t.Format(layoutFormat))+".log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	logger := &Logger{
		traceID: strconv.Itoa(randNum),
	}
	logger.App = &loggerSet{
		parent:   logger,
		debugLog: log.New(appDebugFile, "", log.LstdFlags|log.Lmicroseconds),
		errorLog: log.New(appErrorFile, "", log.LstdFlags|log.Lmicroseconds),
	}

	go cleanupOldLogs(logger, cfgLog)

	return logger
}

type loggerSet struct {
	debugLog *log.Logger
	errorLog *log.Logger
	parent   *Logger
}

type Logger struct {
	traceID    string
	identities []string
	mu         sync.Mutex
	App        *loggerSet
}

func (l *loggerSet) Process(message ...any) {
	l.parent.mu.Lock()
	defer l.parent.mu.Unlock()

	b := fmt.Sprintln(message...)

	identity := "| [" + l.parent.traceID + "]"
	for _, val := range l.parent.identities {
		identity += " | [" + val + "]"
	}
	l.debugLog.Print(identity, " ", b)
}

func (l *loggerSet) Map(message ...any) {
	l.parent.mu.Lock()
	defer l.parent.mu.Unlock()

	b, _ := json.MarshalIndent(message, "", "  ")

	identity := "| [" + l.parent.traceID + "]"
	for _, val := range l.parent.identities {
		identity += " | [" + val + "]"
	}
	l.debugLog.Print(identity, " ", string(b))
}

func (l *loggerSet) Error(message ...any) {
	l.parent.mu.Lock()
	defer l.parent.mu.Unlock()

	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	caller := runtime.FuncForPC(pc[0])
	file, line := caller.FileLine(pc[0])

	fileArr := strings.Split(file, "/")
	lenFile := len(fileArr)

	fileName := "/" + fileArr[lenFile-2] + "/" + fileArr[lenFile-1]

	trace := fmt.Sprintf("%s:%d", fileName, line)

	b := fmt.Sprintln(message...)

	identity := "| [" + l.parent.traceID + "]"
	for _, val := range l.parent.identities {
		identity += " | [" + val + "]"
	}
	l.errorLog.Print(trace, " ", identity, " ", b)
}

func (l *loggerSet) ErrorWithoutTrace(message ...any) {
	l.parent.mu.Lock()
	defer l.parent.mu.Unlock()

	b := fmt.Sprintln(message...)

	identity := "| [" + l.parent.traceID + "]"
	for _, val := range l.parent.identities {
		identity += " | [" + val + "]"
	}
	l.errorLog.Print(identity, " ", b)
}

func (l *loggerSet) close() {
	if l == nil {
		return
	}

	if f, ok := l.debugLog.Writer().(*os.File); ok {
		_ = f.Close()
	}
	if f, ok := l.errorLog.Writer().(*os.File); ok {
		_ = f.Close()
	}
}

func (l *Logger) WithIdentity(identity string) *Logger {
	identitiesCopy := make([]string, len(l.identities))
	copy(identitiesCopy, l.identities)

	identitiesCopy = append(identitiesCopy, identity)

	logger := &Logger{
		traceID:    l.traceID,
		identities: identitiesCopy,
	}
	logger.App = l.cloneSet(l.App, logger)

	return logger
}

func (l *Logger) Close() {
	l.App.close()
}

func (l *Logger) cloneSet(set *loggerSet, parent *Logger) *loggerSet {
	return &loggerSet{
		debugLog: set.debugLog,
		errorLog: set.errorLog,
		parent:   parent,
	}
}

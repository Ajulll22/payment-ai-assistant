package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"

	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
)

func InjectIntoContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, constant.LoggerKey, l)
}

func FromContext(ctx context.Context, cfg constant.LogConfig) *Logger {
	if l, ok := ctx.Value(constant.LoggerKey).(*Logger); ok {
		return l
	}

	dummy := NewLogger(cfg)
	return dummy
}

func cleanupOldLogs(log *Logger, cfgLog constant.LogConfig) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {

		defer wg.Done()
		defer func() {

			if err := recover(); err != nil {
				// fmt.Println("Send mail error")
				err_msg := fmt.Sprintf("%v", err)
				log.App.Error("Delete log app file error, msg:" + err_msg)
			}

		}()

		DeleteExcessLogs(cfgLog, LogTypeApp)

	}()

	wg.Wait()
}

func DeleteExcessLogs(cfgLog constant.LogConfig, logType LogType) error {

	var logDebugList []string
	var logErrorList []string

	reDebug := regexp.MustCompile(fmt.Sprintf(`^%s_(\d{4}-\d{2}-\d{2})$`, regexp.QuoteMeta(cfgLog.DebugFilename)))
	reError := regexp.MustCompile(fmt.Sprintf(`^%s_(\d{4}-\d{2}-\d{2})$`, regexp.QuoteMeta(cfgLog.ErrorFilename)))

	dir := filepath.Join(cfgLog.Path)

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		switch {
		case reDebug.MatchString(file.Name()):
			logDebugList = append(logDebugList, file.Name())
		case reError.MatchString(file.Name()):
			logErrorList = append(logErrorList, file.Name())
		}

	}

	sort.Slice(logDebugList, func(i, j int) bool {
		return logDebugList[i] > logDebugList[j]
	})
	for key, value := range logDebugList {
		if key > cfgLog.MaxFile-1 {
			os.Remove(dir + value)
		}
	}

	sort.Slice(logErrorList, func(i, j int) bool {
		return logErrorList[i] > logErrorList[j]
	})
	for key, value := range logErrorList {
		if key > cfgLog.MaxFile-1 {
			os.Remove(dir + value)
		}
	}

	return nil

}

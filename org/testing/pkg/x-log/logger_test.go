package xlog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestLogger(t *testing.T) {
	type args struct {
		count        int
		loggerType   LoggerType
		startLogging func(i int)
	}

	tests := []struct {
		name         string
		args         args
		wanttedLines []string
	}{
		{
			name: "Test logging with Logrus",
			args: args{
				count:      10,
				loggerType: LoggerTypes.Logrus(),
				startLogging: func(i int) {
					msg := fmt.Sprintf("hello begin [%d]", i)
					WithFields(Fields{"status": "success"}).WithFields(Fields{"status": "failed"}).Info(msg)
					msg = fmt.Sprintf("hello, processing the hello message [%d]", i)
					Info(msg)
					WithFields(Fields{"status": "failed", "value": 6.9}).Errorf("failed to process value [%v] [%d]", 6.9, i)
					msg = fmt.Sprintf("Failed to process value 6.9 [%d]", i)
					Error(msg)
					msg = fmt.Sprintf("error occured [%d]", i)
					Warn(msg)
				},
			},
			wanttedLines: []string{
				`level=info msg="hello begin [%d]" status=failed`,
				`level=info msg="hello, processing the hello message [%d]"`,
				`level=error msg="failed to process value [6.9] [%d]" status=failed value=6.9`,
				`level=error msg="Failed to process value 6.9 [%d]"`,
				`level=warning msg="error occured [%d]"`,
			},
		},
		{
			name: "Test logging with slog",
			args: args{
				count:      10,
				loggerType: LoggerTypes.SLog(),
				startLogging: func(i int) {
					msg := fmt.Sprintf("hello begin [%d]", i)
					WithFields(Fields{"status": "success"}).WithFields(Fields{"status": "failed"}).Info(msg)
					msg = fmt.Sprintf("hello, processing the hello message [%d]", i)
					Info(msg)
					WithFields(Fields{"value": 6.9}).Errorf("failed to process value [%v] [%d]", 6.9, i)
					msg = fmt.Sprintf("Failed to process value 6.9 [%d]", i)
					Error(msg)
					msg = fmt.Sprintf("error occured [%d]", i)
					Warn(msg)
				},
			},
			wanttedLines: []string{
				`"level":"INFO","msg":"hello begin [%d]","status":"failed"`,
				`"level":"INFO","msg":"hello, processing the hello message [%d]"`,
				`"level":"ERROR","msg":"failed to process value [6.9] [%d]","value":6.9`,
				`"level":"ERROR","msg":"Failed to process value 6.9 [%d]"`,
				`"level":"WARN","msg":"error occured [%d]"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := fmt.Sprintf("testing-logs-%s", tt.name)
			fileWriter := &lumberjack.Logger{
				Filename:   filePath,
				MaxSize:    1, // megabytes
				MaxBackups: 1,
				MaxAge:     1, // days
			}
			logWriter := io.MultiWriter(os.Stdout, fileWriter)
			logger, _ := NewLogger(logWriter, DebugLevel, tt.args.loggerType)
			SetDefaultLogger(logger)
			var wg sync.WaitGroup
			for i := 0; i < tt.args.count; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					tt.args.startLogging(i)
				}(i)
			}
			wg.Wait()

			err := fileWriter.Close()
			require.Nil(t, err)

			readLines, err := readLines(filePath)
			require.Nil(t, err)
			for _, want := range tt.wanttedLines {
				for i := 0; i < tt.args.count; i++ {
					line := fmt.Sprintf(want, i)
					found := false
					for _, readLine := range readLines {
						if strings.Contains(readLine, line) {
							found = true
							break
						}
					}
					if !found {
						assert.Fail(t, fmt.Sprintf("log message [%s] not found", line))
					}
				}
			}

			files, err := filepath.Glob(fmt.Sprintf("%s*", filePath))
			if err != nil {
				assert.Fail(t, "failed to list files")
			}
			for _, f := range files {
				if err := os.Remove(f); err != nil {
					assert.Fail(t, "failed to delete file")
				}
			}
		})
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err =  file.Close()
	if scanner.Err() != nil {
		return lines, scanner.Err()
	}
	if err != nil {
		return lines, err
	}
	return lines, nil
}

func TestLoggerFileRotation(t *testing.T) {
	type args struct {
		count        int
		startLogging func(i int)
	}

	tests := []struct {
		name         string
		args         args
		wanttedLines []string
	}{
		{
			name: "Stress Testing to check file rotation",
			args: args{
				count: 9999,
				startLogging: func(i int) {
					msg := fmt.Sprintf("hello begin [%d]", i)
					WithFields(Fields{"status": "success"}).Info(msg)
					msg = fmt.Sprintf("hello, processing the hello message [%d]", i)
					Info(msg)
					WithFields(Fields{"status": "failed", "value": 6.9}).Errorf("failed to process value [%v] [%d]", 6.9, i)
					msg = fmt.Sprintf("Failed to process value 6.9 [%d]", i)
					Error(msg)
					msg = fmt.Sprintf("error occured [%d]", i)
					Warn(msg)
				},
			},
			wanttedLines: []string{
				`level=info msg="hello begin [%d]" status=success`,
				`level=info msg="hello, processing the hello message [%d]"`,
				`level=error msg="failed to process value [6.9] [%d]" status=failed value=6.9`,
				`level=error msg="Failed to process value 6.9 [%d]"`,
				`level=warning msg="error occured [%d]"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := fmt.Sprintf("testing-logs-%s", tt.name)
			fileWriter := &lumberjack.Logger{
				Filename:   filePath,
				MaxSize:    1, // megabytes
				MaxBackups: 1,
				MaxAge:     1, // days
			}
			logWriter := io.MultiWriter(os.Stdout, fileWriter)
			logger, _ := NewLogger(logWriter, DebugLevel, LoggerTypes.Logrus())
			SetDefaultLogger(logger)
			var wg sync.WaitGroup
			for i := 0; i < tt.args.count; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					tt.args.startLogging(i)
				}(i)
			}
			wg.Wait()

			err := fileWriter.Close()
			require.Nil(t, err)

			files, err := filepath.Glob(fmt.Sprintf("%s*", filePath))
			if err != nil {
				assert.Fail(t, "failed to list files")
			}
			for _, f := range files {
				if err := os.Remove(f); err != nil {
					assert.Fail(t, "failed to delete file")
				}
			}
		})
	}
}

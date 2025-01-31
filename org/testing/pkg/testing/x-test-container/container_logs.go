package xtestcontainer

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	"github.com/testcontainers/testcontainers-go"
)

type ContainerLogConsumer struct {
	ServiceName models.ServiceName
	logs        []testcontainers.Log
	errorLogs   []testcontainers.Log
	fdAll       io.WriteCloser
	fdErr       io.WriteCloser
}

func NewContainerLogConsumer(containerLogFolder string, serviceName models.ServiceName) (*ContainerLogConsumer, error) {
	fullLogFileName := fmt.Sprintf("%s.full.log", serviceName)
	errLogFileName := fmt.Sprintf("%s.error.log", serviceName)
	logFilePath := path.Join(containerLogFolder, fullLogFileName)
	errLogFilePath := path.Join(containerLogFolder, errLogFileName)
	fdAll, err := os.Create(logFilePath)
	if err != nil {
		return nil, err
	}
	fdErr, err := os.Create(errLogFilePath)
	if err != nil {
		errClose := fdAll.Close()
		return nil, fx.FlattenErrorsIfAny(err, errClose)
	}

	return &ContainerLogConsumer{
		ServiceName: serviceName,
		fdAll:       fdAll,
		fdErr:       fdErr,
	}, nil
}

func (t *ContainerLogConsumer) Teardown(ctx context.Context) error {
	err1 := t.fdAll.Close()
	err2 := t.fdErr.Close()
	return fx.FlattenErrorsIfAny(err1, err2)
}

func (t *ContainerLogConsumer) Accept(l testcontainers.Log) {
	line := []byte(fmt.Sprintf("[%s] %s", l.LogType, string(l.Content)))
	if t.fdAll != nil {
		_, _ = t.fdAll.Write(line)
	}
	if strings.Contains(string(l.Content), "ERROR") ||
		strings.Contains(string(l.Content), "level=error") ||
		strings.Contains(string(l.Content), "level=fatal") {
		t.errorLogs = append(t.errorLogs, l)
		if t.fdErr != nil {
			_, _ = t.fdErr.Write(line)
		}
	}
	t.logs = append(t.logs, l)
}

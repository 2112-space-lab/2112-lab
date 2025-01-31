package xtestlog

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/org/2112-space-lab/org/testing/pkg/fx"
)

type RequestClientLogger struct {
	httpClient *http.Client
	logger     *slog.Logger
}

func NewRequestClientLogger(cl *http.Client, logger *slog.Logger) *RequestClientLogger {
	return &RequestClientLogger{
		httpClient: cl,
		logger:     logger,
	}
}

func (cl *RequestClientLogger) Do(req *http.Request) (*http.Response, error) {
	logger := cl.logger.With(
		slog.String("requestUrl", req.URL.String()),
		slog.Any("requestMethod", req.Method),
		slog.Any("requestPayload", req.Body),
		slog.Any("requestHeader", req.Header),
	)
	res, err := cl.httpClient.Do(req)
	if res != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return res, err
		}
		defer func() {
			err = fx.FlattenErrorsIfAny(err, res.Body.Close())
		}()
		res.Body = io.NopCloser(bytes.NewReader(body))
		logger = logger.With(
			slog.Any("responseStatus", res.Status),
			slog.Any("responsePayload", json.RawMessage(body)),
			slog.Any("responseHeaders", res.Header),
		)
	}
	if err != nil {
		logger.Error("request failed",
			slog.Any("error", err),
		)
	} else {
		logger.Info("request success")
	}
	return res, err
}

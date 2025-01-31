package models

import (
	"fmt"
	"time"

	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	xtime "github.com/org/2112-space-lab/org/testing/pkg/x-time"
)

type TimeCheckpointName string
type TimeCheckpointValue time.Time

type TimeCheckpointExpression string // name>offset>name

type BackgroundOperationCompoundKey string // checkpointExpr_key

type BackgroundOperation struct {
	BaseKey         string
	Checkpoint      TimeCheckpointExpression
	RegisteredAt    xtime.UtcTime
	StartedAt       fx.Option[xtime.UtcTime]
	CompletedAt     fx.Option[xtime.UtcTime]
	CompletionError error
}

func NewBackgroundOperation(
	baseKey string,
	checkpoint TimeCheckpointExpression,
) *BackgroundOperation {
	return &BackgroundOperation{
		BaseKey:      baseKey,
		Checkpoint:   checkpoint,
		RegisteredAt: xtime.UtcNow(),
	}
}

func (o BackgroundOperation) GetKey() BackgroundOperationCompoundKey {
	key := fmt.Sprintf("%s_%s", o.Checkpoint, o.BaseKey)
	return BackgroundOperationCompoundKey(key)
}

func (o *BackgroundOperation) Start() {
	o.StartedAt = fx.NewValueOption(xtime.UtcNow())
}

func (o *BackgroundOperation) Complete(err error) {
	o.CompletedAt = fx.NewValueOption(xtime.UtcNow())
	o.CompletionError = err
}

package testservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"slices"
	"time"

	"github.com/blues/jsonata-go"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	testservicecontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service-container"
	"github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtesttime "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time"
	models_time "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time/models"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/encoding/protojson"
)

var protojsonOpts = protojson.MarshalOptions{
	AllowPartial:      true,
	UseProtoNames:     true,
	UseEnumNumbers:    false,
	EmitUnpopulated:   true,
	EmitDefaultValues: true,
	Multiline:         true,
}

// SubscribeToPropagator subscribes to RabbitMQ for Propagator events
func SubscribeToPropagator(ctx context.Context, scenarioState PropagatorClientScenarioState, service string, subscriber string, callbacks []models.EventCallbackInfo) (context.CancelFunc, error) {
	conn, err := amqp.Dial(testservicecontainer.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	queueName := fmt.Sprintf("propagator_events_%s", subscriber)
	_, err = ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		queueName, // queue name
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to consume messages: %v", err)
	}

	streamCtx, cancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-streamCtx.Done():
				ch.Close()
				conn.Close()
				return
			case msg := <-msgs:
				var event models.EventRoot
				err := json.Unmarshal(msg.Body, &event)
				if err != nil {
					log.Printf("Error unmarshaling event: %v", err)
					continue
				}

				log.Printf("Received event: %+v", event)

				for _, cb := range callbacks {
					if cb.EventType != event.EventType {
						continue
					}

					go func(cb models.EventCallbackInfo) {
						if cb.ActionDelay != "" {
							waitDur, _ := time.ParseDuration(cb.ActionDelay)
							time.Sleep(waitDur)
						}

						err := processEventCallback(streamCtx, scenarioState, service, cb, event)
						if err != nil {
							log.Printf("Error processing callback for event %s: %v", event.EventType, err)
						}
					}(cb)
				}
			}
		}
	}()

	return cancel, nil
}

// processEventCallback executes the callback for an event
func processEventCallback(ctx context.Context, scenarioState PropagatorClientScenarioState, serviceName string, cb models.EventCallbackInfo, event models_service.EventRoot) error {
	log.Printf("Processing callback: %+v for event: %+v", cb, event)
	scenarioState.SaveReceivedEvent(&event, models_service.ServiceName(serviceName))
	return nil
}

// VerifyPropagatorEvents checks if expected events have been received
func VerifyPropagatorEvents(scenarioState PropagatorClientScenarioState, serviceName string, expectedEvents []models.ExpectedEvent) error {
	for _, exp := range expectedEvents {
		err := checkExpectedEvent(scenarioState, models.ServiceName(serviceName), exp)
		if err != nil {
			return err
		}
	}
	return nil
}

// checkExpectedEvent verifies if a single expected event was received
func checkExpectedEvent(scenarioState PropagatorClientScenarioState, serviceName models.ServiceName, expected models.ExpectedEvent) error {
	var expFrom, expToWarn, expToErr models_time.TimeCheckpointValue
	var errExpFrom, errExpToWarn, errExpToErr error
	expFrom, errExpFrom = xtesttime.EvaluateCheckpoint(scenarioState, expected.FromTime)
	expToErr, errExpToErr = xtesttime.EvaluateCheckpoint(scenarioState, expected.ToTimeErr)
	if expected.ToTimeWarn != "" {
		expToWarn, errExpToWarn = xtesttime.EvaluateCheckpoint(scenarioState, expected.ToTimeWarn)
	} else {
		expToWarn = expToErr
	}

	err := fx.FlattenErrorsIfAny(errExpFrom, errExpToWarn, errExpToErr)
	if err != nil {
		return err
	}

	logger := scenarioState.GetLogger().With(
		slog.Group("expectedEventInfo",
			slog.Any("expectedEvent", expected),
			slog.String("from", time.Time(expFrom).Format(time.RFC3339)),
			slog.String("toWarn", time.Time(expToWarn).Format(time.RFC3339)),
			slog.String("toErr", time.Time(expToErr).Format(time.RFC3339)),
		),
	)

	pauseIncrementDuration := 100 * time.Millisecond
	now := time.Now().UTC()
	if expected.Occurence < 0 {
		sleepDuration := time.Time(expToErr).Sub(time.Now().UTC())
		logger.Info("negative occurance wanting last received - sleeping until end of expected window",
			slog.Duration("sleepDuration", sleepDuration),
			slog.Any("expectedEvent", expected),
		)
		time.Sleep(sleepDuration)
	} else if now.Before(time.Time(expToWarn).Add(pauseIncrementDuration)) {
		sleepDuration := time.Time(expToWarn).Sub(time.Now().UTC())
		logger.Info("sleeping until minimal expected event received time",
			slog.Duration("sleepDuration", sleepDuration),
			slog.Any("expectedEvent", expected),
		)
		time.Sleep(sleepDuration)
	}
	var evt models_service.EventRoot
	var found bool
	for {
		events := scenarioState.GetReceivedEvents(serviceName, time.Time(expFrom).Add(-1*time.Millisecond), time.Time(expToErr))
		evt, found = containsExpectedEvent(logger, events, expected)
		if found {
			break
		}
		if now.After(time.Time(expToErr).Add(pauseIncrementDuration)) {
			break
		}
		time.Sleep(pauseIncrementDuration)
		now = time.Now().UTC()
	}
	if !found || evt.EventTimeUtc.After(time.Time(expToErr)) {
		if expected.IsReject {
			logger.Debug("unwanted event not received during time window - all is fine")
			return nil
		}
		return fmt.Errorf("expected event [%+v] not received between [%s] and [%s]",
			expected,
			time.Time(expFrom).Format(time.RFC3339),
			time.Time(expToErr).Format(time.RFC3339),
		)
	}
	if expected.IsReject {
		return fmt.Errorf("unwanted event with UID [%s] matching [%+v] received between [%s] and [%s]",
			evt.EventUid,
			expected,
			time.Time(expFrom).Format(time.RFC3339),
			time.Time(expToErr).Format(time.RFC3339),
		)
	}

	if evt.EventTimeUtc.After(time.Time(expToWarn)) {
		logger.Warn("got expected event in warning threshold",
			slog.Any("receivedEvent", evt),
		)
	} else {
		logger.Info("got expected event",
			slog.Any("receivedEvent", evt),
		)
	}
	if expected.AssignRef != "" {
		mval, err := json.Marshal(evt)
		if err != nil {
			logger.With(slog.Any("error", err.Error())).Error("failed to marshal event - cannot assign named ref",
				slog.Any("event", evt),
			)
			return err
		}
		scenarioState.RegisterNamedEventReference(expected.AssignRef, mval)
	}
	if expected.ProduceCheckpointEventTime != "" {
		err = scenarioState.RegisterCheckpoint(models_time.TimeCheckpointName(expected.ProduceCheckpointEventTime), models_time.TimeCheckpointValue(evt.EventTimeUtc))
		if err != nil {
			return err
		}
	}
	return nil
}

func containsExpectedEvent(logger *slog.Logger, events []models_service.EventRoot, expected models.ExpectedEvent) (models_service.EventRoot, bool) {
	occurence := 0
	logger.Debug("looking for event",
		slog.Any("expected", expected),
		slog.Any("events", events),
	)
	ordered := events
	expectedOccurence := expected.Occurence
	if expected.Occurence < 0 {
		clone := slices.Clone(events)
		slices.Reverse(clone)
		expectedOccurence = -expected.Occurence
		ordered = clone
	}
	for _, evt := range ordered {
		if expected.EventType != evt.EventType {
			continue
		}
		if expected.XPathQuery != "" {
			mval, err := json.Marshal(evt)
			if err != nil {
				logger.With(slog.Any("error", err.Error())).Error("failed to marshal event - cannot evaluate json query",
					slog.String("eventType", evt.EventType),
					slog.String("eventUID", evt.EventUid),
					slog.Any("expectedEvent", expected),
				)
				continue
			}
			e := jsonata.MustCompile(expected.XPathQuery)
			res, err := e.EvalBytes(mval)
			if err != nil {
				logger.With(slog.Any("error", err.Error())).Error("failed to evaluate json query on event",
					slog.Any("event", string(mval)),
					slog.Any("expectedEvent", expected),
				)
				continue
			}
			if string(sanitizeURLEncodedJSON(res)) != expected.XPathValue {
				logger.Debug("jsonQuery value mismatch from expected - moving on next received event",
					slog.Any("expectedEvent", expected),
					slog.Any("event", string(mval)),
					slog.String("actualJsonQueryValue", string(res)),
				)
				continue
			}
		}
		occurence++
		if occurence == expectedOccurence {
			return evt, true
		}
	}
	return models_service.EventRoot{}, false
}

func sanitizeURLEncodedJSON(b []byte) []byte {
	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	return b
}

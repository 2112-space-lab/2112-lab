package testservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/blues/jsonata-go"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	testservicecontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service-container"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	xtesttime "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time"
	models_time "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time/models"
	"github.com/streadway/amqp"
)

// Subscribe dynamically subscribes to all queues in RabbitMQ
func Subscribe(ctx context.Context, scenarioState RabbitMqClientScenarioState, service models_service.ServiceName, queueName string, callbacks []models_service.EventCallbackInfo) (context.CancelFunc, error) {
	conn, err := amqp.Dial(testservicecontainer.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	queues, err := getRabbitMQQueues(queueName)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to fetch queue list: %v", err)
	}

	for _, queue := range queues {
		log.Printf("🔄 Found queue: %s, subscribing...", queue)

		for _, cb := range callbacks {
			routingKey := fmt.Sprintf("events.%s", cb.EventType)
			err := ch.QueueBind(
				queue,
				routingKey,
				"",
				false,
				nil,
			)
			if err != nil {
				ch.Close()
				conn.Close()
				return nil, fmt.Errorf("failed to bind queue %s to routing key %s: %v", queue, routingKey, err)
			}
			log.Printf("✅ Subscribed queue '%s' to routing key: %s", queue, routingKey)
		}

		msgs, err := ch.Consume(
			queue,
			queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to consume messages from queue %s: %v", queue, err)
		}

		streamCtx, cancel := context.WithCancel(ctx)

		go func(queueName string) {
			defer func() {
				cancel()
				ch.Close()
				conn.Close()
				log.Printf("🔴 Gracefully stopped listening to queue: %s", queueName)
			}()

			for {
				select {
				case <-streamCtx.Done():
					log.Printf("🛑 Stopping listener for queue: %s", queueName)
					return
				case msg := <-msgs:
					var event models_service.EventRoot
					err := json.Unmarshal(msg.Body, &event)
					if err != nil {
						log.Printf("❌ Error unmarshaling event from queue %s: %v", queueName, err)
						continue
					}

					log.Printf("📥 Received event from queue %s: %+v", queueName, event)
					scenarioState.SaveReceivedEvent(&event, service)
					for _, cb := range callbacks {
						if cb.EventType != event.EventType {
							continue
						}

						go func(cb models_service.EventCallbackInfo) {
							if cb.ActionDelay != "" {
								waitDur, _ := time.ParseDuration(cb.ActionDelay)
								time.Sleep(waitDur)
							}

							scenarioState.SaveReceivedEvent(&event, service)
							if err != nil {
								log.Printf("❌ Error processing callback for event %s in queue %s: %v", event.EventType, queueName, err)
							}
						}(cb)
					}
				}
			}
		}(queue)
	}

	return nil, nil
}

// Fetch all queues from RabbitMQ Management API
func getRabbitMQQueues(filter string) ([]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", testservicecontainer.RabbitMQAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.SetBasicAuth(testservicecontainer.RabbitMQAPIUser, testservicecontainer.RabbitMQAPIPassword)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call RabbitMQ API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response from RabbitMQ API: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var queues []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &queues); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	log.Printf("Total queues retrieved from RabbitMQ: %d", len(queues))

	var queueNames []string
	for _, q := range queues {
		trimmedLowerQueue := strings.ToLower(strings.TrimSpace(q.Name))
		trimmedLowerFilter := strings.ToLower(strings.TrimSpace(filter))

		if strings.Contains(trimmedLowerQueue, trimmedLowerFilter) {
			queueNames = append(queueNames, q.Name)
			log.Printf("Queue matched filter '%s': %s", filter, q.Name)
		}
	}

	if len(queueNames) == 0 {
		log.Printf("No queues found matching filter: '%s'", filter)
	} else {
		log.Printf("Total queues matched filter '%s': %d", filter, len(queueNames))
	}

	return queueNames, nil
}

// VerifyEvents checks if expected events have been received
func VerifyEvents(scenarioState RabbitMqClientScenarioState, serviceName string, expectedEvents []models_service.ExpectedEvent) error {
	for _, exp := range expectedEvents {
		err := checkExpectedEvent(scenarioState, models_service.ServiceName(serviceName), exp)
		if err != nil {
			return err
		}
	}
	return nil
}

// checkExpectedEvent verifies if a single expected event was received
func checkExpectedEvent(scenarioState RabbitMqClientScenarioState, serviceName models_service.ServiceName, expected models_service.ExpectedEvent) error {
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
	if expected.Occurrence < 0 {
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
	if !found || evt.GetEventTimeUtc().Inner().After(time.Time(expToErr)) {
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

	if evt.GetEventTimeUtc().Inner().After(time.Time(expToWarn)) {
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
		err = scenarioState.RegisterCheckpoint(models_time.TimeCheckpointName(expected.ProduceCheckpointEventTime), models_time.TimeCheckpointValue(evt.GetEventTimeUtc().Inner()))
		if err != nil {
			return err
		}
	}
	return nil
}

func containsExpectedEvent(logger *slog.Logger, events []models_service.EventRoot, expected models_service.ExpectedEvent) (models_service.EventRoot, bool) {
	occurence := 0
	logger.Debug("looking for event",
		slog.Any("expected", expected),
		slog.Any("events", events),
	)
	ordered := events
	expectedOccurence := expected.Occurrence
	if expected.Occurrence < 0 {
		clone := slices.Clone(events)
		slices.Reverse(clone)
		expectedOccurence = -expected.Occurrence
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

// PublishTestEvent injects a test event into a RabbitMQ queue for testing purposes.
func PublishTestEvent(serviceName models_service.ServiceAppName, queueName string, event models_service.EventRoot) error {
	conn, err := amqp.Dial(testservicecontainer.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("❌ failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("❌ failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Serialize the event
	eventBody, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("❌ failed to serialize test event: %v", err)
	}

	// Ensure queue exists
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("❌ failed to declare queue %s: %v", queueName, err)
	}

	// Publish message to queue
	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventBody,
		},
	)
	if err != nil {
		return fmt.Errorf("❌ failed to publish test event: %v", err)
	}

	log.Printf("✅ Successfully injected test event into queue '%s': %+v\n", queueName, event)
	return nil
}

// RegisterQueuesToRabbitMQ declares and ensures the existence of queues in RabbitMQ.
func RegisterQueuesToRabbitMQ(queues models_cont.EnvVarKeyValueMap) error {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(testservicecontainer.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("❌ failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("❌ failed to open a channel: %v", err)
	}
	defer ch.Close()

	for key, queueName := range queues {
		if queueName == "" {
			log.Printf("⚠️ Skipping empty queue name for key: %s", key)
			continue
		}

		// Declare the queue
		_, err := ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			log.Printf("❌ Failed to declare queue '%s' (key: %s): %v", queueName, key, err)
			continue
		}

		log.Printf("✅ Successfully registered queue '%s' (key: %s)", queueName, key)
	}

	return nil
}

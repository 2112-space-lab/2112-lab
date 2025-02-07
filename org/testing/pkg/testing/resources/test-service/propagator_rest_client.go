package testservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	testservicecontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service-container"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtestlog "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-log"
)

type PropagatorRestClientRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type PropagatorRestClient struct {
	hostName   string
	httpClient PropagatorRestClientRequestDoer
}

func GetPropagatorRestClient(ctx context.Context, scenarioState PropagatorClientScenarioState, serviceName models_service.ServiceName) (*PropagatorRestClient, error) {
	logger := scenarioState.GetLogger()
	cont, err := scenarioState.GetPropagatorServiceContainer(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	httpBoundPort, err := cont.GetBoundPort(ctx, testservicecontainer.PropagatorServiceHttpPort)
	if err != nil {
		return nil, err
	}
	hostname, err := cont.GetHostName(ctx)
	if err != nil {
		return nil, err
	}

	restClient := &PropagatorRestClient{
		hostName:   fmt.Sprintf("http://%s:%s", hostname, httpBoundPort.Port()),
		httpClient: xtestlog.NewRequestClientLogger(&http.Client{}, logger),
	}
	return restClient, err
}

func (cl *PropagatorRestClient) RequestSatellitePropagation(ctx context.Context, config models_service.PropagatorSettings) error {
	url := fmt.Sprintf("%s/satellite/propagate", cl.hostName)

	buf, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal UpdateLiveSettings [%w]", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("failed prepare request UpdateLiveSettings [%w]", err)
	}

	res, err := cl.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed send request UpdateLiveSettings [%w]", err)
	}
	defer func() {
		err = fx.FlattenErrorsIfAny(err, res.Body.Close())
	}()
	_, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body UpdateLiveSettings [%w]", err)
	}
	return nil
}

package testservice

import (
	"fmt"
	"time"

	xtime "github.com/org/2112-space-lab/org/testing/pkg/x-time"
)

type GlobalPropKey string

type AppGlobalPropertyUpdateTableRow struct {
	PropertyKey GlobalPropKey
	Value       string
}

// func GetAppRestClient(ctx context.Context, scenarioState AppClientScenarioState, serviceName models_service.ServiceName) (*app_client.ClientWithResponses, error) {
// 	logger := scenarioState.GetLogger()
// 	cont, err := scenarioState.GetAppServiceContainer(ctx, serviceName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	httpBoundPort, err := cont.GetBoundPort(ctx, testservicecontainer.AppServiceHttpPort)
// 	if err != nil {
// 		return nil, err
// 	}
// 	hostname, err := cont.GetHostName(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	restClient, err := app_client.NewClientWithResponses(
// 		fmt.Sprintf("http://%s:%s", hostname, httpBoundPort.Port()),
// 		app_client.WithHTTPClient(xtestlog.NewRequestClientLogger(&http.Client{}, logger)),
// 	)
// 	return restClient, err
// }

func MapToOptionalTime(value string) *time.Time {
	if value != "" {
		t, err := xtime.FromString(xtime.DateTimeFormat(value))
		if err != nil {
			return nil
		}
		ti := t.Inner()
		return &ti
	}
	return nil
}

func MapToOptionalString(value string) *string {
	if value != "" {
		t := value
		str := fmt.Sprintf("%v", t)
		return &str
	}
	return nil
}

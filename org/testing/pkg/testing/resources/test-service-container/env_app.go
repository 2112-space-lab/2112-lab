package testservicecontainer

import "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"

func GetAppDefaultEnv() models.EnvVarKeyValueMap {
	return map[string]string{
		"SITE_ID":                          "GS_US_SUN",
		"PASSCALC_HOST":                    "http://passcalc:5007",
		"LOG_LEVEL":                        "debug",
		"GEM_TLS":                          "false",
		"GEM_TYPE":                         "GEM",
		"GATEWAY_MODE":                     "MEO",
		"GATEWAY_TYPE":                     "ttnc",
		"GATEWAY_HEIGHT":                   "142",
		"GATEWAY_LATITUDE":                 "21.671",
		"GATEWAY_LONGITUDE":                "-158.031",
		"GATEWAY_CUTOFF":                   "5",
		"BRIDGE_STREAM_COOLDOWN_MS":        "12000",
		"BRIDGE_MINIMAL_UPTIME":            "1000",
		"SERVER_TLS":                       "false",
		"PLAN_MONITOR_TICKER_TIME_SPAN_MS": "10",
		"TRACK_RETRY_MAXIMUM_WAIT_MS":      "100",
		"TRACK_RETRY_BACKOFF_BASELINE":     "20",
		"PAST_EVENT_MAX_BUFFER_STREAM_V2":  "2000",
	}
}

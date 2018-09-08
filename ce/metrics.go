package ce

import (
	"fmt"
	"log"
)

const (
	MetricsAPI                     = "/metrics/api"
	MetricsBulkJobsAPI             = "/metrics/bulk-jobs"
	MetricsElementsCreated         = "/metrics/elements-created"
	MetricsElementInstancesCreated = "/metrics/element-instances-created"
	MetricsEvents                  = "/metrics/events"
	MetricsFormulaExecutions       = "/metrics/formula-executions"
	MetricsFormulasCreated         = "/metrics/formulas-created"
	MetricsVDRsCreated             = "/metrics/vdrs-created"
	MetricsVDRsInvoked             = "/metrics/vdrs-invoked"
	MetricsHubAPI                  = "/metrics/hub-api"
	MetricsHubsCreated             = "/metrics/hubs-created"
)

// GetJSONMetricsFor provides JSON return for the provided url
func GetJSONMetricsFor(url string, base, auth string, debug bool) ([]byte, int, string, error) {
	if debug {
		log.Println("GET", url)
	}
	bodybytes, status, curlcmd, err := Execute("GET", url, auth)
	if debug {
		log.Printf("Status %v", status)
	}
	if err != nil {
		if debug {
			log.Printf("%s", bodybytes)
		}
		return bodybytes, status, curlcmd, err
	}
	if status != 200 {
		return bodybytes, status, curlcmd, fmt.Errorf("Status code %v", status)
	}
	return bodybytes, status, curlcmd, nil
}

// GetMetricsHubAPI returns raw JSON metrics
func GetMetricsHubAPI(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s", base, MetricsHubAPI)
	return GetJSONMetricsFor(url, base, auth, debug)
}

// GetMetricsHubsCreated returns raw JSON metrics
func GetMetricsHubsCreated(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s", base, MetricsHubsCreated)
	return GetJSONMetricsFor(url, base, auth, debug)
}

// GetMetricsVDRsInvoked returns raw JSON metrics
func GetMetricsVDRsInvoked(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsVDRsInvoked,
	)
	return GetJSONMetricsFor(url, base, auth, debug)
}

// GetMetricsVDRsCreated returns raw JSON metrics
func GetMetricsVDRsCreated(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsVDRsCreated,
	)
	return GetJSONMetricsFor(url, base, auth, debug)
}

// GetMetricsFormulasCreated returns raw JSON metrics
func GetMetricsFormulasCreated(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsFormulasCreated,
	)
	return GetJSONMetricsFor(url, base, auth, debug)
}

// GetMetricsFormulaExecutions returns raw JSON metrics
func GetMetricsFormulaExecutions(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsFormulaExecutions,
	)
	return GetJSONMetricsFor(url, base, auth, debug)
}

// GetMetricsEvents returns raw JSON metrics
func GetMetricsEvents(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsEvents,
	)
	return GetJSONMetricsFor(url, base, auth, debug)
}

// GetMetricsElementsCreated returns raw JSON metrics
func GetMetricsElementsCreated(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsElementsCreated,
	)
	return GetJSONMetricsFor(url, base, auth, debug)
}

// GetMetricsElementInstancesCreated returns raw JSON metrics
func GetMetricsElementInstancesCreated(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s", base, MetricsElementInstancesCreated)
	return GetJSONMetricsFor(url, base, auth, debug)
}

// GetMetricsBulkJobs returns raw JSON metrics
func GetMetricsBulkJobs(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsBulkJobsAPI,
	)
	return GetJSONMetricsFor(url, base, auth, debug)
}

// GetMetrics returns raw JSON metrics
func GetMetrics(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsAPI,
	)
	return GetJSONMetricsFor(url, base, auth, debug)
}

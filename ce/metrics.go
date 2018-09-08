package ce

import (
	"fmt"
	"log"
)

const (
	MetricsAPI                   = "/metrics/api"
	MetricsBulkJobs              = "/metrics/bulk-jobs"
	MetricsElementInstaceCreated = "/metrics/element-instance-created"
	MetricsElementsCreated       = "/metrics/elements-created"
	MetricsEvents                = "/metrics/events"
	MetricsFormulaExecutions     = "/metrics/formula-executions"
	MetricsFormulaCreated        = "/metrics/formulas-created"
	MetricsHubAPI                = "/metrics/hub-api"
	MetricsHubsCreated           = "/metrics/hubs-created"
	MetricsVDRsCreated           = "/metrics/vdrs-created"
	MetricsVDRsInvoked           = "/metrics/vdrs-invoked"
)

// GetMetricsAPI Retrieve the API metrics for the accounts provided.
// Any customer or organization IDs provided will be used to identify accounts within those entities.
func GetMetricsAPI(base, auth string, instanceID string, enable bool, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsAPI,
	)
	if debug {
		log.Println("Getting metrics for APIs...")
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

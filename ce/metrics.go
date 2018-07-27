package ce

import (
	"fmt"
	"log"
)

const (
	MetricsAPI               = "/metrics/api"
	MetricsBulkJobsAPI       = "/metrics/bulk-jobs"
	MetricsElementsCreated   = "/metrics/elements-created"
	MetricsEvents            = "/metrics/events"
	MetricsFormulaExecutions = "/metrics/formula-executions"
	MetricsFormulasCreated   = "/metrics/formulas-created"
	MetricsVDRsCreated       = "/metrics/vdrs-created"
	MetricsVDRsInvoked       = "/metrics/vdrs-invoked"
)

// GetMetricsVDRsInvoked returns raw JSON metrics
func GetMetricsVDRsInvoked(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsVDRsInvoked,
	)
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

// GetMetricsVDRsCreated returns raw JSON metrics
func GetMetricsVDRsCreated(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsVDRsCreated,
	)
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

// GetMetricsFormulasCreated returns raw JSON metrics
func GetMetricsFormulasCreated(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsFormulasCreated,
	)
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

// GetMetricsFormulaExecutions returns raw JSON metrics
func GetMetricsFormulaExecutions(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsFormulaExecutions,
	)
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

// GetMetricsEvents returns raw JSON metrics
func GetMetricsEvents(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsEvents,
	)
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

// GetMetricsElementsCreated returns raw JSON metrics
func GetMetricsElementsCreated(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsElementsCreated,
	)
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

// GetMetricsBulkJobs returns raw JSON metrics
func GetMetricsBulkJobs(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsBulkJobsAPI,
	)
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

// GetMetrics returns raw JSON metrics
func GetMetrics(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		MetricsAPI,
	)
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

package ce

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/moul/http2curl"
)

// Job represents an scheduled job on the platform
type Job struct {
	ID                 string     `json:"id"`
	DisallowConcurrent bool       `json:"disallowConcurrent"`
	Data               JobData    `json:"data"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Trigger            JobTrigger `json:"trigger"`
}

// JobData represents the data of the scheduled job
type JobData struct {
	ID            int         `json:"id"`
	ElementKey    string      `json:"elementKey"`
	Topic         string      `json:"topic"`
	Notifications interface{} `json:"notifications"`
}

// JobTrigger is the trigger that kicks off the job
type JobTrigger struct {
	ID           string `json:"ID"`
	CalendarName string `json:"calendarName"`
	MayFireAgain bool   `json:"mayFireAgain"`
	NextFireTime int    `json:"nextFireTime"`
	Description  string `json:"Description"`
	StartTime    int    `json:"startTime"`
	EndTime      int    `json:"endTime"`
	Priority     int    `json:"priority"`
	State        string `json:"state"`
}

// ListJobs lists jobs on the Platform
func ListJobs(base, auth string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s", base, "/jobs")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	curlCmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		return bodybytes, resp.StatusCode, curl, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodybytes, resp.StatusCode, curl, nil
}

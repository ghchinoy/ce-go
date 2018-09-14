package ce

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
)

var testjob = []byte(`{
	"trigger": {
		"cron": "0 0/15 * 1/1 * ? *"
	},
	"name": "Test Job",
	"description": "My test job",
	"method": "GET",
	"uri": "/elements/api-v2/instances"
  }`)

var job struct {
	Trigger struct {
		CRON string `json:"cron"`
	} `json:"trigger"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Method      string `json:"method"`
	URI         string `json:"uri"`
}

func TestListJobs(t *testing.T) {
	bodybytes, code, _, err := ListJobs(base, auth)
	if err != nil {
		t.Errorf("error %s", err)
	}
	if code != 200 {
		t.Errorf("non-200 code %v", code)
	}
	if len(bodybytes) < 1 {
		t.Errorf("body length too small %v", len(bodybytes))
	}
}

func uniqueTestJobBytes() ([]byte, error) {
	err := json.Unmarshal(testjob, &job)
	if err != nil {
		return nil, err
	}
	namedesc := fmt.Sprintf("test job %s", time.Now().Format(time.Stamp))
	job.Description = namedesc
	job.Name = namedesc
	jobbytes, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}
	return jobbytes, nil
}

func TestCreateJob(t *testing.T) {
	jobbytes, err := uniqueTestJobBytes()
	if err != nil {
		t.Errorf("setup error %s", err)
	}
	_, code, _, err := CreateJob(base, auth, jobbytes)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if code != 200 {
		t.Errorf("non-200 error code %v", code)
	}
}

func TestDeleteJob(t *testing.T) {
	jobbytes, err := uniqueTestJobBytes()
	if err != nil {
		t.Errorf("setup error %s", err)
	}

	resultbytes, code, _, err := CreateJob(base, auth, jobbytes)
	if err != nil {
		t.Errorf("couldn't set up - error: %s", err)
	}
	if code != 200 {
		t.Errorf("coudn't set up - non-200 code: %v", code)
		log.Printf("%s", resultbytes)
	}

	var resp map[string]interface{}
	err = json.Unmarshal(resultbytes, &resp)
	if err != nil {
		t.Errorf("unable to parse body: %s", err)
	}
	_, code, _, err = DeleteJob(base, auth, resp["id"].(string))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if code != 200 {
		t.Errorf("non-200 error code %v", code)
	}
}

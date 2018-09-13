package ce

import (
	"testing"
)

func TestListJobs(t *testing.T) {

	//[]byte, int, string, error)
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

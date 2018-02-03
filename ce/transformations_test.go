package ce

import "testing"

func TestGetTransformations(t *testing.T) {
	bodybytes, status, _, err := GetTransformations(base, auth)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if status != 200 {
		t.Errorf("Error: %v", status)
		t.Logf("%s", bodybytes)
	}
}

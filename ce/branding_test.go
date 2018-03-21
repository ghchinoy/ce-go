package ce

import "testing"

func TestSetBranding(t *testing.T) {

	bodybytes, status, _, err := SetBranding(base, auth, DefaultBranding, false)
	if err != nil {
		t.Errorf("Test failed %s", err.Error())
	}
	if status != 200 {
		t.Errorf("Test failed with non-200, %v", status)
		t.Logf("%s", bodybytes)
	}
}

func TestGetBranding(t *testing.T) {

	bodybytes, status, _, err := GetBranding(base, auth, false)
	if err != nil {
		t.Errorf("Test failed %s", err.Error())
	}
	if status != 200 {
		t.Errorf("Test failed with non-200, %v", status)
		t.Logf("%s", bodybytes)
	}
}

func TestResetBranding(t *testing.T) {

	bodybytes, status, _, err := ResetBranding(base, auth, false)
	if err != nil {
		t.Errorf("Test failed %s", err.Error())
	}
	if status != 200 {
		t.Errorf("Test failed with non-200, %v", status)
		t.Logf("%s", bodybytes)
	}
}

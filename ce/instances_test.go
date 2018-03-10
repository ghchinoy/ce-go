package ce

import (
	"strconv"
	"testing"
)

func TestEnableElementInstanceEvents(t *testing.T) {
	bodybytes, status, _, err := EnableElementInstanceEvents(base, auth, strconv.Itoa(452326), true, true)
	if err != nil {
		t.Errorf("Something went wrong: %s", err.Error())
	}
	if status != 200 {
		t.Logf("%s", bodybytes)
		t.Errorf("Status: %v", status)
	}
}

func TestDisableElementInstanceEvents(t *testing.T) {
	bodybytes, status, _, err := EnableElementInstanceEvents(base, auth, strconv.Itoa(452326), false, true)
	if err != nil {
		t.Errorf("Something went wrong: %s", err.Error())
	}
	if status != 200 {
		t.Logf("%s", bodybytes)
		t.Errorf("Status: %v", status)
	}
}

func TestEnableElementInstance(t *testing.T) {
	bodybytes, status, _, err := EnableElementInstance(base, auth, strconv.Itoa(452319), true, true)
	if err != nil {
		t.Errorf("Something went wrong: %s", err.Error())
	}
	if status != 200 {
		t.Logf("%s", bodybytes)
		t.Errorf("Status: %v", status)
	}
}

func TestDisableElementInstance(t *testing.T) {
	bodybytes, status, _, err := EnableElementInstance(base, auth, strconv.Itoa(452319), false, true)
	if err != nil {
		t.Errorf("Something went wrong: %s", err.Error())
	}
	if status != 200 {
		t.Logf("%s", bodybytes)
		t.Errorf("Status: %v", status)

	}
}

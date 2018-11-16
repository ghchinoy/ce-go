package ce

import (
	"testing"
)

func TestAddToElementsDenyList(t *testing.T) {
	deny := []string{"dropbox", "jira", "sendgrid", "twilio"}
	bodybytes, code, _, err := AddToElementsDenyList(base, auth, deny)
	if err != nil {
		t.Errorf("error %s", err)
	}
	if code == 403 {
		t.Logf("Insufficient privs %v", code)

	} else if code != 200 {
		t.Errorf("non-200 code %v", code)
	}
	if len(bodybytes) < 1 {
		t.Errorf("body length too small %v", len(bodybytes))
	}
}

func TestResetElementsDenyList(t *testing.T) {
	bodybytes, code, _, err := ResetElementsDenyList(base, auth)
	if err != nil {
		t.Errorf("error %s", err)
	}
	if code == 403 {
		t.Logf("Insufficient privs %v", code)

	} else if code != 200 {
		t.Errorf("non-200 code %v", code)
	}
	if len(bodybytes) < 1 {
		t.Errorf("body length too small %v", len(bodybytes))
	}

}

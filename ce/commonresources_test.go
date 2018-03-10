package ce

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var (
	base string
	auth string
)

func TestMain(m *testing.M) {

	base = os.Getenv("CE_BASE")
	auth = os.Getenv("CE_AUTH")
	os.Exit(m.Run())
}

func TestImportResource(t *testing.T) {

	commoncontact := `{"fields":[{"type":"string","path":"country"},{"type":"string","path":"firstName"},{"type":"string","path":"lastName"},{"type":"string","path":"city"},{"type":"string","path":"phone"},{"type":"string","path":"street"},{"type":"string","path":"postalCode"},{"type":"string","path":"name"},{"type":"string","path":"id"},{"type":"string","path":"state"},{"type":"string","path":"email"}],"level":"organization"}`
	err := ioutil.WriteFile("/tmp/common-contact.cro.json", []byte(commoncontact), 0644)
	if err != nil {
		t.Errorf("Error writing common contact for test: %s", err)
	}

	bodybytes, status, _, err := ImportResource(base, auth, "Test-Resource", "/tmp/common-contact.cro.json")
	if err != nil {
		t.Errorf("Test failed: %s", err)
	}
	if status != 200 {
		t.Errorf("Couldn't add contact: %s", err)
		t.Logf("%s\n", bodybytes)
	}
}

func TestCopyResource(t *testing.T) {

	now := time.Now()
	bodybytes, status, _, err := CopyResource(base, auth, "Test-Resource", fmt.Sprintf("Test-Resource-%s", now.Format("2006-01-02-15-04")))
	if err != nil {
		t.Errorf("Test failed: %s", err)
	}
	if status != 200 {
		t.Errorf("Couldn't add contact: %s", err)
		t.Logf("%s\n", bodybytes)
	}
}

package ce

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCreateFormulaInstance(t *testing.T) {

	base := "https://staging.cloud-elements.com/elements/api-v2"
	auth := "Organization fa6a81bb6138009f5a41bd4a20a5776a, User ghHckE3/EM3ntlNO0yGoMK+6bobxax6tZEdueY7P8Vg="

	// test without config
	var config FormulaInstanceConfig
	config.Name = "TestFormula"

	bodybytes, status, _, err := CreateFormulaInstance(base, auth, strconv.Itoa(19547), config)
	if err != nil {
		t.Errorf("Something went wrong: %s", err.Error())
	}
	if status != 200 {
		fmt.Printf("%s", bodybytes)
		t.Errorf("Status: %v", status)

	}

}

func TestTriggerFormulaInstanceNoTrigger(t *testing.T) {
	base := "https://staging.cloud-elements.com/elements/api-v2"
	auth := "Organization fa6a81bb6138009f5a41bd4a20a5776a, User ghHckE3/EM3ntlNO0yGoMK+6bobxax6tZEdueY7P8Vg="

	bodybytes, status, _, err := TriggerFormulaInstance(base, auth, "199701", "{}")

	if err != nil {
		t.Errorf("Something went wrong: %s", err.Error())
	}
	if status != 200 {
		fmt.Printf("%s", bodybytes)
		t.Errorf("Status: %v", status)
	}
}

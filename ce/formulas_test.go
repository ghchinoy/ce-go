package ce

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCreateFormulaInstance(t *testing.T) {
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
	bodybytes, status, _, err := TriggerFormulaInstance(base, auth, "199701", "{}")

	if err != nil {
		t.Errorf("Something went wrong: %s", err.Error())
	}
	if status != 200 {
		fmt.Printf("%s", bodybytes)
		t.Errorf("Status: %v", status)
	}
}

func TestGetFormulaInstanceExecutions(t *testing.T) {
	bodybytes, status, _, err := GetFormulaInstanceExecutions(base, auth, "199701")
	if err != nil {
		t.Errorf("Something went wrong: %s", err.Error())
	}
	if status != 200 {
		fmt.Printf("%s", bodybytes)
		t.Errorf("Status: %v", status)
	}
}

func TestCancelFormulaExecution(t *testing.T) {

}

func TestImportFormula(t *testing.T) {

	// create a dummy formula
	properties := struct {
		Body string `json:"body"`
	}{
		"done();",
	}
	steps := []Step{
		Step{
			Name:       "dummystep",
			Type:       "script",
			OnFailure:  []string{},
			OnSuccess:  []string{},
			Properties: properties,
		},
	}
	f := Formula{
		Name:  "Dummy",
		Steps: steps,
		//CreatedDate: time.Now(),
	}
	//fmt.Printf("%+v\n", f)

	bodybytes, status, _, err := ImportFormula(base, auth, f)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if status != 200 {
		fmt.Printf("%s", bodybytes)
		t.Errorf("Status: %v", status)
	}

}

/*
func TestDeleteFormula(t *testing.T) {

	bodybytes, status, _, err := DeleteFormula(base, auth, strconv.Itoa(formula.ID))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if status != 200 {
		fmt.Printf("%s", bodybytes)
		t.Errorf("Status: %v", status)
	}
}

*/

/*
func TestGetFormulaInstances(t *testing.T) {
	bodybytes, status, _, err := GetFormulaInstances(base, auth, formulaID)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if status != 200 {
		t.Errorf("Status: %v", status)
	}
}
*/

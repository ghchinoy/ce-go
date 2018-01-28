package ce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/moul/http2curl"
	"github.com/olekukonko/tablewriter"
)

const (
	// FormulasURI is the base URI for Formulas
	FormulasURI = "/formulas"
	// FormulaCancelExecutionURIFormat is the API URI for cancelling a formula instance execution
	FormulaCancelExecutionURIFormat = "/formulas/instances/executions/%s"
	// FormulaExecutionsURIFormat is the URI to obtain executions of a Formula Instance
	FormulaExecutionsURIFormat = "/formulas/instances/%s/executions"
	// FormulaRetryExecutionURI is the URI to retry a Formula execution
	FormulaRetryExecutionURI = "/formulas/instances/executions/%s/retries"
	// FormulaURIFormat is the main, partial API URI for Formula
	FormulaURIFormat = "/formulas/%s"
	// FormulaInstancesURI is the main API URI for Formula Instances
	FormulaInstancesURI = "/formulas/instances"
	// FormulaInstancesURIFormat is the URI to obtain instances of a Formula template
	FormulaInstancesURIFormat       = "/formulas/%s/instances"
	FormulaInstanceDetailsURIFormat = "/formulas/instances/%s"
	FormulaInstanceDeleteURIFormat  = "/formulas/%v/instances/%s"
)

// Formula represents the structure of a CE Formula
type Formula struct {
	ID             int               `json:"id,omitempty"`
	Name           string            `json:"name"`
	UserID         int               `json:"userId"`
	AccountID      int               `json:"accountId"`
	CreatedDate    time.Time         `json:"createdDate"`
	Steps          []Step            `json:"steps"`
	Triggers       []Trigger         `json:"triggers"`
	Active         bool              `json:"active"`
	SingleThreaded bool              `json:"singleThreaded"`
	Configuration  []Configuration   `json:"configuration"`
	API            string            `json:"api"`
	Instances      []FormulaInstance `json:"instances,omitempty"`
}

// Step represents a Formula step
type Step struct {
	ID         int         `json:"id"`
	OnSuccess  []string    `json:"onSuccess"`
	OnFailure  []string    `json:"onFailure"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
}

// Trigger represents an action that starts a Formula
type Trigger struct {
	ID         int         `json:"id"`
	Type       string      `json:"type"`
	OnSuccess  []string    `json:"onSuccess"`
	OnFailure  []string    `json:"onFailure"`
	Async      bool        `json:"async"`
	Name       string      `json:"name"`
	Properties interface{} `json:"properties"`
}

// Configuration represents a configuration for a formula
type Configuration struct {
	ID       int `json:"id"`
	Key      string
	Name     string
	Type     string
	Required bool
}

// FormulaInstance represents a configured instance of a Formula
type FormulaInstance struct {
	ID            int         `json:"id"`
	Formula       Formula     `json:"formula"`
	Name          string      `json:"name"`
	CreatedDate   time.Time   `json:"createdDate"`
	Settings      interface{} `json:"settings"`
	Active        bool        `json:"active"`
	Configuration interface{} `json:"configuration"`
}

// FormulaInstanceConfig represents a configuration used when creating an Instance of a Formula
type FormulaInstanceConfig struct {
	Name          string      `json:"name"`
	Active        bool        `json:"active"`
	Configuration interface{} `json:"configuration,omitempty"`
}

// FormulaInstanceCreationResponse is the response returned when a Formula Instance is triggered
type FormulaInstanceCreationResponse struct {
	ID        int    `json:"id"`
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
}

// FormulaInstanceExecution is a brief info about an instance Execution
type FormulaInstanceExecution struct {
	ID                int       `json:"id"`
	FormulaInstanceID int       `json:"formulaInstanceId"`
	Status            string    `json:"status"`
	CreateDate        time.Time `json:"createdDate"`
	UpdatedDate       time.Time `json:"updatedDate"`
}

// DeleteFormula deletes a Formula
func DeleteFormula(base, auth string, formulaID string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(FormulaURIFormat, formulaID),
	)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlCmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		return bodybytes, -1, "", err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodybytes, resp.StatusCode, curl, nil
}

// ImportFormula imports a Formula template, given a Formula
func ImportFormula(base, auth string, f Formula) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		"/formulas",
	)

	fbytes, err := json.Marshal(f)
	if err != nil {
		return bodybytes, -1, "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(fbytes))
	if err != nil {
		//fmt.Println("Can't construct request", err.Error())
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlCmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Cannot process response", err.Error())
		return bodybytes, -1, "", err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodybytes, resp.StatusCode, curl, nil
}

// CancelFormulaExecution cancels an execution given an Execution ID
func CancelFormulaExecution(base, auth string, executionID string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(FormulaCancelExecutionURIFormat, executionID),
	)
	// construct a fixed json body for sending cancelled status
	cancelmessage := struct {
		Status string `json:"status"`
	}{"cancelled"}
	cancelbytes, err := json.Marshal(cancelmessage)
	if err != nil {
		//fmt.Println("Can't even")
		return bodybytes, -1, "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(cancelbytes))
	if err != nil {
		//fmt.Println("Can't construct request", err.Error())
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlCmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Cannot process response", err.Error())
		return bodybytes, -1, "", err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodybytes, resp.StatusCode, curl, nil
}

// GetFormulaInstanceExecutions returns a list of Formula Instance Executions given a Formula Instance ID
func GetFormulaInstanceExecutions(base, auth string, formulaInstanceID string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(FormulaExecutionsURIFormat, formulaInstanceID),
	)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//fmt.Println("Can't construct request", err.Error())
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	curlCmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Cannot process response", err.Error())
		return bodybytes, -1, "", err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodybytes, resp.StatusCode, curl, nil
}

// TriggerFormulaInstance invokes a Formula Instance with the given trigger
func TriggerFormulaInstance(base, auth string, formulaTemplateID, triggerBody string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(FormulaExecutionsURIFormat, formulaTemplateID),
	)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(triggerBody)))
	if err != nil {
		//fmt.Println("Can't construct request", err.Error())
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlCmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Cannot process response", err.Error())
		return bodybytes, -1, "", err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodybytes, resp.StatusCode, curl, nil
}

// CreateFormulaInstance creates an instance of a Formula given a FormulaInstanceConfig
func CreateFormulaInstance(base, auth string, formulaTemplateID string, config FormulaInstanceConfig) ([]byte, int, string, error) {
	var bodybytes []byte

	url := fmt.Sprintf("%s%s", base, fmt.Sprintf(FormulaInstancesURIFormat, formulaTemplateID))

	fibytes, err := json.Marshal(config)
	//fmt.Println(url)
	//fmt.Printf("%s\n", fibytes)
	if err != nil {
		//fmt.Println("Unable to convert to Formula Instance configuration json", err.Error())
		return bodybytes, -1, "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(fibytes))
	if err != nil {
		//fmt.Println("Unable to create request", err.Error())
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlCmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Cannot process response", err.Error())
		return bodybytes, -1, "", err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodybytes, resp.StatusCode, curl, nil

}

// ExportAllFormulasToDir creates a directory given and exports all Formula JSON files
func ExportAllFormulasToDir(base, auth string, dirname string) error {
	formulaListByes, _, _, err := FormulasList(base, auth)
	if err != nil {
		return err
	}
	var formulas []Formula
	err = json.Unmarshal(formulaListByes, &formulas)
	if err != nil {
		return err
	}

	// create formulas dir
	err = os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		return err
	}
	for _, f := range formulas {
		name := fmt.Sprintf("%s.formula.json", strings.Replace(f.Name, " ", "", -1))
		formulaBytes, err := json.Marshal(f)
		if err != nil {
			break
		}
		fmt.Printf("Exporting '%s' to %s/%s\n", f.Name, dirname, name)
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dirname, name), formulaBytes, 0644)
	}

	return nil
}

// DeleteFormulaInstance deletes an Instance of a Formula
func DeleteFormulaInstance(base, auth string, instanceID string) ([]byte, int, string, error) {
	var bodybytes []byte

	// Get the Instance info
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(FormulaInstanceDetailsURIFormat, instanceID),
	)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//fmt.Println("Can't construct request", err.Error())
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlCmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		// unable to reach CE API
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var fi FormulaInstance
	err = json.Unmarshal(bodybytes, &fi)
	if err != nil {
		// unable to create Formula Instance from response
		return bodybytes, -1, curl, err
	}

	formulaID := fi.Formula.ID

	// Delete the Instance
	url = fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(FormulaInstanceDeleteURIFormat, formulaID, instanceID),
	)
	client = &http.Client{}
	req, err = http.NewRequest("DELETE", url, nil)
	if err != nil {
		//fmt.Println("Can't construct request", err.Error())
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlCmd, _ = http2curl.GetCurlCommand(req)
	curl = fmt.Sprintf("%s", curlCmd)
	resp, err = client.Do(req)
	if err != nil {
		// unable to reach CE API
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return bodybytes, resp.StatusCode, curl, nil
}

// GetInstancesOfFormula returns an Instance array, given a Formula ID and an Auth header
func GetInstancesOfFormula(id int, baseurl string, auth string) ([]FormulaInstance, error) {
	var instances []FormulaInstance

	url := fmt.Sprintf("%s%s", baseurl,
		fmt.Sprintf("/formulas/%v/instances", id))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//fmt.Println("Can't construct request", err.Error())
		return instances, err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Cannot process response", err.Error())
		return instances, err
	}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	err = json.Unmarshal(bodybytes, &instances)
	if err != nil {
		return instances, err
	}

	return instances, nil
}

// FormulaDetailsTableOutput prints to stdout an ASCII rendered table of the details of a Formula
func FormulaDetailsTableOutput(f Formula) error {

	// basic formula info
	data := [][]string{}

	if len(f.Triggers) < 1 {
		fmt.Printf("Formula %v is malformed, no trigger present\n", f.ID)

	} else {
		data = append(data, []string{
			strconv.Itoa(f.ID),
			f.Name,
			strconv.FormatBool(f.Active),
			strconv.Itoa(len(f.Steps)),
			f.Triggers[0].Type,
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "active", "steps", "trigger"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

		fmt.Println()

		// Triggers

		data = [][]string{}

		for _, v := range f.Triggers {
			data = append(data, []string{
				strconv.Itoa(v.ID),
				v.Name,
				v.Type,
				strconv.FormatBool(v.Async),
				fmt.Sprintf("%s", v.OnSuccess),
			})
		}

		table = tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Type", "Async", "Success"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

		// Steps

		fmt.Println("\nSteps")

		data = [][]string{}

		for _, v := range f.Steps {
			data = append(data, []string{
				strconv.Itoa(v.ID),
				v.Name,
				v.Type,
				fmt.Sprintf("%s", v.OnSuccess),
				fmt.Sprintf("%s", v.OnFailure),
			})
		}

		table = tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Type", "Success", "Failure"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

		// Configuration parameters
		fmt.Println("\nConfiguration")

		if len(f.Configuration) > 0 {
			data = [][]string{}
			for _, v := range f.Configuration {
				data = append(data, []string{
					strconv.Itoa(v.ID),
					v.Name,
					v.Key,
					v.Type,
					strconv.FormatBool(v.Required),
				})
			}
			table = tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Key", "Value", "Required"})
			table.SetBorder(false)
			table.AppendBulk(data)
			table.Render()
		} else {
			fmt.Println("No configuration parameters needed.")
		}

		if f.API != "" {
			fmt.Printf("\n%s -H 'Elements-Formula-Instance-Id: '\n", f.API)
		}
	}

	return nil
}

// FormulaDetailsAsBytes returns Formula template details as bytes
func FormulaDetailsAsBytes(formulaID, base, auth string) ([]byte, int, string, error) {

	var bodybytes []byte

	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(FormulaURIFormat, formulaID),
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Can't construct request", err.Error())
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
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

// FormulaUpdate performs a PATCH with a Formula
func FormulaUpdate(formulaID, base, auth string, formula Formula) ([]byte, int, error) {
	var bodybytes []byte

	formulaRequestBytes, err := json.Marshal(formula)
	if err != nil {
		return nil, -1, err
	}
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(FormulaURIFormat, formulaID),
	)
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(formulaRequestBytes))
	if err != nil {
		//fmt.Println("Can't construct request", err.Error())
		return bodybytes, -1, err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Cannot process response", err.Error())
		return bodybytes, -1, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return bodybytes, resp.StatusCode, nil
}

// FormulasList retruns a list of formulas
func FormulasList(base, auth string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s", base, FormulasURI)
	t := &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client := &http.Client{Transport: t}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// cant construct request
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accpet", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlCmd, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return bodybytes, -1, "", err
	}
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		return bodybytes, -1, curl, err
	}
	defer resp.Body.Close()
	bodybytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return bodybytes, resp.StatusCode, curl, err
	}

	return bodybytes, resp.StatusCode, curl, nil
}

// GetFormulaInstanceExecutionID returns the output of the instances/execution/{id} call
func GetFormulaInstanceExecutionID(executionID, base, auth string) ([]byte, int, string, error) {

	var bodybytes []byte
	url := fmt.Sprintf("%s%s", base,
		fmt.Sprintf(FormulaCancelExecutionURIFormat, executionID),
	)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// cant construct request
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlCmd, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return bodybytes, -1, "", err
	}
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		return bodybytes, -1, curl, err
	}
	defer resp.Body.Close()
	bodybytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return bodybytes, resp.StatusCode, curl, err
	}

	return bodybytes, resp.StatusCode, curl, nil
}

// CombinedFormulaAndInstances returns a list of Formulas with Instances
func CombinedFormulaAndInstances(formulabytes []byte, base, auth string) ([]Formula, error) {
	var formulas []Formula
	err := json.Unmarshal(formulabytes, &formulas)
	if err != nil {
		return formulas, err
	}
	for i, v := range formulas {
		if len(v.Triggers) < 1 {
			log.Printf("Formula %v is malformed, no trigger present\n", v.ID)
			break
		}
		instances, err := GetInstancesOfFormula(v.ID, base, auth)
		if err != nil {
			break
		}
		// note use of index here, since range makes a copy of slice
		// https://golang.org/ref/spec#RangeClause
		formulas[i].Instances = instances
	}

	return formulas, nil
}

// OutputFormulasList writes a nice table of formulas to stdout
func OutputFormulasList(formulabytes []byte, base, auth string) error {
	data := [][]string{}

	var formulas []Formula
	err := json.Unmarshal(formulabytes, &formulas)
	if err != nil {
		return err
	}
	for _, v := range formulas {
		if len(v.Triggers) < 1 {
			fmt.Printf("Formula %v is malformed, no trigger present\n", v.ID)
			break
		}

		var instancecount string
		instances, err := GetInstancesOfFormula(v.ID, base, auth)
		if err != nil {
			// unable to retrieve instances of formula!
			instancecount = "N/A"
		}
		instancecount = strconv.Itoa(len(instances))

		for _, t := range v.Triggers {

			api := "N/A"
			if v.Triggers[0].Type == "manual" {
				api = v.API
			}

			data = append(data, []string{
				strconv.Itoa(v.ID),
				v.Name,
				strconv.FormatBool(v.Active),
				strconv.Itoa(len(v.Steps)),
				instancecount,
				t.Type,
				strconv.Itoa(t.ID),
				fmt.Sprintf("%s", t.OnSuccess),
				api,
			},
			)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "active", "steps", "instances", "trigger", "id", "success", "api"})
	table.SetBorder(false)
	table.SetAutoMergeCells(true)
	table.AppendBulk(data)
	table.Render()

	return nil
}

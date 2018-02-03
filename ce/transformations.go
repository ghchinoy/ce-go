package ce

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/moul/http2curl"
)

// Transformation structure represents a Transformation
type Transformation struct {
	Level      string `json:"level"`
	ObjectName string `json:"objectName"`
	VendorName string `json:"vendorName"`
	StartDate  string `json:"startDate"`
	Fields     []struct {
		Type       string `json:"type"`
		Path       string `json:"path"`
		VendorPath string `json:"vendorPath"`
		Level      string `json:"level"`
	} `json:"fields"`
	Configuration []struct {
		Type       string `json:"type"`
		Properties struct {
			FromVendor bool `json:"fromVendor"`
			ToVendor   bool `json:"toVendor"`
		} `json:"properties,omitempty"`
	} `json:"configuration"`
	IsLegacy bool                 `json:"isLegacy"`
	Script   TransformationScript `json:"script,omitempty"`
}

// TransformationScript represents a script attached to a Transformation
type TransformationScript struct {
	Body                string `json:"body,omitempty"`
	MimeType            string `json:"mimeType"`
	FilterEmptyResponse bool   `json:"filterEmptyResponse"`
}

// AccountElement is the association of an account and Element with a Transformation
type AccountElement struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Level   string `json:"level"`
	Account struct {
		Status         string `json:"status"`
		Environment    string `json:"environment"`
		Active         bool   `json:"active"`
		ID             int    `json:"id"`
		DefaultAccount bool   `json:"defaultAccount"`
	} `json:"account"`
	Element Element `json:"element"`
}

// GetTransformationAssocation returns Elements associated with the given Transformation
// the expected result is an array of AccountElement
func GetTransformationAssocation(base, auth string, txname string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf("/organizations/objects/%s/transformations", txname))
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlcmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlcmd)
	res, err := client.Do(req)
	if err != nil {
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	return bodybytes, res.StatusCode, curl, nil
}

// GetTransformationsPerElement returns the transformations associated with a particular Element
func GetTransformationsPerElement(base, auth string, elementkey string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf("/organizations/elements/%s/transformations", elementkey))
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlcmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlcmd)
	res, err := client.Do(req)
	if err != nil {
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	return bodybytes, res.StatusCode, curl, nil
}

// GetTransformations lists the Transformations on the Platform
func GetTransformations(base, auth string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s/organizations/objects/definitions", base)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	curlcmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlcmd)
	res, err := client.Do(req)
	if err != nil {
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	return bodybytes, res.StatusCode, curl, nil
}

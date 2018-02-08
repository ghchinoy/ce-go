package ce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/moul/http2curl"
	"github.com/olekukonko/tablewriter"
)

const (
	// CommonResourcesURI is the base URI of the hidden API
	// for Common Object Resources; this one provides an array of common objects
	// with the element instance IDs associated, as well as details about
	// the field's heirarchy (org, account, instance)
	CommonResourcesURI = "/common-resources"
	// CommonResourcesDefinitionURIFormat provides details for a specific resource
	CommonResourcesDefinitionURIFormat = "/common-resources/%s"
	// CommonResourceURI is the base URI for common object resources
	// this is a simple object with keys being the common object names and no
	// details about associated elements or field level hierarchy
	CommonResourceURI = "/organizations/objects/definitions"
	// CommonResourceDefinitionsFormatURI is a string format for the URI of Common Object Resource definition, given a name of a Common Object
	CommonResourceDefinitionsFormatURI = "/organizations/objects/%s/definitions"
	// CommonResourceTransformationsFormatURI is the string format for the URI of an Element's transformation / mapping, given an element key and an object name
	CommonResourceTransformationsFormatURI = "/organizations/elements/%s/transformations/%s"
)

// CommonResource represents a normalized data object (resource)
type CommonResource struct {
	Name               string  `json:"name,omitempty"`
	ElementInstanceIDs []int   `json:"elementInstanceIds,omitempty"`
	Fields             []Field `json:"fields"`
	Level              string  `json:"level,omitempty"`
}

// Field is a set of  a common resource fields
type Field struct {
	Type            string `json:"type"`
	Path            string `json:"path"`
	AssociatedLevel string `json:"organization,omitempty"`
	AssociatedID    int    `json:"associatedId,omitempty"`
}

// ImportResource imports a common resource object to the Platform
func ImportResource(base, auth string, name, filepath string) ([]byte, int, string, error) {
	var bodybytes []byte

	// read in file
	filebytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return bodybytes, -1, "", err
	}

	bodybytes, status, curlcmd, err := createResource(base, auth, name, filebytes)
	if err != nil {
		return bodybytes, -1, "", err
	}

	return bodybytes, status, curlcmd, nil
}

func createResource(base, auth string, name string, resourcebytes []byte) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(CommonResourceDefinitionsFormatURI, name),
	)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(resourcebytes))
	if err != nil {
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
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

// CopyResource copies a Resource to another
func CopyResource(base, auth string, source, target string) ([]byte, int, string, error) {
	var bodybytes []byte
	originalbytes, status, curlcmd1, err := GetResourceDefinition(base, auth, source, false)
	if err != nil {
		return bodybytes, -1, "", err
	}
	if status != 200 {
		return bodybytes, status, "", err
	}

	bodybytes, status, curlcmd2, err := createResource(base, auth, target, originalbytes)
	if err != nil {
		return bodybytes, -1, "", err
	}

	return bodybytes, status, fmt.Sprintf("%s\n%s", curlcmd1, curlcmd2), nil
}

// DeleteResource deletes a common resource object
func DeleteResource(base, auth, resourceName string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(CommonResourceDefinitionsFormatURI, resourceName),
	)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Can't construct request", err.Error())
		os.Exit(1)
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
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

// GetResourceDefinition returns a Resource's definition
func GetResourceDefinition(base, auth string, resourceName string, details bool) ([]byte, int, string, error) {
	var bodybytes []byte
	var url string
	if details {
		url = fmt.Sprintf("%s%s",
			base,
			fmt.Sprintf(CommonResourcesDefinitionURIFormat, resourceName),
		)
	} else {
		url = fmt.Sprintf("%s%s",
			base,
			fmt.Sprintf(CommonResourceDefinitionsFormatURI, resourceName),
		)
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// cant construct request
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accpet", "application/json")
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

// ResourcesList retruns a list of common resource objects
func ResourcesList(base, auth string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s", base, CommonResourcesURI)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// cant construct request
		return bodybytes, -1, "", err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accpet", "application/json")
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

// OutputResourcesList prints a nicely formatted table to stdout
func OutputResourcesList(resourcesbytes []byte) error {
	data := [][]string{}

	var commonResources []CommonResource
	err := json.Unmarshal(resourcesbytes, &commonResources)
	if err != nil {
		fmt.Printf("Response not a list of Common Resources, %s", err.Error())
		return err
	}

	for _, v := range commonResources {

		var fieldList string
		if len(v.Fields) > 0 {
			var fields []string
			for _, f := range v.Fields {
				fields = append(fields, f.Path)
			}
			fieldList = strings.Join(fields[:], ", ")
			fieldList = " [" + fieldList + "]"
		}

		var instanceList string
		if len(v.ElementInstanceIDs) > 0 {
			var ids []string
			for _, i := range v.ElementInstanceIDs {
				ids = append(ids, strconv.Itoa(i))
			}
			instanceList = strings.Join(ids[:], ", ")
			instanceList = " [" + instanceList + "]"
		}

		data = append(data, []string{
			v.Name,
			strconv.Itoa(len(v.ElementInstanceIDs)) + instanceList,
			strconv.Itoa(len(v.Fields)),
			fieldList,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Mapped Instances", "#", "Fields"})
	table.SetBorder(false)
	table.SetColWidth(40)
	table.AppendBulk(data)
	table.Render()

	return nil
}

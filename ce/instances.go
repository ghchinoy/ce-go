package ce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/moul/http2curl"
	"github.com/olekukonko/tablewriter"
)

const (
	InstancesURI                         = "/instances"
	InstancesFormatURI                   = "/instances/%s"
	InstanceConfigurationURI             = "/instances/configuration"
	InstanceConfigurationFormatURI       = "/instances/configuration/%s"
	InstanceDocsURI                      = "/instances/docs"
	InstanceOperationDocsFormatURI       = "/instances/docs/%s"
	InstancesEventsURI                   = "/instances/events"
	InstancesEventsAnalyticsAccountsURI  = "/instances/events/analytics/accounts"
	InstancesEventsAnalyticsInstancesURI = "/instances/events/analytics/instances"
	InstancesEventsFormatURI             = "/instances/events/%s"
	InstancesTransformationsURI          = "/instances/transformations"
	InstanceTransformationsFormatURI     = "/instances/%s/transformations"
	InstanceDocFormatURI                 = "/instances/%s/docs"
	InstancesEnableURI                   = "/instances/enabled"

	InstanceDefinitions_ID       = "/instances/%s/objects/definitions"
	InstanceDefinitions_Token    = "/instances/objects/definitions"
	InstanceOperationOAI_ID      = "/instances/%s/docs/%s"
	InstanceOAIByOperation_ID    = "/instances/%s/docs/%s/definitions"
	InstanceOAIByOperation_Token = "/instances/docs/%s/definitions"
)

// Instance represents an Element Instance
type Instance struct {
	ID                     int
	Name                   string
	CreatedDate            string
	Token                  string
	Element                Element
	ElementID              int
	Tags                   []string
	ProvisionInteractions  interface{}
	Valid                  bool
	Disabled               bool
	MaxCacheSize           int
	CacheTimeToLive        int
	Configuration          InstanceConfiguration
	EventsEnabled          bool
	TraceLoggingEnabled    bool
	CachingEnabled         bool
	ExternalAuthentication string
	User                   User
	TransformationData     []struct {
		ObjectName string `json:"objectName"`
		VendorName string `json:"vendorName"`
	} `json:"transformationData"`
}

// Execute is a HTTP command that returns bytes, HTTP status, and a curl command
func Execute(method, url, auth string) ([]byte, int, string, error) {

	var bodybytes []byte

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Can't construct request", err.Error())
		os.Exit(1)
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

	// verify it's a collection of Element Instances
	var instances []ElementInstance
	err = json.Unmarshal(bodybytes, &instances)
	if err != nil {
		//fmt.Println("Unable to read Element instances")
		bodybytes, _ = json.Marshal(instances)
	}

	return bodybytes, resp.StatusCode, curl, nil
}

// EnableElementInstance enables or disables an instance given an instance ID and an enable status
func EnableElementInstance(base, auth string, instanceID string, enable bool, debug bool) ([]byte, int, string, error) {

	// get the instance info
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstancesFormatURI, instanceID),
	)
	if debug {
		log.Println("Getting instance info...")
		log.Println(url)
	}
	bodybytes, status, curlcmd, err := Execute("GET", url, auth)
	if err != nil {
		if debug {
			log.Printf("%s", bodybytes)
		}
		return bodybytes, status, curlcmd, err
	}

	var instance ElementInstance
	err = json.Unmarshal(bodybytes, &instance)
	if err != nil {
		if debug {
			log.Printf("%s", bodybytes)
		}
		return bodybytes, status, curlcmd, err
	}
	if debug {
		log.Printf("Instance %v %s/%s", instance.ID, instance.Element.Key, instance.Name)
	}

	// enable | disable an Element Instance
	method := "PUT"
	if !enable {
		method = "DELETE"
	}
	auth = fmt.Sprintf("%s, Element %s", auth, instance.Token)
	url = fmt.Sprintf("%s%s", base, fmt.Sprintf(InstancesFormatURI, instanceID))
	if debug {
		log.Printf("%s %s", method, url)
	}
	bodybytes, status, curlcmd, err = Execute(method, url, auth)
	if err != nil {
		if debug {
			log.Printf("%s", bodybytes)
		}
		return bodybytes, status, curlcmd, err
	}

	return bodybytes, status, curlcmd, nil

}

// GetAllInstances returns the Element Instances for the authed user
func GetAllInstances(base, auth string) ([]byte, int, string, error) {
	var bodybytes []byte

	url := fmt.Sprintf("%s%s", base, InstancesURI)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Can't construct request", err.Error())
		os.Exit(1)
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

	// verify it's a collection of Element Instances
	var instances []ElementInstance
	err = json.Unmarshal(bodybytes, &instances)
	if err != nil {
		//fmt.Println("Unable to read Element instances")
		bodybytes, _ = json.Marshal(instances)
	}

	return bodybytes, resp.StatusCode, curl, nil
}

// DeleteElementInstance deletes an instance given its ID
func DeleteElementInstance(base, auth string, instanceID string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstancesFormatURI, instanceID),
	)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Can't construct request", err.Error())
		os.Exit(1)
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
	return bodybytes, resp.StatusCode, curl, nil
}

// GetInstanceInfo obtains details of an Instance
func GetInstanceInfo(base, auth, instanceID string) ([]byte, int, string, error) {
	var bodybytes []byte

	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstancesFormatURI, instanceID),
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Can't construct request", err.Error())
		os.Exit(1)
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

	return bodybytes, resp.StatusCode, curl, nil
}

// GetInstanceOAI returns the OAI Spec for an Instance ID
func GetInstanceOAI(base, auth, instanceID string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstanceDocFormatURI, instanceID),
	)
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
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodybytes, resp.StatusCode, curl, nil
}

// GetInstanceTransformations retrieves transformations given an Element Instance, uses the Element Token in a header
func GetInstanceTransformations(base, auth string, id string) ([]byte, int, string, error) {
	var bodybytes []byte

	// Get the Element Instance token
	bodybytes, _, _, err := GetInstanceInfo(base, auth, id)
	if err != nil {
		fmt.Println("Unable to retrieve instance", err)
		os.Exit(1)
	}
	var instance Instance
	err = json.Unmarshal(bodybytes, &instance)
	token := instance.Token
	auth = fmt.Sprintf("%s, Element %s", auth, token)

	url := fmt.Sprintf("%s%s", base, InstancesTransformationsURI)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Can't construct request", err.Error())
		os.Exit(1)
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

	return bodybytes, resp.StatusCode, curl, nil
}

// GetInstanceObjectDefinitions returns the schema definitions for an Instance
func GetInstanceObjectDefinitions(base, auth, instanceID string) ([]byte, int, string, error) {
	var bodybytes []byte

	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstanceDefinitions_ID, instanceID),
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Can't construct request", err.Error())
		os.Exit(1)
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

	return bodybytes, resp.StatusCode, curl, nil
}

// GetInstanceOperationDefinition returns the bytes of a call to get Instance schema definitions
func GetInstanceOperationDefinition(base, auth, instanceID, operationName string) ([]byte, int, string, error) {
	var bodybytes []byte

	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstanceOAIByOperation_ID, instanceID, operationName),
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Can't construct request", err.Error())
		os.Exit(1)
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

	return bodybytes, resp.StatusCode, curl, nil
}

// OutputInstanceDetails outputs Instance details
func OutputInstanceDetails(bodybytes []byte) error {
	var i Instance
	err := json.Unmarshal(bodybytes, &i)
	if err != nil {
		return err
	}
	data := [][]string{}

	data = append(data, []string{
		strconv.Itoa(i.ID),
		i.Element.Key,
		i.Name,
		strconv.FormatBool(i.Valid),
		strconv.FormatBool(i.Disabled),
		strconv.FormatBool(i.EventsEnabled),
		fmt.Sprintf("%s", i.Tags),
		i.Token,
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Key", "Name", "Valid", "Disabled", "Events", "Tags", "Token"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}

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
	Disabled               bool `json:"disabled"`
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

// EnableElementInstanceEvents will enable or disable events on an Element Instance without requiring reauthentication
func EnableElementInstanceEvents(base, auth string, instanceID string, enable bool, debug bool) ([]byte, int, string, error) {
	// get the Instance, since the element key is needed for the PUT
	// get the instance info
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstancesFormatURI, instanceID),
	)
	if debug {
		log.Println("Getting instance info...")
		log.Println("GET", url)
	}
	bodybytes, status, curlcmd, err := Execute("GET", url, auth)
	if debug {
		log.Printf("Status %v", status)
	}
	if err != nil {
		if debug {
			log.Printf("%s", bodybytes)
		}
		return bodybytes, status, curlcmd, err
	}
	if status != 200 {
		return bodybytes, status, curlcmd, fmt.Errorf("Status code %v", status)
	}
	var instance ElementInstance
	if debug {
		log.Printf("bodybytes len %v", len(bodybytes))
	}
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

	// change "configuration" "event.notification.enabled" to enable
	if enable {
		instance.Configuration.EventNotificationEnabled = "true"
	} else {
		instance.Configuration.EventNotificationEnabled = "false"
	}

	// PUT to /elements/ELEMENT.KEY/instances/INSTANCEID?reAuthenticate=false the full body with configuration change
	url = fmt.Sprintf("%s%s", base,
		fmt.Sprintf(ElementInstancesFormatURINoReauthURI, strconv.Itoa(instance.Element.ID), strconv.Itoa(instance.ID)),
	)
	requestbytes, err := json.Marshal(instance)
	if err != nil {
		return bodybytes, -1, curlcmd, err
	}
	bodybytes, status, curlcmd, err = ExecuteWithBody("PUT", url, auth, requestbytes)
	if err != nil {
		return bodybytes, status, curlcmd, err
	}
	if status != 200 {
		return bodybytes, status, curlcmd, fmt.Errorf("Non-200 status %v", status)
	}
	return bodybytes, status, curlcmd, nil
}

// EnableElementInstanceTraceLogging enables or disables an Element Instance's
// trace logging
func EnableElementInstanceTraceLogging(base, auth string, instanceID string, enable, debug bool) ([]byte, int, string, error) {

	// Get the Element Instance
	bodybytes, status, curlcmd, err := GetInstanceInfo(base, auth, instanceID)
	if err != nil {
		return bodybytes, status, curlcmd, err
	}
	var i Instance
	err = json.Unmarshal(bodybytes, &i)
	if err != nil {
		return bodybytes, status, curlcmd, err
	}

	// enable/disable trace logging
	i.TraceLoggingEnabled = enable
	requestbytes, err := json.Marshal(i)
	if err != nil {
		return bodybytes, status, curlcmd, err
	}

	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(InstancesFormatURI, instanceID),
	)
	if debug {
		log.Printf("Setting Element Instance %s trace logging to %v ...", instanceID, enable)
		log.Println("GET", url)
	}
	bodybytes, status, curlcmd, err = ExecuteWithBody("POST", url, auth, requestbytes)
	if err != nil {
		return bodybytes, status, curlcmd, err
	}
	if status != 200 {
		return bodybytes, status, curlcmd, fmt.Errorf("Non-200 status %v", status)
	}
	return bodybytes, status, curlcmd, nil
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
		log.Println("GET", url)
	}
	bodybytes, status, curlcmd, err := Execute("GET", url, auth)
	if debug {
		log.Printf("Status %v", status)
	}
	if err != nil {
		if debug {
			log.Printf("%s", bodybytes)
		}
		return bodybytes, status, curlcmd, err
	}
	if status != 200 {
		return bodybytes, status, curlcmd, fmt.Errorf("Status code %v", status)
	}

	var instance ElementInstance
	if debug {
		log.Printf("bodybytes len %v", len(bodybytes))
	}
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
	url = fmt.Sprintf("%s%s", base, InstancesEnableURI)
	if debug {
		log.Printf("%s %s", method, url)
	}
	enablebytes, status, curlcmd, err := Execute(method, url, auth)
	if err != nil {
		if debug {
			log.Printf("%s", enablebytes)
		}
		return enablebytes, status, curlcmd, err
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
		strconv.FormatBool(i.TraceLoggingEnabled),
		fmt.Sprintf("%s", i.Tags),
		i.Token,
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Key", "Name", "Valid", "Disabled", "Events", "Trace", "Tags", "Token"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}

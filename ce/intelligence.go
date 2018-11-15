package ce

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/moul/http2curl"
)

// metadata is only available in production
// and only via a certain role
// TODO proper error messages for the previous
// TODO good output - nice table
// TODO csv output
// TODO filtered output

// per Element metadata
// https://console.cloud-elements.com/elements/api-v2/elements/201/metadata
// is in ce.elements ElementsMetadataFormatURI

// TODO decide whether this should be rolled into ce.element

// IntelligenceURI provides the endpoint for metadata
var IntelligenceURI = "/elements/metadata"

// Metadata is the individual Intelligence result
type Metadata struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Key    string `json:"key,omitempty"`
	Active bool   `json:"active,omitempty"`
	API    struct {
		Type        string `json:"type,omitempty"`
		Version     string `json:"version,omitempty"`
		ContentType string `json:"contentType,omitempty"`
	} `json:"api,omitempty"`
	AuthenticationTypes []string `json:"authenticationTypes,omitempty"`
	Beta                bool     `json:"beta,omitempty"`
	Bulk                struct {
		Upload   bool `json:"upload,omitempty"`
		Download bool `json:"download,omitempty"`
	} `json:"bulk,omitempty"`
	Description string `json:"description,omitempty"`
	Discovery   struct {
		CustomFields          bool `json:"customFields,omitempty"`
		CustomObjects         bool `json:"customObjects,omitempty"`
		EndpointCustomFields  bool `json:"endpointCustomFields,omitempty"`
		EndpointCustomObjects bool `json:"endpointCustomObjects,omitempty"`
		CRUDSSupported        bool `json:"crudsSupported,omitempty"`
	} `json:"discovery,omitempty"`
	DisplayOrder           int              `json:"displayOrder,omitempty"`
	VendorDocumentationURL string           `json:"vendorDocumentationUrl,omitempty"`
	ElementType            string           `json:"elementType,omitempty"`
	Events                 []string         `json:"events,omitempty"`
	Image                  string           `json:"image,omitempty"`
	Notes                  string           `json:"notes,omitempty"`
	Objects                []MetadataObject `json:"objects,omitempty"`
	Usage                  struct {
		InstanceCount int `json:"instanceCount,omitempty"`
		CustomerCount int `json:"customerCount,omitempty"`
		Traffic       int `json:"traffic,omitempty"`
	} `json:"usage,omitempty"`
	Captured    bool   `json:"captured,omitempty"`
	Extended    bool   `json:"extended,omitempty"`
	Hub         string `json:"hub,omitempty"`
	PricingTier string `json:"pricingTier,omitempty"`

	Transformations    bool   `json:"transformations,omitempty"`
	Churros            bool   `json:"churros,omitempty"`
	ElementClass       string `json:"elementClass,omitempty"`
	NormalizedPaging   bool   `json:"normalizedPaging,omitempty"`
	SwaggerValidated   bool   `json:"swaggerValidated,omitempty"`
	Cloneable          bool   `json:"cloneable,omitempty"`
	AuthenticationType string `json:"authenticationType,omitempty"`
	Extendable         bool   `json:"extendable,omitempty"`
	/*
		Events             struct {
			Supported      bool     `json:"supported"`
			Methods        []string `json:"methods"`
			PollingVersion string   `json:"pollingVersion"`
			Polling        struct {
				Subscriptions PollType `json:"subscriptions"`
				Invoices      PollType `json:"invoices"`
				Plans         PollType `json:"plans"`
				Customers     PollType `json:"customers"`
				Transactions  PollType `json:"transactions"`
			} `json:"polling"`
		} `json:"events"`
	*/

}

// MetadataObject is an object structure
type MetadataObject struct {
	ElementID           int      `json:"elementId,omitempty"`
	ID                  int      `json:"id,omitempty"`
	Name                string   `json:"name,omitempty"`
	CustomFields        bool     `json:"customFields,omitempty"`
	OwnerAccountID      int      `json:"ownerAccountID,omitempty"`
	EventsEnabled       bool     `json:"eventsEnabled,omitempty"`
	Description         string   `json:"description,omitempty"`
	OperationsSupported []string `json:"operationsSupported,omitempty"` // this appears to be HTTP method
	IsNestedObject      bool     `json:"isNestedObject,omitempty"`
	MetadataDiscovery   bool     `json:"metadataDiscovery,omitempty"`
	NativeBulkUpload    bool     `json:"nativeBulkUpload,omitempty"`
	NativeBulkDownload  bool     `json:"nativeBulkDownload,omitempty"`
	VendorEventTypes    string   `json:"vendorEventTypes,omitempty"`
	Notes               string   `json:"notes,omitempty"`
}

// PollType is a type of polling reference
type PollType struct {
	PollURL    string   `json:"pollUrl"`
	EventTypes []string `json:"eventTypes"`
}

// Intelligence is a metadata container for use for sorting, etc
type Intelligence []Metadata

// sort by ID
func (is Intelligence) Len() int           { return len(is) }
func (is Intelligence) Less(i, j int) bool { return is[i].ID < is[j].ID }
func (is Intelligence) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }

// ByCustomerCount implements sort.Interface for Intelligence
type ByCustomerCount []Metadata

func (e ByCustomerCount) Len() int { return len(e) }
func (e ByCustomerCount) Less(i, j int) bool {
	return e[i].Usage.CustomerCount > e[j].Usage.CustomerCount
}
func (e ByCustomerCount) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// ByInstanceCount implements sort.Interface for Intelligence
type ByInstanceCount []Metadata

func (e ByInstanceCount) Len() int { return len(e) }
func (e ByInstanceCount) Less(i, j int) bool {
	return e[i].Usage.InstanceCount > e[j].Usage.InstanceCount
}
func (e ByInstanceCount) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// ByTraffic implements sort.Interface for Intelligence
type ByTraffic []Metadata

func (e ByTraffic) Len() int { return len(e) }
func (e ByTraffic) Less(i, j int) bool {
	return e[i].Usage.Traffic > e[j].Usage.Traffic
}
func (e ByTraffic) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// ByIntHub implements sort.Interface for Intelligence
type ByIntHub []Metadata

func (e ByIntHub) Len() int           { return len(e) }
func (e ByIntHub) Less(i, j int) bool { return e[i].Hub < e[j].Hub }
func (e ByIntHub) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// ByIntName implements sort.Interface for Intelligence
type ByIntName []Metadata

func (e ByIntName) Len() int           { return len(e) }
func (e ByIntName) Less(i, j int) bool { return strings.ToLower(e[i].Name) < strings.ToLower(e[j].Name) }
func (e ByIntName) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// ByAPIType implements sort.Interface for Intelligence
type ByAPIType []Metadata

func (e ByAPIType) Len() int { return len(e) }
func (e ByAPIType) Less(i, j int) bool {
	return strings.ToLower(e[i].API.Type) < strings.ToLower(e[j].API.Type)
}
func (e ByAPIType) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// ByAuthn implements sort.Interface for Intelligence
type ByAuthn []Metadata

func (e ByAuthn) Len() int { return len(e) }
func (e ByAuthn) Less(i, j int) bool {
	return strings.ToLower(e[i].AuthenticationType) < strings.ToLower(e[j].AuthenticationType)
}
func (e ByAuthn) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// GetIntelligence returns all Elements as bytes
func GetIntelligence(base, auth string) ([]byte, int, string, error) {

	var bodybytes []byte

	url := fmt.Sprintf("%s%s", base, fmt.Sprintf("%s?expand=true", IntelligenceURI))

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

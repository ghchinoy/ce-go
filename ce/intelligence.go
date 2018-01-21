package ce

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/moul/http2curl"
	"github.com/olekukonko/tablewriter"
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

var IntelligenceURI = "/elements/metadata"

// Metadata is the individual Intelligence result
type Metadata struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Key                string `json:"key"`
	Image              string `json:"image"`
	DisplayOrder       int    `json:"displayOrder"`
	Active             bool   `json:"active"`
	Beta               bool   `json:"beta"`
	Description        string `json:"description"`
	Transformations    bool   `json:"transformations"`
	ElementType        string `json:"elementType"`
	Churros            bool   `json:"churros"`
	ElementClass       string `json:"elementClass"`
	NormalizedPaging   bool   `json:"normalizedPaging"`
	SwaggerValidated   bool   `json:"swaggerValidated"`
	Cloneable          bool   `json:"cloneable"`
	AuthenticationType string `json:"authenticationType"`
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
	Discovery struct {
		CustomFields          bool `json:"customFields"`
		CustomObjects         bool `json:"customObjects"`
		EndpointCustomFields  bool `json:"endpointCustomFields"`
		EndpointCustomObjects bool `json:"endpointCustomObjects"`
	} `json:"discovery"`
	Bulk struct {
		Upload   bool `json:"upload"`
		Download bool `json:"download"`
	} `json:"bulk"`
	Usage struct {
		InstanceCount int
		CustomerCount int
		Traffic       int
	} `json:"usage"`
	API struct {
		Type        string `json:"type"`
		ContentType string `json:"contentType"`
	} `json:"api"`
	Hub string `json:"hub"`
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

	url := fmt.Sprintf("%s%s", base, IntelligenceURI)

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

// OutputMetadataTable writes out either a tabular or csv view of the metadata
func OutputMetadataTable(metadatabytes []byte, orderBy string, filterBy string, asCsv bool) {
	var intelligence Intelligence
	err := json.Unmarshal(metadatabytes, &intelligence)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sort.Sort(intelligence)
	if orderBy == "customers" {
		sort.Sort(ByCustomerCount(intelligence))
	} else if orderBy == "hub" {
		sort.Sort(ByIntHub(intelligence))
	} else if orderBy == "name" {
		sort.Sort(ByIntName(intelligence))
	} else if orderBy == "instances" {
		sort.Sort(ByInstanceCount(intelligence))
	} else if orderBy == "traffic" {
		sort.Sort(ByTraffic(intelligence))
	} else if orderBy == "api" {
		sort.Sort(ByAPIType(intelligence))
	} else if orderBy == "authn" {
		sort.Sort(ByAuthn(intelligence))
	}

	data := [][]string{}
	for _, v := range intelligence {
		//configcount := strconv.Itoa(len(v.Configuration))
		data = append(data, []string{
			strconv.Itoa(v.ID),
			v.Key,
			v.Name,
			v.Hub,
			v.API.Type,
			v.AuthenticationType,
			strconv.FormatBool(v.Transformations),
			strconv.FormatBool(v.Active),
			strconv.FormatBool(v.Beta),
			v.ElementClass,
			strconv.Itoa(v.Usage.Traffic),
			strconv.Itoa(v.Usage.CustomerCount),
			strconv.Itoa(v.Usage.InstanceCount),
		})
	}

	if asCsv == true {
		w := csv.NewWriter(os.Stdout)
		for _, record := range data {
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}
		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatal(err)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Key", "Name", "Hub", "API", "Authn", "Transforms", "Hidden", "Beta", "Class", "Traffic", "Customers", "Instances"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	}
}

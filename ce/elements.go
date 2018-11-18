package ce

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/moul/http2curl"
	"github.com/olekukonko/tablewriter"
)

const (
	ElementsURI           = "/elements"
	ElementsKeysURI       = "/elements/keys"
	ElementsDocsFormatURI = "/elements/%s/docs"
	// ElementsLBDocsFormatURI if `force`, get latest swagger from repo and generate new lbdocs data
	// `version`` will retrieve a specific version (force is not applicable with `version`)
	ElementsLBDocsFormatURI              = "/elements/%s/lbdocs"
	ElementsMetadataFormatURI            = "/elements/%s/metadata"
	ElementFormatURI                     = "/elements/%s"
	ElementInstancesFormatURI            = "/elements/%s/instances"
	ElementInstanceFormatURI             = "/elements/%s/instances/%s"
	ElementInstancesFormatURINoReauthURI = "/elements/%s/instances/%s?reAuthenticate=false"
	ElementsOAuthTokenFormatURI          = "/elements/%s/oauth/token"
	ElementsOAuthURLTokenFormatURI       = "/elements/%s/oauth/url"
	ElementValidateModelsFormatURI       = "/elements/%s/validate"
	ElementsDenyList                     = "/customers/elements/blacklist"
)

// Element represents an Element resulting from a global Element list
type Element struct {
	ID                       int                      `json:"id,omitempty"`
	Name                     string                   `json:"name,omitempty"`
	Key                      string                   `json:"key,omitempty"`
	Description              string                   `json:"description,omitempty"`
	Image                    string                   `json:"image,omitempty"`
	Active                   bool                     `json:"active,omitempty"`
	Deleted                  bool                     `json:"deleted,omitempty"`
	OAuth                    bool                     `json:"typeOauth,omitempty"`
	TrialAccount             bool                     `json:"trialAccount,omitempty"`
	ConfigurationDescription string                   `json:"configuration_description,omitempty"`
	SignupURL                string                   `json:"signup_url,omitempty"`
	DefaultTransformations   []InstanceTransformation `json:"default_transformations,omitempty"`
	Configuration            []ElementConfiguration   `json:"configuration,omitempty"`
	Resources                []ElementResources       `json:"resources,omitempty"`
	Objects                  interface{}              `json:"objects,omitempty"`
	TransformationsEnabled   bool                     `json:"transformationsEnabled,omitempty"`
	BulkDownloadEnabled      bool                     `json:"bulkDownloadEnabled,omitempty"`
	BulkUploadEnabled        bool                     `json:"bulkUploadEnabled,omitempty"`
	Cloneable                bool                     `json:"cloneable,omitempty"`
	Extendable               bool                     `json:"extendable,omitempty"`
	Beta                     bool                     `json:"beta,omitempty"`
	Authentication           ElementAuthentication    `json:"authentication,omitempty"`
	Hooks                    interface{}              `json:"hooks,omitempty"`
	Extended                 bool                     `json:"extended,omitempty"`
	Hub                      string                   `json:"hub,omitempty"`
	ProtocolType             string                   `json:"protocolType,omitempty"`
	Parameters               interface{}              `json:"parameters,omitempty"`
	Private                  bool                     `json:"private,omitempty"`
	HookName                 string                   `json:"hookName,omitempty"`
}

// ElementConfiguration represents an element's configuration
type ElementConfiguration struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Key             string `json:"key"`
	Description     string `json:"description"`
	DefaultValue    string `json:"defaultValue"`
	ResellerConfig  bool   `json:"resellerConfig"`
	CompanyConfig   bool   `json:"companyConfig"`
	Active          bool   `json:"active"`
	Internal        bool   `json:"internal"`
	GroupControl    bool   `json:"groupControl"`
	DisplayOrder    int    `json:"displayOrder"`
	Type            string `json:"type"`
	HideFromConsole bool   `json:"hideFromConsole"`
	Required        bool   `json:"required"`
}

// ElementResources represents an Element's resources
type ElementResources struct {
	ID             int           `json:"id,omitempty"`
	CreatedDate    string        `json:"createdDate,omitempty"`
	UpdatedDate    string        `json:"updateDate,omitempty"`
	Description    string        `json:"description,omitempty"`
	Path           string        `json:"path,omitempty"`
	VendorPath     string        `json:"vendorPath,omitempty"`
	Method         string        `json:"method,omitempty"`
	VendorMethod   string        `json:"vendorMethod,omitempty"`
	Parameters     interface{}   `json:"parameters,omitempty"`
	Type           string        `json:"type,omitempty"`
	Hooks          []interface{} `json:"hooks,omitempty"`
	Response       interface{}   `json:"response,omitempty"`
	PaginationType string        `json:"paginationType,omitempty"`
	OwnerAccountID int           `json:"ownerAccountId,omitempty"`
}

// InstanceTransformation is a transformation for a field on an Element Instance
type InstanceTransformation struct {
	Name       string `json:"name,omitempty"`
	VendorName string `json:"vendor_name,omitempty"`
}

// ElementInstance represents an Element Instance
type ElementInstance struct {
	ID                     int                   `json:"id,omitempty"`
	Name                   string                `json:"name,omitempty"`
	Token                  string                `json:"token,omitempty"`
	Element                Element               `json:"element"`
	Tags                   []string              `json:"tags"`
	Valid                  bool                  `json:"valid"`
	Disabled               bool                  `json:"disabled"`
	Configuration          InstanceConfiguration `json:"configuration,omitempty"`
	EventsEnabled          bool                  `json:"eventsEnabled"`
	ExternalAuthentication string                `json:"externalAuthentication"`
	User                   InstanceUser          `json:"user"`
	TraceLoggingEnabled    bool                  `json:"traceLoggingEnabled"`
}

// InstanceConfiguration is the configuration associated with an Element Instance
// - this may be too variable to capture in a structure, may want to leave as-is interface
type InstanceConfiguration struct {
	BaseURL                         string      `json:"base_url,omitempty"`
	EventNotificationSubscriptionID string      `json:"event_notification_subscription_id,omitempty"`
	EventMetadata                   string      `json:"event_metadata,omitempty"`
	EventVendorType                 string      `json:"event.vendor.type,omitempty"`
	EventNotificationSignatureKey   interface{} `json:"event.notification.signature.key,omitempty"`
	EventNotificationEnabled        string      `json:"event.notification.enabled,omitempty"`
	EventObjects                    string      `json:"event.objects,omitempty"`
	EventHelperKey                  string      `json:"event.helper.key,omitempty"`
	EventPollerRefreshInterval      string      `json:"event.poller.refresh_interval,omitempty"`
	EventPollerURLs                 interface{} `json:"event.poller.urls,omitempty"`
	EventNotificationCallbackURL    string      `json:"event.notification.callback.url,omitempty"`
	FilterResponseNulls             string      `json:"filter.response.nulls,omitempty"`
	BulkQueryDateMask               string      `json:"bulk.query.date_mask,omitempty"`
	BulkAttributeCreatedTime        string      `json:"bulk.attribute.created_time,omitempty"`
	BulkAttributeModifiedTime       string      `json:"bulk.attribute.modified_time,omitempty"`
	CRMSessionRefreshTime           string      `json:"crm.session.refresh.time,omitempty"`
	CRMSessionRefreshInterval       string      `json:"crm.session.refresh.interval,omitempty"`
	OAuthCallbackURL                string      `json:"oauth.callback.url,omitempty"`
	OAuthUserRefreshToken           string      `json:"oauth.user.refresh_token,omitempty"`
	OAuthAPIKey                     string      `json:"oauth.api.key,omitempty"`
	OAuthAPISecret                  string      `json:"oauth.api.secret,omitempty"`
	OAuthScope                      string      `json:"oauth.scope,omitempty"`
	OAuthUserToken                  string      `json:"oauth.user.token,omitempty"`
	OAuthUserRefreshTime            string      `json:"oauth.user.refresh_time,omitempty"`
	SessionID                       string      `json:"session.id,omitempty"`
	SFDCUserID                      string      `json:"sfdc_user_id,omitempty"`
	SFDCPassword                    string      `json:"sfdc_password,omitempty"`
	SFDCAPIKey                      string      `json:"sfdc.api.key,omitempty"`
	SFDCAPISecret                   string      `json:"sfdc.api.secret,omitempty"`
	SFDCRevokeURL                   string      `json:"sfdc.revoke.url,omitempty"`
	SFDCSessionSignature            string      `json:"sfdc.session.signature,omitempty"`
	SFDCSessionInstanceURL          string      `json:"sfdc.session.instance.url,omitempty"`
	SFDCSessionAPIVersionURI        string      `json:"sfdc.session.api.version.uri,omitempty"`
	SFDCTokenURL                    string      `json:"sfdc.token.url,omitempty"`
	SFDCUserIDURL                   string      `json:"sfdc.user.id.url,omitempty"`
	SFDCUsername                    string      `json:"sfdc.username,omitempty"`
	SFDCUserDisplayName             string      `json:"sfdc.user.display.name,omitempty"`
	SFDCSecurityToken               string      `json:"sfdc.security.token,omitempty"`
}

// InstanceUser is the user associated with an Element Instance
type InstanceUser struct {
	ID           int    `json:"id,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
}

// ElementAuthentication represents an element's authentication
type ElementAuthentication struct {
	Type string `json:"type,omitempty"`
}

// Elements is a struct container for a list of elements, used in sorting
type Elements []Element

func (elements Elements) Len() int           { return len(elements) }
func (elements Elements) Less(i, j int) bool { return elements[i].ID < elements[j].ID }
func (elements Elements) Swap(i, j int)      { elements[i], elements[j] = elements[j], elements[i] }

// ByHub implements sort.Interface for Elements
type ByHub []Element

func (e ByHub) Len() int           { return len(e) }
func (e ByHub) Less(i, j int) bool { return e[i].Hub < e[j].Hub }
func (e ByHub) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// ByName implements sort.Interface for Elements
type ByName []Element

func (e ByName) Len() int           { return len(e) }
func (e ByName) Less(i, j int) bool { return strings.ToLower(e[i].Name) < strings.ToLower(e[j].Name) }
func (e ByName) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// DeleteElement deletes an Element on the Platform
func DeleteElement(base, auth string, elementID int) ([]byte, int, string, error) {
	elementIDstr := strconv.Itoa(elementID)
	url := fmt.Sprintf("%s%s", base, fmt.Sprintf(ElementFormatURI, elementIDstr))
	return Execute("DELETE", url, auth)
}

// ImportElement imports an Element to the Platform
func ImportElement(base, auth string, element Element) ([]byte, int, string, error) {
	var bodybytes []byte
	elementBytes, err := json.Marshal(element)
	if err != nil {
		return bodybytes, -1, "", err
	}
	url := fmt.Sprintf("%s%s", base, ElementsURI)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(elementBytes))
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	curlCmd, _ := http2curl.GetCurlCommand(req)
	curl := fmt.Sprintf("%s", curlCmd)
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Cannot process response", err.Error())
		return bodybytes, -1, curl, err
	}
	bodybytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return bodybytes, resp.StatusCode, curl, nil
}

// GetAllElements returns all Elements as bytes
func GetAllElements(base, auth string) ([]byte, int, string, error) {

	var bodybytes []byte

	url := fmt.Sprintf("%s%s", base, ElementsURI)

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

// GetElementModelValidation validates the models for a provided Element id
func GetElementModelValidation(base, auth, elementid string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(ElementValidateModelsFormatURI, elementid),
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

// GetElementLBDocs returns the LoopBack model document for this
// force is a boolean, and will force a refresh of the latest version
// version is an int, referring to a version number of LBDocs; version has no effect on force
func GetElementLBDocs(base, auth, elementid string, force bool, version string) ([]byte, int, string, error) {
	urlstr := fmt.Sprintf("%s%s", base, fmt.Sprintf(ElementsLBDocsFormatURI, elementid))
	u, _ := url.Parse(urlstr)
	q := u.Query()
	if force {
		q.Set("force", fmt.Sprintf("%v", force))
	}
	if version != "" {
		q.Set("version", version)
	}
	u.RawQuery = q.Encode()
	return Execute("GET", u.String(), auth)
}

// GetElementOAI returns the OAI for an Element id
func GetElementOAI(base, auth, elementid string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(ElementsDocsFormatURI, elementid),
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

// GetExportElement returns the JSON of the Element
func GetExportElement(base, auth, elementid string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(ElementFormatURI, elementid),
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

// GetElementMetadata returns the metadata for an Element id
func GetElementMetadata(base, auth, elementid string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(ElementsMetadataFormatURI, elementid),
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

// GetElementInstances returns the instances for an Element key/id
func GetElementInstances(base, auth, elementid string) ([]byte, int, string, error) {
	var bodybytes []byte
	url := fmt.Sprintf("%s%s",
		base,
		fmt.Sprintf(ElementInstancesFormatURI, elementid),
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

// AddToElementsDenyList adds a list of Element keys to the deny list
// requires Customer Admin privileges
func AddToElementsDenyList(base, auth string, elementkeys []string) ([]byte, int, string, error) {
	urlstr := fmt.Sprintf("%s%s", base, ElementsDenyList)
	u, _ := url.Parse(urlstr)
	elementarray, err := json.Marshal(&elementkeys)
	if err != nil {
		return nil, -1, "", err
	}
	return ExecuteWithBody("PUT", u.String(), auth, elementarray)
}

// ResetElementsDenyList clears out the Element deny list
func ResetElementsDenyList(base, auth string) ([]byte, int, string, error) {
	urlstr := fmt.Sprintf("%s%s", base, ElementsDenyList)
	u, _ := url.Parse(urlstr)
	return Execute("DELETE", u.String(), auth)
}

// ElementKeyToID returns the ID (int) of an Element Key (string)
func ElementKeyToID(key string, profilemap map[string]string) (int, error) {
	var elementid int
	elementid, err := strconv.Atoi(key)
	if err != nil {

		// Get elements
		bodybytes, _, _, err := GetAllElements(profilemap["base"], profilemap["auth"])
		if err != nil {
			return elementid, err
		}
		var elements Elements
		err = json.Unmarshal(bodybytes, &elements)
		if err != nil {
			return elementid, err
		}
		// find Element ID given Element key
		for _, v := range elements {
			if v.Key == key {
				elementid = v.ID
			}
		}
		if elementid == 0 {
			err := fmt.Errorf("unable to find Element ID for Element Key %s", key)
			return elementid, err
		}

	}
	return elementid, nil
}

// OutputElementInstancesTable writes out a tabular view of the instances list
func OutputElementInstancesTable(instancesbytes []byte) error {
	var instances []ElementInstance
	err := json.Unmarshal(instancesbytes, &instances)
	if err != nil {
		return err
	}

	data := [][]string{}
	for _, i := range instances {
		data = append(data, []string{
			strconv.Itoa(i.ID),
			i.Element.Key,
			i.Name,
			strconv.FormatBool(i.Valid),
			strconv.FormatBool(i.Disabled),
			strconv.FormatBool(i.EventsEnabled),
			//strconv.FormatBool(i.)
			fmt.Sprintf("%s", i.Tags),
			i.Token,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Key", "Name", "Valid", "Disabled", "Events", "Tags", "Token"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}

// FilterCustomElements returns only the custom elements
func FilterCustomElements(elementsbytes []byte) ([]byte, error) {
	var elements Elements
	err := json.Unmarshal(elementsbytes, &elements)
	if err != nil {
		fmt.Println(err.Error())
		return elementsbytes, err
	}
	var filteredElements Elements
	for _, v := range elements {
		if v.Private == true {
			filteredElements = append(filteredElements, v)
		}
	}
	elementsbytes, err = json.Marshal(filteredElements)
	if err != nil {
		return elementsbytes, err
	}
	return elementsbytes, nil
}

// FilterElementFromList returns an array of Elements whose key matches the input
func FilterElementFromList(keyfilter string, elementsbytes []byte) ([]byte, error) {
	var elements Elements
	err := json.Unmarshal(elementsbytes, &elements)
	if err != nil {
		fmt.Println(err.Error())
		return elementsbytes, err
	}
	var filteredElements Elements
	for _, v := range elements {
		if v.Key == keyfilter {
			filteredElements = append(filteredElements, v)
		}
	}
	elementsbytes, err = json.Marshal(filteredElements)
	if err != nil {
		return elementsbytes, err
	}
	return elementsbytes, nil
}

// OutputElementsTable writes out a tabular view of the elements list
func OutputElementsTable(elementsbytes []byte, orderBy string, filterBy string) {
	var elements Elements
	err := json.Unmarshal(elementsbytes, &elements)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sort.Sort(elements)
	if orderBy == "name" {
		sort.Sort(ByName(elements))
	} else if orderBy == "hub" {
		sort.Sort(ByHub(elements))
	}
	data := [][]string{}
	for _, v := range elements {
		authtype := v.Authentication.Type
		configcount := strconv.Itoa(len(v.Configuration))
		data = append(data, []string{
			strconv.Itoa(v.ID),
			v.Key,
			v.Name,
			v.Hub,
			authtype,
			configcount,
			strconv.FormatBool(v.Private),
			strconv.FormatBool(v.Active),
			strconv.FormatBool(v.Extendable),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Key", "Name", "Hub", "Auth", "Configs", "Private", "Active", "Extendable"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}

// OutputElementsTableAsCSV writes out a csv view of the elements list
func OutputElementsTableAsCSV(elementsbytes []byte, orderBy string, filterBy string) {
	var elements Elements
	err := json.Unmarshal(elementsbytes, &elements)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sort.Sort(elements)
	if orderBy == "name" {
		sort.Sort(ByName(elements))
	} else if orderBy == "hub" {
		sort.Sort(ByHub(elements))
	}
	data := [][]string{}
	for _, v := range elements {
		authtype := v.Authentication.Type
		configcount := strconv.Itoa(len(v.Configuration))
		data = append(data, []string{
			strconv.Itoa(v.ID),
			v.Key,
			v.Name,
			v.Hub,
			authtype,
			configcount,
			strconv.FormatBool(v.Private),
			strconv.FormatBool(v.Active),
			strconv.FormatBool(v.Extendable),
		})
	}

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

}

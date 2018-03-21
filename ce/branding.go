package ce

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	// BrandingURI can be used with GET, PUT, DELETE
	BrandingURI = "/organizations/branding"
	// BrandingLogoURI can be used with PATCH
	BrandingLogoURI = `/organizations/branding/logo`
	// BrandingFaviconURI can be used with PATCH
	BrandingFaviconURI = `/organizations/branding/favicon`
)

// BrandingConfig represents the configuration of a Platform's branding
// TODO generate this struct via a GET to /organizations/branding
// as the configuration values change frequently; possibly during test suite
type BrandingConfig struct {
	HeaderFont                     string `json:"headerFont"`
	HeaderColor                    string `json:"headerColor"`
	BodyFont                       string `json:"bodyFont"`
	BodyColor                      string `json:"bodyColor"`
	Logo                           string `json:"logo"`
	Favicon                        string `json:"favicon"`
	ThemePrimaryColor              string `json:"themePrimaryColor"`
	ThemeSecondaryColor            string `json:"themeSecondaryColor"`
	ThemeHighlightColor            string `json:"themeHighlightColor"`
	ButtonPrimaryBackgroundColor   string `json:"buttonPrimaryBackgroundColor"`
	ButtonPrimaryTextColor         string `json:"buttonPrimaryTextColor"`
	ButtonSecondaryBackgroundColor string `json:"buttonSecondaryBackgroundColor"`
	ButtonSecondaryTextColor       string `json:"buttonSecondaryTextColor"`
	ButtonDeleteBackgroundColor    string `json:"buttonDeleteBackgroundColor"`
	ButtonDeleteTextColor          string `json:"buttonDeleteTextColor"`
	LogoBackgroundColor            string `json:"logoBackgroundColor"`
	TopBarBackgroundColor          string `json:"topBarBackgroundColor"`
	NavigationBackgroundColor      string `json:"navigationBackgroundColor"`
	ContextBackgroundColor         string `json:"contextBackgroundColor"`
	CardHeaderColor                string `json:"cardHeaderColor"`
	CardBackground                 string `json:"cardBackground"`
	CardMenuBackground             string `json:"cardMenuBackground"`
	CardMenuLinkColor              string `json:"cardMenuLinkColor"`
	NavigationLinkBackground       string `json:"navigationLinkBackground"`
	NavigationLinkForeground       string `json:"navigationLinkForeground"`
	NavigationLinkBackgroundHover  string `json:"navigationLinkBackgroundHover"`
	NavigationLinkForegroundHover  string `json:"navigationLinkForegroundHover"`
	NavigationLinkBackgroundActive string `json:"navigationLinkBackgroundActive"`
	NavigationLinkForegroundActive string `json:"navigationLinkForegroundActive"`
	TopBarNavigationColor          string `json:"topBarNavigationColor"`
	TableHeaderBackground          string `json:"tableHeaderBackground"`
	TableHeaderForeground          string `json:"tableHeaderForeground"`
	TableBodyBackground            string `json:"tableBodyBackground"`
	TableBodyForeground            string `json:"tableBodyForeground"`
	IntercomEnabled                bool   `json:"intercomEnabled"`
	PendoEnabled                   bool   `json:"pendoEnabled"`
	DocumentationURL               string `json:"documentationUrl"`
	ElementsEnabled                bool   `json:"elementsEnabled"`
	InstanceEnabled                bool   `json:"instanceEnabled"`
	FormulasEnabled                bool   `json:"formulasEnabled"`
	VirtualDataEnabled             bool   `json:"virtualDataEnabled"`
	ReportsEnabled                 bool   `json:"reportsEnabled"`
	NavigationIconPosition         string `json:"navigationIconPosition"`
	NavigationIconSize             string `json:"navigationIconSize"`
	NavigationLabelSize            string `json:"navigationLabelSize"`
}

// DefaultBranding represents the defaults of CE branding
var DefaultBranding = BrandingConfig{
	HeaderFont:                     "museo-sans",
	HeaderColor:                    "#4d82bf",
	BodyFont:                       "Open Sans",
	BodyColor:                      "#333333",
	ThemePrimaryColor:              "#4d82bf",
	ThemeSecondaryColor:            "#44c8f5",
	ThemeHighlightColor:            "#761299",
	ButtonPrimaryBackgroundColor:   "#4d82bf",
	ButtonPrimaryTextColor:         "#ffffff",
	ButtonSecondaryBackgroundColor: "#44c8f5",
	ButtonSecondaryTextColor:       "#ffffff",
	ButtonDeleteBackgroundColor:    "#FF4E4E",
	ButtonDeleteTextColor:          "#ffffff",
	LogoBackgroundColor:            "#44c8f5",
	TopBarBackgroundColor:          "#ffffff",
	NavigationBackgroundColor:      "#172330",
	ContextBackgroundColor:         "#edf1f2",
	CardHeaderColor:                "#4d82bf",
	CardBackground:                 "#ffffff",
	CardMenuBackground:             "#d1d1d1",
	CardMenuLinkColor:              "#172330",
	NavigationLinkBackground:       "#172330",
	NavigationLinkForeground:       "#ffffff",
	NavigationLinkBackgroundHover:  "#303a47",
	NavigationLinkForegroundHover:  "#ffffff",
	NavigationLinkForegroundActive: "#44c8f5",
	NavigationLinkBackgroundActive: "#101922",
	TopBarNavigationColor:          "#172330",
	TableHeaderBackground:          "#172330",
	TableBodyBackground:            "#ffffff",
	TableHeaderForeground:          "#ffffff",
	TableBodyForeground:            "#333333",
	NavigationIconSize:             "20px",
	NavigationLabelSize:            "9px",
}

// GetBranding returns the Platform's branding
func GetBranding(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		BrandingURI,
	)
	if debug {
		log.Println("Retrieving Platform branding ...")
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
	return bodybytes, status, curlcmd, nil
}

// SetBranding sets branding on the Platform, given a JSON object
func SetBranding(base, auth string, branding interface{}, debug bool) ([]byte, int, string, error) {
	var bodybytes []byte
	var curlcmd string
	var status int
	url := fmt.Sprintf("%s%s", base,
		BrandingURI,
	)
	requestbytes, err := json.Marshal(branding)
	if err != nil {
		return bodybytes, -1, curlcmd, err
	}
	if debug {
		log.Println("Updating Platform branding ...")
		log.Println("PUT", url)
	}
	bodybytes, status, curlcmd, err = ExecuteWithBody("PUT", url, auth, requestbytes)
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
		return bodybytes, status, curlcmd, fmt.Errorf("Non-200 status %v", status)
	}
	return bodybytes, status, curlcmd, nil
}

// ResetBranding returns the Platform branding to the default
func ResetBranding(base, auth string, debug bool) ([]byte, int, string, error) {
	url := fmt.Sprintf("%s%s",
		base,
		BrandingURI,
	)
	if debug {
		log.Println("Resetting Platform branding ...")
		log.Println("DELETE", url)
	}
	bodybytes, status, curlcmd, err := Execute("DELETE", url, auth)
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
	return bodybytes, status, curlcmd, nil
}

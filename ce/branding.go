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
type BrandingConfig struct {
	HeaderFont                     string `json:"headerFont,omitempty"`
	HeaderColor                    string `json:"headerColor,omitempty"`
	BodyFont                       string `json:"bodyFont,omitempty"`
	BodyColor                      string `json:"bodyColor,omitempty"`
	Logo                           string `json:"logo,omitempty"`
	Favicon                        string `json:"favicon,omitempty"`
	ThemePrimaryColor              string `json:"themePrimaryColor,omitempty"`
	ThemeHighlightColor            string `json:"themeHighlightColor,omitempty"`
	ButtonPrimaryBackgroundColor   string `json:"buttonPrimaryBackgroundColor,omitempty"`
	ButtonPrimaryTextColor         string `json:"buttonPrimaryTextColor,omitempty"`
	ButtonSecondaryBackgroundColor string `json:"buttonSecondaryBackgroundColor,omitempty"`
	ButtonSecondaryTextColor       string `json:"buttonSecondaryTextColor,omitempty"`
	ButtonDeleteBackgroundColor    string `json:"buttonDeleteBackgroundColor,omitempty"`
	ButtonDeleteTextColor          string `json:"buttonDeleteTextColor,omitempty"`
	LogoBackgroundColor            string `json:"logoBackgroundColor,omitempty"`
	TopBarBackgroundColor          string `json:"topBarBackgroundColor,omitempty"`
	NavigationBackgroundColor      string `json:"navigationBackgroundColor,omitempty"`
	ContextBackgroundColor         string `json:"contextBackgroundColor,omitempty"`
	Errors                         struct {
		Image string `json:"image,omitempty"`
	} `json:"errors,omitempty"`
	ShouldUploadImage bool `json:"shouldUploadImage,omitempty"`
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

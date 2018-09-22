package ce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client is a CE client
type Client struct {
	BaseURL      *url.URL
	Organization string
	User         string
	Element      string

	httpClient *http.Client
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	c.httpClient = &http.Client{}
	rel := &url.URL{Path: fmt.Sprintf("/elements/api-v2%s", path)}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Organization %s, User %s", c.Organization, c.User))

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

// ListElements lists Elements on the Platform
func (c *Client) ListElements() ([]Element, error) {
	req, err := c.newRequest("GET", "/elements", nil)
	if err != nil {
		return nil, err
	}
	var elements []Element
	_, err = c.do(req, &elements)
	return elements, err
}

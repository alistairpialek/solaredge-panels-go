package solaredge

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

const (
	defaultBaseURL = "https://monitoring.solaredge.com/solaredge-apigw/api/"
	version        = "0.0.1"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	BaseURL       *url.URL
	UserAgent     string
	Authorization string

	common service

	client HTTPClient
	Site   *SiteService
}

type service struct {
	client *Client
}

func UserAgent() string {
	return "solaredge-panels-go/" + version + " (" + runtime.GOOS + "; " + runtime.GOARCH + ")"
}

// NewClient returns a new SolarEdge API client. If a nil httpClient is
// provided, a new http.Client will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewClient(httpClient HTTPClient, username string, password string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	c := &Client{
		client:        httpClient,
		BaseURL:       SetBaseURL(),
		UserAgent:     UserAgent(),
		Authorization: createAuthorizationString(username, password),
	}
	c.common.client = c
	c.Site = (*SiteService)(&c.common)

	return c
}

func SetBaseURL() *url.URL {
	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		panic(err)
	}

	return baseURL
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.Authorization)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	// Show the request.
	//log.Printf("request: %+v", req)

	// Make the response.
	resp, err := c.client.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	// Show the response.
	//log.Printf("response: %+v", resp)

	// Read the response body now so that we can show it for debugging purposes.
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	//log.Printf("response body: %s", string(bodyBytes))

	// Use unmarshal rather than decode since the reader now has already been read.
	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return resp, err
	}

	// Check the response status code is 2xx
	if resp.StatusCode >= 299 {
		return resp, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func createAuthorizationString(username string, password string) string {
	authString := fmt.Sprintf("%s:%s", username, password)
	b64AuthString := base64.StdEncoding.EncodeToString([]byte(authString))

	return fmt.Sprintf("Basic %s", b64AuthString)
}

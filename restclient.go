package xiqrestclient

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// RESTClient encapsulates the actual HTTP client that communicates with XIQ.
// Use New() to obtain an usable instance. All fields should be treated as read-only; functions are provided where changes shall be possible.
type RESTClient struct {
	httpClient      http.Client
	HTTPHost        string
	HTTPPort        uint
	HTTPTimeout     uint
	XIQOwnerID      string
	XIQAccessToken  string
	XIQClientID     string
	XIQClientSecret string
	XIQRedirectURI  string
	UserAgent       string
}

// New is used to create an usable instance of RESTClient.
// By default a new instance will use HTTPS to port 443 with strict certificate checking. The HTTP timeout is set to 5 seconds. Authentication must be set manually before trying to send a query to XIQ.
func New(host string, ownerID string) RESTClient {
	var c RESTClient

	c.httpClient = http.Client{}
	c.HTTPHost = host
	c.HTTPPort = 443
	c.XIQOwnerID = ownerID
	c.SetTimeout(5)
	c.SetUserAgent(fmt.Sprintf("%s/%s", moduleName, moduleVersion))

	return c
}

// SetTimeout sets the HTTP timeout in seconds for the RESTClient instance.
func (c *RESTClient) SetTimeout(seconds uint) error {
	if httpMinTimeout <= seconds && httpMaxTimeout >= seconds {
		c.httpClient.Timeout = time.Second * time.Duration(seconds)
		return nil
	}
	return fmt.Errorf("timeout out of range (%d - %d)", httpMinTimeout, httpMaxTimeout)
}

// SetUserAgent sets the User-Agent HTTP header.
func (c *RESTClient) SetUserAgent(ua string) {
	c.UserAgent = ua
}

// SetAuth sets the authentication credentials.
func (c *RESTClient) SetAuth(accessToken string, clientID string, clientSecret string, redirectURI string) {
	c.XIQAccessToken = accessToken
	c.XIQClientID = clientID
	c.XIQClientSecret = clientSecret
	c.XIQRedirectURI = redirectURI
}

// SanitizeEndpoint prepares the provided API endpoint for concatenation.
func SanitizeEndpoint(endpoint *string) {
	if !strings.HasPrefix(*endpoint, "/") {
		*endpoint = fmt.Sprintf("/%s", *endpoint)
	}
	if !strings.HasPrefix(*endpoint, "/xapi") {
		*endpoint = fmt.Sprintf("/xapi%s", *endpoint)
	}
}

// SetRequestHeaders sets the usual headers required for requests to XIQ.
func SetRequestHeaders(client *RESTClient, req *http.Request, payload *[]byte) {
	req.Header.Set("User-Agent", client.UserAgent)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", jsonMimeType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.XIQAccessToken))
	req.Header.Set("X-AH-API-CLIENT-ID", client.XIQClientID)
	req.Header.Set("X-AH-API-CLIENT-SECRET", client.XIQClientSecret)
	req.Header.Set("X-AH-API-CLIENT-REDIRECT-URI", client.XIQRedirectURI)
	if payload != nil {
		req.Header.Set("Content-Type", jsonMimeType)
	}
}

// GetRequest returns a prepared HTTP GET request instance.
func (c *RESTClient) GetRequest(endpoint string) (*http.Request, error) {
	SanitizeEndpoint(&endpoint)
	endpointURL := fmt.Sprintf("https://%s:%d%s?ownerId=%s", c.HTTPHost, c.HTTPPort, endpoint, c.XIQOwnerID)

	req, reqErr := http.NewRequest(http.MethodGet, endpointURL, nil)
	if reqErr != nil {
		return req, fmt.Errorf("could not create request: %s", reqErr)
	}
	SetRequestHeaders(c, req, nil)

	return req, nil
}

// PerformRequest sends a request to XIQ and returns the result.
func (c *RESTClient) PerformRequest(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

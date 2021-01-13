package flow

const (
	ProductionURL = "https://www.flow.cl/api"
	SandboxURL = "https://sandbox.flow.cl/api"
)

// Client contains the authentication data and the canonical URL used for requests to the Flow API/
type Client struct {
	// APIKey is the access key provided by Flow.
	APIKey string

	// SecretKey is uses to sign the API requests.
	SecretKey string

	// URL is the base URL to be used in the requests. Can be ProductionURL or SandboxURL.
	URL string
}

// NewClient creates a *Client with the given keys. By default it's set to sandbox mode.
func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		URL:       SandboxURL,
	}
}

// SetProduction sets the Client's URL to the production base URL.
func (c *Client) SetProduction()  {
	c.URL = ProductionURL
}

// SetProduction sets the Client's URL to the sandbox base URL.
func (c *Client) SetSandbox()  {
	c.URL = SandboxURL
}

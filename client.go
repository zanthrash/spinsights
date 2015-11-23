package spinsights

import (
	"net/http"
	"net/url"
)

var DefalutClient = NewClient(nil)

type Client struct {
	BaseUrl *url.URL
	ApplicationPath string
	UserAgent string
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		cloned := *http.DefaultClient
		httpClient = &cloned
	}

	c := &Client{
		BaseUrl: &url.URL{
			Scheme:"https",
			Host:"spinnaker-api.prod.netflix.net",
		},
		ApplicationPath:"/application/",
		UserAgent:"spinsights.go",
		httpClient: httpClient,
	}

	return c
}

func (c *Client) NewRequest(s string) (*http.Request, error) {
	rel, err := url.Parse(c.BasePath + s)
	if err != nil {
		return nil, err
	}



}

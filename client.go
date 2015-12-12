package spinsights

import (
	"net/url"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"fmt"
	"github.com/gizak/termui/debug"
	"bytes"
)

var DefalutClient = NewClient(nil)



func (c *Client) GetExecutionById(id string) (*Execution, error) {
	rel, err := url.Parse("/pipelines/" + id)

	if err != nil {
		return &Execution{}, err
	}

	url := c.BaseUrl.ResolveReference(rel)

	_, body, _:= gorequest.New().Get(url.String()).End()


	var execution *Execution
	json.Unmarshal([]byte(body), &execution)

	return execution, nil
}

func (c *Client) GetInstanceDetails(instance Instance) (*InstanceDetail, error) {
	urlString := fmt.Sprintf("/instances/%s/%s/%s", instance.Account, instance.Region, instance.InstanceId)
	rel, err := url.Parse(urlString)
	if err != nil {
		return &InstanceDetail{}, err
	}

	url := c.BaseUrl.ResolveReference(rel)
	_, body, _ := gorequest.New().Get(url.String()).End()

	var instanceDetail *InstanceDetail
	json.Unmarshal([]byte(body), &instanceDetail)
	return instanceDetail, nil
}

func (c *Client) InstanceSearch(instanceId string) ([]SearchResult, error) {
	rel, err := url.Parse("/search")
	if err != nil {
		return []SearchResult{}, err
	}
	q := rel.Query()
	q.Set("q", instanceId)
	q.Set("type", "instances")

	rel.RawQuery = q.Encode()

	url := c.BaseUrl.ResolveReference(rel)

	debug.Log(url.String())
	_, body, _ := gorequest.New().Get(url.String()).End()

	prettyPrintJson(body)
	var results []SearchResult
	json.Unmarshal([]byte(body), &results)
	return results, nil

}

func unmarshallExecutionJSON(jsonString string) *Execution  {
	var parsed Execution
	json.Unmarshal([]byte(jsonString), &parsed)
	return &parsed
}

func prettyPrintJson(jsonString string) {
	var prettyJson bytes.Buffer
	json.Indent(&prettyJson, []byte(jsonString), "", "    ")

	debug.Log(string(prettyJson.Bytes()))
}

func NewClient(agent *gorequest.SuperAgent) *Client {
	if agent == nil {
		cloned := *gorequest.New()
		agent = &cloned
	}

	c := &Client{
		BaseUrl: &url.URL{
			Scheme:"https",
			Host:"spinnaker-api.prod.netflix.net",
		},
		UserAgent:"spinsights.go",
		agent: agent,
	}

	return c
}


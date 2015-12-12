package spinsights

import (
	"net/url"
	"github.com/parnurzeal/gorequest"
)

type Client struct {
	BaseUrl *url.URL
	UserAgent string
	agent *gorequest.SuperAgent
}

type Execution struct {
	Application string
	Status string
	Name string
	StartTime int
	EndTime int
	ExecutingInstance string
	Stages []Stage
	Trigger Trigger
}

type Context struct {
	Exception Exception
}

type Exception struct {
	Details ExcptionDetails
}

type ExcptionDetails struct {
	Error string
	Errors []string
	StackTrace string
}

type Trigger struct {
	User string
	Type string
}

type Stage struct {
	Type string
	Name string
	Status string
	StartTime int
	EndTime int
	Tasks []Task
	Context Context
}

type Task struct {
	Id string
	Name string
	StartTime int
	EndTime int
	Status string
}

type SearchResult struct {
	Results []Instance
}

type Instance struct {
	Provider string
	Type string
	Account string
	Region string
	InstanceId string
	Application string
	Cluster string
	ServerGroup string
}

type InstanceDetail struct {
	Name string
	HealthState string
	PrivateIpAddress string
	InsightActions []InsightActions
}

type InsightActions struct {
	Url string
	Label string
}
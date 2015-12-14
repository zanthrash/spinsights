package spinsights

import (
	"net/url"
	"github.com/parnurzeal/gorequest"
	"strings"
	"fmt"
)

type Client struct {
	BaseUrl *url.URL
	UserAgent string
	agent *gorequest.SuperAgent
}

type Execution struct {
	Id string
	Application string
	Status string
	Name string
	StartTime int
	EndTime int
	PipelineConfigId string
	ExecutingInstance string
	Stages []Stage
	Trigger Trigger
	Context Context
}

func (exe *Execution) getScalingActivitiesUrls() []string {
	appName := exe.Application
	deployStages := exe.getDeployStages()

	urls := []string{}
	for i := range deployStages {
		stage := deployStages[i]
		region := stage.Context.Source.Region
		account := stage.Context.Source.Account
		serverGroup := stage.Context.getDeployServerGroupName()
		clusterName := stage.Context.getClusterName()
		url := fmt.Sprintf("https://spinnaker-api.prod.netflix.net/applications/%s/clusters/%s/%s/serverGroups/%s/scalingActivities?region=%s", appName, account, clusterName, serverGroup, region)
		urls = append(urls, url)
	}

	return urls
}

func (exe *Execution) getDeployStages() []Stage {
	deployStages := make([]Stage, 0)

	for i := range exe.Stages {
		if exe.Stages[i].Type == "deploy" {
			deployStages = append(deployStages, exe.Stages[i])
		}
	}

	return deployStages
}

type Context struct {
	Exception Exception
	Source Source
	DeployServerGroup map[string]interface{} `json:"deploy.server.groups"`
}

func (context *Context) getDeployRegion() string {
	if context.DeployServerGroup == nil {
		return ""
	}
	keys := make([]string, 0, len(context.DeployServerGroup))
	for k := range context.DeployServerGroup {
		keys = append(keys, k)
	}
	if len(keys) > 0 {
		return keys[0]
	}
	return ""
}

func (context *Context) getDeployServerGroupName() string {
	if context.DeployServerGroup == nil {
		return ""
	}

	key := context.getDeployRegion()
	if key != "" {
		m := context.DeployServerGroup[key].([]interface{})
		return m[0].(string)
	}
	return ""
}

func (context *Context) getClusterName() string {
	if serverGroup := context.getDeployServerGroupName(); serverGroup != "" {
		split := strings.Split(serverGroup, "-")
		return strings.Join(split[:len(split)-1], "-")
	}
	return ""
}



type Source struct {
	AsgName string
	Account string
	Region string
}

type Exception struct {
	Details ExceptionDetails
}

type ExceptionDetails struct {
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
	StartTime int64
	EndTime int64
	ParentStageId string
	Tasks []Task
	Context Context
}

type Task struct {
	Id string
	Name string
	StartTime int64
	EndTime int64
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

type AutoScalingActivity struct {
	ActivityId string
	AutoScalingGroupName string
	Description string
	Cause string
	StartTime int64
	EndTime int64
	StatusCode string
	Progress int
	Details string
}


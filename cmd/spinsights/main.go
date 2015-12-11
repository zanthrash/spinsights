package main

import (
	"os"
	"os/exec"
	"github.com/codegangsta/cli"
	"github.com/zanthrash/spinsights"
	"github.com/gizak/termui/debug"

)

func main() {
	app := cli.NewApp()
	app.Name = "spinsights"
	app.Usage = "Troubleshooting tool for spinnaker"
	app.EnableBashCompletion = true


	go func() { panic(debug.ListenAndServe()) }()

//	client := spinsights.DefalutClient
	var executionId string
	app.Commands = []cli.Command{
		{
			Name: "exec",
			Usage: "Takes a execution id and displays details about a pipeline execution",
			Aliases: []string{"e"},
			Action: func(c * cli.Context) {
				executionId = c.Args().First()
//				client.GetExecutionById(executionId)
				spinsights.RenderPipeline(executionId)
			},
		},
		{
			Name: "open",
			Aliases: []string{"o"},
			Action: func(c *cli.Context) {
				url := c.Args().First()
				cm := exec.Command("open", url)
				cm.Run()

			},
		},
	}

	app.Run(os.Args)
}

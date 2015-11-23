package main

import (
	"os"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "spinsights"
	app.Usage = "Troubleshooting tool for spinnaker"
	app.EnableBashCompletion = true

	var executionId string
	app.Commands = []cli.Command{
		{
			Name: "exec",
			Aliases: []string{"e"},
			Action: func(c * cli.Context) {
				executionId = c.Args().First()
				println("boom", executionId)
			},

		},
	}

	app.Run(os.Args)
}

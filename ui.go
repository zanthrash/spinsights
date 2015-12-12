package spinsights

import (
	ui "github.com/gizak/termui"
	"time"
	"os/exec"
	"fmt"
	"github.com/gizak/termui/debug"
)

var (
	client = DefalutClient
	showSucceeded = true
	orcaDetails *InstanceDetail
	actions = map[string]string{"q":"Quit", "f":"Toggle Success"}
	instructions = ui.NewPar("")
	info = ui.NewList()
	stages = ui.NewList()
	exception = ui.NewPar("None")
)

func drawInstructions() {
	instructions.Height = 3
	instructions.BorderLabel = "Spinisghts"

	if orcaDetails != nil {
		actions["l"] = "Orca Logs"
	}

	if instructions.Text == "" {
		for key, action := range actions {
			instructions.Text += fmt.Sprintf(" %s: %s ", key, action)
		}
	}

}

func drawInfo(exe *Execution) {
	info.BorderLabel = "Info"
	info.ItemFgColor = ui.ColorYellow
	info.Items = []string{
		"Name: " + exe.Name,
		fmt.Sprintf("Status: [%s]%s",  exe.Status, getStatusColor(exe.Status)),
	}
	info.Height = len(info.Items) + 2
}

func drawStages(exe *Execution) {
	stages.BorderLabel = "Stages"
	stages.ItemFgColor = ui.ColorYellow
	stageList := make([]string, 0)

	exception.Border = false
	exception.BorderLabel = "Exception"
	exception.Height = 3
	exception.Text = ""

	for i := range exe.Stages {
		stage := exe.Stages[i]

		statusColor := getStatusColor(stage.Status)
		stageInfo := fmt.Sprintf("%s [%s]%s", stage.Name, stage.Status, statusColor )
		stageList = append(stageList, stageInfo)

		if &stage.Context.Exception != nil {
			exception.Text = fmt.Sprintf("%s \n%s", stage.Context.Exception.Details.Error, stage.Context.Exception.Details.StackTrace)
			exception.Height = 20
		}

		for t := range exe.Stages[i].Tasks {
			task := exe.Stages[i].Tasks[t]
			if task.Status != "SUCCEEDED" || stage.Status == "SUCCEEDED" && showSucceeded {
				statusColor = getStatusColor(task.Status)
				taskInfo := fmt.Sprintf("  %s [%s]%s", task.Name, task.Status, statusColor)
				stageList = append(stageList, taskInfo)
			}
		}
	}
	stages.Items = stageList
	stages.Height = len(stageList) + 2
}

func RenderPipeline(executionId string) {
	if err := ui.Init(); err != nil {
		panic(err)
	}

	defer ui.Close()

	ui.Merge("timer", ui.NewTimerCh(time.Second * 5) )

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, instructions),
		),
		ui.NewRow(
			ui.NewCol(5, 0, info),
		),
		ui.NewRow(
			ui.NewCol(5, 0, stages),
			ui.NewCol(7, 0, exception),
		),
	)

	ui.Body.Align()

	draw := func(exe *Execution ) {
		if exe != nil {
			drawInstructions()
			drawInfo(exe)
			drawStages(exe)
			ui.Body.Align()
		}

		ui.Render(ui.Body)
	}

	fetchAndDraw := func() {
		execution, _ := client.GetExecutionById(executionId)

		if orcaDetails == nil {
			debug.Log("fetching orca details\n")
			searchResults, _ := client.InstanceSearch(execution.ExecutingInstance)
			if len(searchResults) > 0  && len(searchResults[0].Results) > 0{
				result := searchResults[0].Results[0]
				orcaDetails, _ = client.GetInstanceDetails(result)
			}
		}
		draw(execution)
	}

	fetchAndDraw()

	ui.Handle("/timer/5s", func(e ui.Event) {
		fetchAndDraw()
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/kbd/f", func(e ui.Event) {
		showSucceeded = !showSucceeded
		fetchAndDraw()
	})

	ui.Handle("/sys/kbd/q", func(e ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/l", func(e ui.Event) {
		if orcaDetails != nil {
			tomcatLogUrl := fmt.Sprintf("http://%s:7001/AdminLogs/list?view=tomcat/catalina.out", orcaDetails.PrivateIpAddress)
			debug.Log(orcaDetails.PrivateIpAddress)
			cmd := exec.Command("open", tomcatLogUrl)
			go cmd.Start()
		}
	})

	ui.Loop()

}

func getStatusColor(status string) string {
	switch status {
	case "SUCCEEDED": return "(fg-green)"
	case "TERMINAL": return "(fg-red)"
	case "NOT_STARTED": return "(fg-cyan)"
	default: return "(fg-cyan)"
	}
}


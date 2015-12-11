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
)

func RenderPipeline(executionId string) {
	if err := ui.Init(); err != nil {
		panic(err)
	}

	var orcaDetails *InstanceDetail

	defer ui.Close()
	ui.Merge("timer", ui.NewTimerCh(time.Second * 5) )

	instructions := ui.NewPar("Loading...")
	instructions.Height = 3
	instructions.BorderLabel = "Spinisghts"

	info := ui.NewList()
	info.Items = []string{"Loading.."}
	info.BorderLabel = "Info"
	info.ItemFgColor = ui.ColorYellow
	info.Height = 3

	stages := ui.NewList()
	stages.Items = []string{"Loading.."}
	stages.BorderLabel = "Stages"
	stages.ItemFgColor = ui.ColorYellow
	stages.Height = 3

	exception := ui.NewPar("None")
	exception.Border = false
	exception.BorderLabel = "Exception"
	exception.Height = 3

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

			if &orcaDetails != nil {
				instructions.Text = "l => Orca Logs"
			}

			info.Items = []string{
				"Name: " + exe.Name,
				fmt.Sprintf("Status: [%s]%s",  exe.Status, getStatusColor(exe.Status)),
			}
			info.Height = len(info.Items) + 2

			stageList := make([]string, 0)
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
					statusColor = getStatusColor(task.Status)
					taskInfo := fmt.Sprintf("  %s [%s]%s", task.Name, task.Status, statusColor)
					stageList = append(stageList, taskInfo)
				}
			}
			stages.Items = stageList
			stages.Height = len(stageList) + 2


			exception.Height = 20

			ui.Body.Align()

		}
		ui.Render(ui.Body)
	}

	fetchAndDraw := func() {
		execution, _ := client.GetExecutionById(executionId)

		if orcaDetails == nil {
			searchResults, _ := client.InstanceSearch(execution.ExecutingInstance)
			result := (*searchResults)[0].Results[0]
			orcaDetails, _ = client.GetInstanceDetails(&result)
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

	ui.Handle("/sys/kbd/q", func(e ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/l", func(e ui.Event) {
		tomcatLogUrl := fmt.Sprintf("http://%s:7001/AdminLogs/list?view=tomcat/catalina.out", orcaDetails.PrivateIpAddress)
		debug.Log(orcaDetails.PrivateIpAddress)
		cmd := exec.Command("open", tomcatLogUrl)
		go cmd.Start()
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


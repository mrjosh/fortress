package table

import (
	"os"

	"github.com/mrjosh/fortress/pkg/pipelines/sdk"
	"github.com/olekukonko/tablewriter"
)

func PrintPipelines(rows []*sdk.Pipeline) error {

	data := [][]string{}
	for i, row := range rows {
		data = append(data, row.GetData()...)
		if i != len(rows)-1 {
			data = append(data, []string{"", "", "", "", "", "", "", "", ""})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"PIPELINE", "JOB", "STATUS", "STAGE", "BRANCH", "COMMIT", "DURATION", "TRIGGERER", "CREATED_AT"})
	table.SetBorder(false)
	table.AppendBulk(data)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
		tablewriter.Colors{},
	)

	// Change table lines
	table.SetAutoWrapText(false)

	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetHeaderLine(false)

	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAutoMergeCellsByColumnIndex([]int{0, 4, 5, 7, 8})

	table.Render()
	return nil
}

package helpers

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/pterm/pterm"
)

type PipelineRow struct {
	Index int
	*PipelineDetails
}

type PipelineStatus int

func (s PipelineStatus) String() string {
	switch s {
	case PipelineStatusNeutral:
		return "neutral"
	case PipelineStatusSkipped:
		return "skipped"
	case PipelineStatusActionRequired:
		return "action_required"
	case PipelineStatusManual:
		return "manual"
	case PipelineStatusTimedOut:
		return "timed_out"
	case PipelineStatusStale:
		return "stale"
	case PipelineStatusInProgress:
		return "in_progress"
	case PipelineStatusRunning:
		return "running"
	case PipelineStatusQueued:
		return "queued"
	case PipelineStatusPending:
		return "pending"
	case PipelineStatusWaiting:
		return "waiting"
	case PipelineStatusSuccess:
		return "success"
	case PipelineStatusCompleted:
		return "completed"
	case PipelineStatusFailed:
		return "failed"
	case PipelineStatusFailure:
		return "failure"
	case PipelineStatusCancelled:
		return "cancelled"
	case PipelineStatusCanceled:
		return "canceled"
	}
	return "unknown"
}

const (
	PipelineStatusUnknown PipelineStatus = iota
	PipelineStatusCreated
	PipelineStatusSuccess
	PipelineStatusCompleted
	PipelineStatusFailed
	PipelineStatusFailure
	PipelineStatusCancelled
	PipelineStatusCanceled
	PipelineStatusInProgress
	PipelineStatusRunning
	PipelineStatusQueued
	PipelineStatusPending
	PipelineStatusWaiting
	PipelineStatusNeutral
	PipelineStatusActionRequired
	PipelineStatusManual
	PipelineStatusTimedOut
	PipelineStatusSkipped
	PipelineStatusStale
)

type PipelineDetails struct {
	ID       int
	CommitID string
	User     string
	Duration int
	Branch   string
	Status   PipelineStatus
}

func (p *PipelineRow) initTable() {

	pterm.Description.Prefix = pterm.Prefix{
		Text:  "PENDING",
		Style: pterm.NewStyle(pterm.FgYellow),
	}

	pterm.Debug.Prefix = pterm.Prefix{
		Text:  "Unknown",
		Style: pterm.NewStyle(pterm.FgDarkGray),
	}

	pterm.Info.Prefix = pterm.Prefix{
		Text:  "RUNNING",
		Style: pterm.NewStyle(pterm.FgBlue),
	}

	pterm.Warning.Prefix = pterm.Prefix{
		Text:  "CANCELLED",
		Style: pterm.NewStyle(pterm.FgGray),
	}

	pterm.Success.Prefix = pterm.Prefix{
		Text:  "PASSED",
		Style: pterm.NewStyle(pterm.FgGreen),
	}

	pterm.Error.Prefix = pterm.Prefix{
		Text:  "FAILED",
		Style: pterm.NewStyle(pterm.FgRed),
	}

}

func (p *PipelineRow) GetTableRowStrings() []string {

	p.initTable()

	var (
		pipStatus string
		commitID  = []rune(p.CommitID)
	)

	switch p.Status {
	case PipelineStatusUnknown:
		pipStatus = pterm.Debug.Sprint()
	case PipelineStatusInProgress, PipelineStatusRunning:
		pipStatus = pterm.Info.Sprint()
	case PipelineStatusCanceled, PipelineStatusCancelled:
		pipStatus = pterm.Warning.Sprint()
	case PipelineStatusSuccess, PipelineStatusCompleted:
		pipStatus = pterm.Success.Sprint()
	case PipelineStatusFailed:
		pipStatus = pterm.Error.Sprint()
	case PipelineStatusPending, PipelineStatusCreated, PipelineStatusQueued:
		pipStatus = pterm.Description.Sprint()
	}

	var (
		branch string
		ref    = []rune(p.Branch)
	)

	if len(ref) > 15 {
		branch = fmt.Sprintf("%s...", string(ref[0:15]))
	} else {
		branch = string(ref)
	}

	duration := fmt.Sprintf("%d", p.Duration)
	if p.Duration == 0 {
		duration = "-"
	}

	return []string{
		pipStatus,
		fmt.Sprintf("#%d", p.ID),
		branch,
		string(commitID[0:8]),
		duration,
		color.BlueString("@%s", p.User),
	}
}

func PrintTable(rows []*PipelineRow) error {

	defaultTable := pterm.TablePrinter{
		Style:          &pterm.ThemeDefault.TableStyle,
		HeaderStyle:    &pterm.ThemeDefault.TableHeaderStyle,
		Separator:      " | ",
		SeparatorStyle: &pterm.ThemeDefault.TableSeparatorStyle,
	}

	tableData := pterm.TableData{
		{"STATUS", "PIPELINE", "BRANCH", "COMMIT", "DURATION", "TRIGGERER"},
	}

	for _, row := range rows {
		tableData = append(tableData, row.GetTableRowStrings())
	}

	table := defaultTable.
		WithHasHeader(true).
		WithBoxed(true).
		WithData(tableData)

	return table.Render()
}

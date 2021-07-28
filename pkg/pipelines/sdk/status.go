package sdk

import "github.com/fatih/color"

type PipelineStatus int

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

func GetGithubPipelineStatus(status string) PipelineStatus {
	switch status {
	case PipelineStatusNeutral.String():
		return PipelineStatusNeutral
	case PipelineStatusSkipped.String():
		return PipelineStatusSkipped
	case PipelineStatusActionRequired.String():
		return PipelineStatusActionRequired
	case PipelineStatusTimedOut.String():
		return PipelineStatusTimedOut
	case PipelineStatusStale.String():
		return PipelineStatusStale
	case PipelineStatusInProgress.String():
		return PipelineStatusInProgress
	case PipelineStatusQueued.String():
		return PipelineStatusQueued
	case PipelineStatusWaiting.String():
		return PipelineStatusWaiting
	case PipelineStatusSuccess.String(), PipelineStatusCompleted.String():
		return PipelineStatusSuccess
	case PipelineStatusFailure.String():
		return PipelineStatusFailed
	case PipelineStatusCancelled.String():
		return PipelineStatusCancelled
	}
	return PipelineStatusUnknown
}

func GetPipelineRowStatus(status string) PipelineStatus {
	switch status {
	case PipelineStatusSkipped.String():
		return PipelineStatusSkipped
	case PipelineStatusCreated.String():
		return PipelineStatusCreated
	case PipelineStatusManual.String():
		return PipelineStatusManual
	case PipelineStatusRunning.String():
		return PipelineStatusRunning
	case PipelineStatusPending.String():
		return PipelineStatusPending
	case PipelineStatusSuccess.String():
		return PipelineStatusSuccess
	case PipelineStatusFailed.String():
		return PipelineStatusFailed
	case PipelineStatusCanceled.String():
		return PipelineStatusCanceled
	}
	return PipelineStatusUnknown
}

func GetColumnStatus(status PipelineStatus) string {
	switch status {
	case PipelineStatusInProgress, PipelineStatusRunning:
		return color.BlueString("RUNNING")
	case PipelineStatusCanceled, PipelineStatusCancelled:
		return color.YellowString("CANCELED")
	case PipelineStatusSuccess, PipelineStatusCompleted:
		return color.GreenString("PASSED")
		//if p.HasWarnings {
		//return color.YellowString("PASSED")
		//}
	case PipelineStatusManual, PipelineStatusActionRequired:
		return color.HiYellowString("MANUAL")
	case PipelineStatusSkipped:
		return color.HiYellowString("SKIPPED")
	case PipelineStatusFailed:
		return color.RedString("FAILED")
	case PipelineStatusPending, PipelineStatusCreated, PipelineStatusQueued:
		return color.YellowString("PENDING")
	}
	return color.BlackString("UNKNOWN")
}

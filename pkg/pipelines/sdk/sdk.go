package sdk

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/justincampbell/timeago"
)

type Pipeline struct {
	Index int
	*PipelineDetails
}

type Stage struct {
	ID       int
	Status   PipelineStatus
	Duration string
	Name     string
}

type PipelineDetails struct {
	ID          int
	CommitID    string
	User        string
	Duration    string
	Branch      string
	Status      PipelineStatus
	Stages      []Stage
	HasWarnings bool
	CreatedAt   time.Time
}

func (p *Pipeline) GetData() [][]string {

	var (
		branch   string
		ref      = []rune(p.Branch)
		commitID = []rune(p.CommitID)
	)

	if len(ref) > 15 {
		branch = fmt.Sprintf("%s...", string(ref[0:15]))
	} else {
		branch = string(ref)
	}

	sort.SliceStable(p.Stages, func(i, j int) bool {
		return true
	})

	rows := [][]string{}
	for _, stage := range p.Stages {
		duration := stage.Duration
		if stage.Duration == "0" {
			duration = "-"
		}
		rows = append(rows, []string{
			fmt.Sprintf("#%d", p.ID),
			fmt.Sprintf("#%d", stage.ID),
			GetColumnStatus(stage.Status),
			stage.Name,
			branch,
			string(commitID[0:8]),
			duration,
			fmt.Sprintf("@%s", p.User),
			timeago.FromTime(p.CreatedAt),
		})
	}

	return rows
}

func HumanizeDuration(duration int64) string {
	timeStamp := time.Unix(duration, 0).UTC()
	hours, minutes, seconds := timeStamp.Clock()
	hoursString, minutesString, secondsString :=
		strconv.Itoa(hours),
		strconv.Itoa(minutes),
		strconv.Itoa(seconds)
	if len(hoursString) == 1 {
		hoursString = fmt.Sprintf("0%s", hoursString)
	}
	if len(minutesString) == 1 {
		minutesString = fmt.Sprintf("0%s", minutesString)
	}
	if len(secondsString) == 1 {
		secondsString = fmt.Sprintf("0%s", secondsString)
	}
	return fmt.Sprintf("%s:%s:%s", hoursString, minutesString, secondsString)
}

package gitlab

import "time"

type PipelineRow struct {
	Index int
	*Pipeline
}

type PipelineStatus int

type Pipeline struct {
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

func (p *Pipeline) GetID() int {
	return p.ID
}

func (p *Pipeline) GetCommitID() string {
	return p.CommitID
}

func (p *Pipeline) GetUser() string {
	return p.User
}

func (p *Pipeline) GetDuration() string {
	return p.Duration
}

func (p *Pipeline) GetBranch() string {
	return p.Branch
}

func (p *Pipeline) GetStages() []Stage {
	return p.Stages
}

func (p *Pipeline) GetCreatedAt() time.Time {
	return p.CreatedAt
}

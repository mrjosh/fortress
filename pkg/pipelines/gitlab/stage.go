package gitlab

type Stage struct {
	ID       int
	Status   PipelineStatus
	Duration string
	Name     string
}

func (s *Stage) GetID() int {
	return s.ID
}

func (s *Stage) GetStatus() PipelineStatus {
	return s.Status
}

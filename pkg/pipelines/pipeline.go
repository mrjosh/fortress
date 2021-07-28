package pipelines

import (
	"github.com/mrjosh/fortress/cmds/flags"
	"github.com/mrjosh/fortress/pkg/git"
	"github.com/mrjosh/fortress/pkg/pipelines/github"
	"github.com/mrjosh/fortress/pkg/pipelines/gitlab"
	"github.com/mrjosh/fortress/pkg/pipelines/sdk"
)

//type PipelineRow interface {
//GetData() [][]string
//}

//type Pipeline interface {
//GetRows() ([]*PipelineRow, error)
//GetStages() ([]*Stage, error)
//}

func List(cfg *flags.JobsCommandFlags) ([]*sdk.Pipeline, error) {
	switch cfg.Service {
	case int(git.GitServiceGithub):
		return github.GetPipelines(cfg)
	default:
		return gitlab.GetPipelines(cfg)
	}
}

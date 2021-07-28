package github

import (
	"context"
	"strings"

	githubClient "github.com/google/go-github/v37/github"
	"github.com/mrjosh/fortress/cmds/flags"
	"github.com/mrjosh/fortress/pkg/pipelines/sdk"
)

func GetPipelines(cfg *flags.JobsCommandFlags) ([]*sdk.Pipeline, error) {

	var opt *githubClient.ListWorkflowRunsOptions
	if cfg.Branch != "" && cfg.Branch != flags.AllBranchesOption {
		opt = &githubClient.ListWorkflowRunsOptions{}
		opt.Branch = cfg.Branch
	}

	status := sdk.GetPipelineRowStatus(cfg.Status)
	if status != sdk.PipelineStatusUnknown {
		if opt == nil {
			opt = &githubClient.ListWorkflowRunsOptions{}
		}
		opt.Status = cfg.Status
	}

	repoSlice := strings.Split(cfg.Repository, "/")

	git := githubClient.NewClient(nil)
	workflows, _, err := git.Actions.ListRepositoryWorkflowRuns(context.Background(), repoSlice[0], repoSlice[1], opt)
	if err != nil {
		return nil, err
	}

	rows := make([]*sdk.Pipeline, 0)

	for index, pipeline := range workflows.WorkflowRuns {

		status := pipeline.GetConclusion()
		if status == "" {
			status = pipeline.GetStatus()
		}

		rows = append(rows, &sdk.Pipeline{
			Index: index,
			PipelineDetails: &sdk.PipelineDetails{
				ID:       int(pipeline.GetID()),
				Branch:   pipeline.GetHeadBranch(),
				User:     pipeline.Repository.Owner.GetLogin(),
				CommitID: pipeline.HeadCommit.GetID(),
				Status:   sdk.GetGithubPipelineStatus(status),
				Duration: "0",
			},
		})
	}

	return rows, nil
}

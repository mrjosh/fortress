package gitlab

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/mrjosh/fortress/cmds/flags"
	"github.com/mrjosh/fortress/pkg/pipelines/sdk"
	gitlabClient "github.com/xanzy/go-gitlab"
)

func GetPipelines(cfg *flags.JobsCommandFlags) ([]*sdk.Pipeline, error) {

	gitBaseURL := fmt.Sprintf("https://%s/api/v4", cfg.BaseURL)
	git, err := gitlabClient.NewClient(cfg.AccessToken, gitlabClient.WithBaseURL(gitBaseURL))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	opt := &gitlabClient.ListProjectPipelinesOptions{}
	opt.PerPage = cfg.PerPage
	if cfg.Latest {
		opt.PerPage = 1
	}

	if cfg.Branch != "" && cfg.Branch != flags.AllBranchesOption {
		opt.Ref = gitlabClient.String(cfg.Branch)
	}

	status := sdk.GetPipelineRowStatus(cfg.Status)
	if status != sdk.PipelineStatusUnknown {
		opt.Status = gitlabClient.BuildState(gitlabClient.BuildStateValue(status.String()))
	}

	pipelinesData, _, err := git.Pipelines.ListProjectPipelines(cfg.Repository, opt)
	if err != nil {
		return nil, err
	}

	var (
		wg   = sync.WaitGroup{}
		rows = make([]*sdk.Pipeline, 0)
	)

	for index, pip := range pipelinesData {

		wg.Add(1)
		go func(index int, pip *gitlabClient.PipelineInfo) {

			jobs, _, err := git.Jobs.ListPipelineJobs(cfg.Repository, pip.ID, nil)
			if err != nil {
				log.Fatal(err)
				return
			}

			if len(jobs) > 0 {
				var (
					hasWarnings = false
					//duration    int64
					stages = make([]sdk.Stage, 0)
				)
				for _, job := range jobs {
					if jobs[0].Pipeline.Status == sdk.PipelineStatusSuccess.String() {
						if job.Status == sdk.PipelineStatusCanceled.String() {
							hasWarnings = true
						}
					}
					//duration += int64(job.Duration)
					stages = append(stages, sdk.Stage{
						ID:       job.ID,
						Status:   sdk.GetPipelineRowStatus(job.Status),
						Name:     job.Name,
						Duration: sdk.HumanizeDuration(int64(job.Duration)),
					})
				}
				rows = append(rows, &sdk.Pipeline{
					Index: index,
					PipelineDetails: &sdk.PipelineDetails{
						ID:          pip.ID,
						Branch:      pip.Ref,
						User:        jobs[0].User.Username,
						CommitID:    pip.SHA,
						Status:      sdk.GetPipelineRowStatus(pip.Status),
						Stages:      stages,
						HasWarnings: hasWarnings,
						CreatedAt:   *pip.CreatedAt,
					},
				})
			}

			wg.Done()

		}(index, pip)
	}

	wg.Wait()

	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].Index < rows[j].Index
	})

	return rows, nil
}

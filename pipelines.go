package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	github "github.com/google/go-github/v37/github"
	gitlab "github.com/xanzy/go-gitlab"
)

func getGithubPipelineStatus(status string) PipelineStatus {
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

func getGitlabPipelineStatus(status string) PipelineStatus {
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

func getGitlabPipelines(cfg *GitConfig) ([]*PipelineRow, error) {

	gitBaseURL := fmt.Sprintf("https://%s/api/v4", cfg.BaseURL)
	git, err := gitlab.NewClient(cfg.AccessToken, gitlab.WithBaseURL(gitBaseURL))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	opt := &gitlab.ListProjectPipelinesOptions{}
	if cfg.Branch != "" && cfg.Branch != AllBranchedOption {
		opt.Ref = gitlab.String(cfg.Branch)
	}

	status := getGitlabPipelineStatus(cfg.Status)
	if status != PipelineStatusUnknown {
		opt.Status = gitlab.BuildState(gitlab.BuildStateValue(status.String()))
	}

	pipelinesData, _, err := git.Pipelines.ListProjectPipelines(cfg.Repository, opt)
	if err != nil {
		return nil, err
	}

	var (
		wg   = sync.WaitGroup{}
		rows = make([]*PipelineRow, 0)
	)

	for index, pip := range pipelinesData {
		wg.Add(1)
		go func(index, pipID int) {

			pipeline, _, err := git.Pipelines.GetPipeline(cfg.Repository, pipID)
			if err != nil {
				log.Fatal(err)
				return
			}

			rows = append(rows, &PipelineRow{
				Index: index,
				PipelineDetails: &PipelineDetails{
					ID:       pipeline.ID,
					Branch:   pipeline.Ref,
					User:     pipeline.User.Username,
					CommitID: pipeline.SHA,
					Status:   getGitlabPipelineStatus(pipeline.Status),
					Duration: pipeline.Duration,
				},
			})

			wg.Done()

		}(index, pip.ID)
	}

	wg.Wait()

	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].Index < rows[j].Index
	})

	return rows, nil
}

func getGithubPipelines(cfg *GitConfig) ([]*PipelineRow, error) {

	var opt *github.ListWorkflowRunsOptions
	if cfg.Branch != "" && cfg.Branch != AllBranchedOption {
		opt = &github.ListWorkflowRunsOptions{}
		opt.Branch = cfg.Branch
	}

	status := getGitlabPipelineStatus(cfg.Status)
	if status != PipelineStatusUnknown {
		if opt == nil {
			opt = &github.ListWorkflowRunsOptions{}
		}
		opt.Status = cfg.Status
	}

	repoSlice := strings.Split(cfg.Repository, "/")

	git := github.NewClient(nil)
	workflows, _, err := git.Actions.ListRepositoryWorkflowRuns(context.Background(), repoSlice[0], repoSlice[1], opt)
	if err != nil {
		return nil, err
	}

	rows := make([]*PipelineRow, 0)

	for index, pipeline := range workflows.WorkflowRuns {

		status := pipeline.GetConclusion()
		if status == "" {
			status = pipeline.GetStatus()
		}

		rows = append(rows, &PipelineRow{
			Index: index,
			PipelineDetails: &PipelineDetails{
				ID:       int(pipeline.GetID()),
				Branch:   pipeline.GetHeadBranch(),
				User:     pipeline.Repository.Owner.GetLogin(),
				CommitID: pipeline.HeadCommit.GetID(),
				Status:   getGithubPipelineStatus(status),
				Duration: 0,
			},
		})
	}

	return rows, nil
}

func GetPipelines(cfg *GitConfig) ([]*PipelineRow, error) {
	if cfg.Service == GitConfigServiceGithub {
		return getGithubPipelines(cfg)
	}
	return getGitlabPipelines(cfg)
}

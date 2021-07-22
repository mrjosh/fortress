package helpers

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

type GitConfigService int

const (
	GitConfigServiceGithub GitConfigService = iota
	GitConfigServiceGitlab
)

type GitConfig struct {
	BaseURL     string
	Branch      string
	Status      string
	AccessToken string
	Repository  string
	Service     GitConfigService
}

func GetLocalGitConfig() (gitConfig *GitConfig, err error) {

	gitConfig = new(GitConfig)

	gitURLCmd := exec.Command("git", "remote", "get-url", "origin")
	gitURLName, err := gitURLCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git remote origin: %v", err)
	}
	gitBaseURL := strings.TrimSpace(string(gitURLName))

	gitBaseURLParsed := strings.Split(gitBaseURL, "@")
	if len(gitBaseURLParsed) > 1 {

		gitDomain := strings.Split(gitBaseURLParsed[1], ":")
		gitConfig.Repository = strings.ReplaceAll(gitDomain[1], ".git", "")
		gitConfig.BaseURL = gitDomain[0]

	} else {

		gitParsedURL, err := url.Parse(gitBaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to get git remote origin: %v", err)
		}

		gitConfig.Repository = strings.ReplaceAll(gitParsedURL.Path[1:], ".git", "")
		gitConfig.BaseURL = gitParsedURL.Host

	}

	gitBranchCmd := exec.Command("git", "branch", "--show-current")
	gitBranch, err := gitBranchCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git branch: %v", err)
	}
	gitConfig.Branch = strings.TrimSpace(string(gitBranch))

	gitConfig.Service = GitConfigServiceGitlab
	if gitConfig.BaseURL == "github.com" {
		gitConfig.Service = GitConfigServiceGithub
	}

	return
}

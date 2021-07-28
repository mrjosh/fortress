package git

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/mrjosh/fortress/cmds/flags"
)

type Service int

const (
	GitServiceGithub Service = iota
	GitServiceGitlab
)

func IsGitlab(s int) bool {
	return Service(s) == GitServiceGitlab
}

func IsGithub(s int) bool {
	return Service(s) == GitServiceGithub
}

func getGitRemoteBranch(gitExecutable string) (string, error) {
	gitBranchCmd := exec.Command(gitExecutable, "branch", "--show-current")
	gitBranchOutput, err := gitBranchCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git branch: %v", err)
	}
	return strings.TrimSpace(string(gitBranchOutput)), nil
}

func getGitRemoteOrigin(gitExecutable string) (string, error) {
	gitURLCmd := exec.Command(gitExecutable, "remote", "get-url", "origin")
	gitURLName, err := gitURLCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git remote origin: %v", err)
	}
	return strings.TrimSpace(string(gitURLName)), nil
}

func AutoLoadGitConfig(cFlags *flags.JobsCommandFlags) error {

	gitExecutable, err := exec.LookPath(cFlags.GitExecutable)
	if err != nil {
		return fmt.Errorf("git command not found, make sure git is installed on your system: %v", err)
	}

	gitBaseURL := cFlags.BaseURL
	gitRepo := cFlags.Repository

	if cFlags.BaseURL == "" {
		gitBaseURL, err = getGitRemoteOrigin(gitExecutable)
		if err != nil {
			return fmt.Errorf("failed to get git remote origin: %v", err)
		}
	}

	gitBaseURLParsed := strings.Split(gitBaseURL, "@")
	if len(gitBaseURLParsed) > 1 {

		gitDomain := strings.Split(gitBaseURLParsed[1], ":")
		if gitRepo == "" {
			gitRepo = strings.ReplaceAll(gitDomain[1], ".git", "")
		}
		cFlags.BaseURL = gitDomain[0]

	} else {

		if !strings.Contains(gitBaseURL, "http") {
			gitBaseURL = fmt.Sprintf("https://%s", gitBaseURL)
		}

		gitParsedURL, err := url.Parse(gitBaseURL)
		if err != nil {
			return fmt.Errorf("failed to get git remote origin: %v", err)
		}

		if gitRepo == "" {
			gitRepo = strings.ReplaceAll(gitParsedURL.Path[1:], ".git", "")
		}
		cFlags.BaseURL = gitParsedURL.Host
	}

	cFlags.Repository = gitRepo

	if cFlags.Branch == "" {
		cFlags.Branch, err = getGitRemoteBranch(gitExecutable)
		if err != nil {
			return fmt.Errorf("failed to get git branch: %v", err)
		}
	}

	if cFlags.AllBranches {
		cFlags.Branch = flags.AllBranchesOption
	}

	cFlags.Service = int(GitServiceGitlab)
	if cFlags.BaseURL == "github.com" {
		cFlags.Service = int(GitServiceGithub)
	}

	return nil
}

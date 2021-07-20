package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/mrjosh/fortress/config"
	"github.com/xanzy/go-gitlab"
)

const (
	AllBranchedOption = "*all-branches*"
)

func main() {

	cfg, err := config.Loadfile(fmt.Sprintf("%s/.config/fortress.yaml", os.Getenv("HOME")))
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to get git remote origin url: %v", err))
	}

	gitConfig, err := GetLocalGitConfig()
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to get git remote origin url: %v", err))
	}

	branch := flag.String("branch", gitConfig.Branch, "Branch name")
	allBranches := flag.Bool("all-branches", false, "Get pipelines on all branches")
	accessToken := flag.String("access-token", "", "Git Service Personal Access Token")
	status := flag.String("status", "", "Filter pipelines by status")
	flag.Parse()

	cfgDomain, err := cfg.GetConfigByDomain(gitConfig.BaseURL)
	if err == nil {
		accessToken = gitlab.String(cfgDomain.AccessToken)
	}

	gitConfig.AccessToken = *accessToken
	gitConfig.Branch = *branch
	gitConfig.Status = *status

	if *allBranches {
		gitConfig.Branch = AllBranchedOption
	}

	projectWithBranchName := fmt.Sprintf(
		"%s %s",
		gitConfig.Repository,
		color.YellowString("%s", gitConfig.Branch),
	)

	s := spinner.New(spinner.CharSets[14], time.Millisecond*100)
	s.Start()
	s.Prefix = fmt.Sprintf(
		"[%s] Getpping repository pipelines [%s]",
		color.CyanString(gitConfig.BaseURL),
		projectWithBranchName,
	)

	pipelines, err := GetPipelines(gitConfig)
	if err != nil {
		log.Fatalln(err)
		return
	}

	s.FinalMSG = fmt.Sprintf(
		"\n[%s] %s [%s]\n",
		color.CyanString(gitConfig.BaseURL),
		color.GreenString(fmt.Sprintf("✔️ Found %d pipelines", len(pipelines))),
		projectWithBranchName,
	)
	s.Stop()

	if err := PrintTable(pipelines); err != nil {
		log.Fatalln(err)
		return
	}

}

package cmds

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/mrjosh/fortress/cmds/flags"
	"github.com/mrjosh/fortress/config"
	"github.com/mrjosh/fortress/pkg/git"
	"github.com/mrjosh/fortress/pkg/pipelines"
	"github.com/mrjosh/fortress/pkg/table"
	"github.com/spf13/cobra"
)

func newJobsCommand() *cobra.Command {
	cFlags := new(flags.JobsCommandFlags)
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "Get pipeline jobs of the current git repository",
		RunE: func(_ *cobra.Command, _ []string) error {

			log.SetFlags(log.Lshortfile)

			cfg, err := config.LoadFile(fmt.Sprintf("%s/.config/fortress.yaml", os.Getenv("HOME")))
			if err != nil {
				return fmt.Errorf("failed to load config file: %v", err)
			}

			if err := git.AutoLoadGitConfig(cFlags); err != nil {
				return fmt.Errorf("failed to get git config: %v", err)
			}

			dc, ok := cfg.GetGitDomainConfig(cFlags.BaseURL)
			if ok {
				cFlags.AccessToken = dc.AccessToken
			}

			projectWithBranchName := fmt.Sprintf(
				"%s %s",
				cFlags.Repository,
				color.YellowString("%s", cFlags.Branch),
			)

			s := spinner.New(spinner.CharSets[14], time.Millisecond*100)
			s.Start()
			s.Prefix = fmt.Sprintf(
				"[%s] Getting repository pipelines [%s]",
				color.CyanString(cFlags.BaseURL),
				projectWithBranchName,
			)

			pips, err := pipelines.List(cFlags)
			if err != nil {
				return fmt.Errorf("failed to get git pipelines: %v", err)
			}

			s.FinalMSG = fmt.Sprintf(
				"\n[%s] %s [%s]\n",
				color.CyanString(cFlags.BaseURL),
				color.GreenString(fmt.Sprintf("✔️ Found %d pipelines", len(pips))),
				projectWithBranchName,
			)
			s.Stop()

			return table.PrintPipelines(pips)
		},
	}
	cmd.Flags().BoolVarP(&cFlags.AllBranches, "all-branches", "A", false, "Get all pipelines on all branches")
	cmd.Flags().StringVarP(&cFlags.Branch, "branch", "b", "", "Get pipelines on a specific branch")
	defaultConfigFilePath := fmt.Sprintf("%s/.config/fortress.yaml", os.Getenv("HOME"))
	cmd.Flags().StringVarP(&cFlags.ConfigFile, "config-file", "c", defaultConfigFilePath, "Fortress config file path")
	cmd.Flags().StringVarP(&cFlags.Status, "status", "s", "", "Get pipelines with specific status")
	cmd.Flags().StringVarP(&cFlags.AccessToken, "token", "t", "", "Git service personal access token")
	cmd.Flags().BoolVarP(&cFlags.Latest, "latest", "l", false, "Latest pipeline")
	cmd.Flags().IntVarP(&cFlags.PerPage, "limit", "L", 5, "Pipeline display limit number")
	cmd.Flags().StringVarP(&cFlags.BaseURL, "base-url", "U", "", "Git repository BaseURL")
	cmd.Flags().StringVarP(&cFlags.Repository, "repo", "R", "", "Git repository name")
	cmd.Flags().BoolVarP(&cFlags.Wide, "wide", "W", false, "With job details")
	defaultGitExecutable, err := exec.LookPath("git")
	if err != nil {
		defaultGitExecutable = "git"
	}
	cmd.Flags().StringVarP(&cFlags.GitExecutable, "git-executable", "G", defaultGitExecutable, "Git executable")
	return cmd
}

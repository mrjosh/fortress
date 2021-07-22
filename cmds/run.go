package cmds

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/mrjosh/fortress/config"
	"github.com/mrjosh/fortress/helpers"
	"github.com/spf13/cobra"
)

type JobsFlags struct {
	Branch      string
	AllBranches bool
	ConfigFile  string
	Status      string
	AccessToken string
}

func newJobsCommand() *cobra.Command {
	cFlags := new(JobsFlags)
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "Get pipeline jobs of the current git repository",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := config.Loadfile(fmt.Sprintf("%s/.config/fortress.yaml", os.Getenv("HOME")))
			if err != nil {
				log.Fatalln(fmt.Errorf("failed to get git remote origin url: %v", err))
			}

			gitConfig, err := helpers.GetLocalGitConfig()
			if err != nil {
				log.Fatalln(fmt.Errorf("failed to get git remote origin url: %v", err))
			}

			cfgDomain, err := cfg.GetConfigByDomain(gitConfig.BaseURL)
			if err == nil {
				cFlags.AccessToken = cfgDomain.AccessToken
			}

			gitConfig.AccessToken = cFlags.AccessToken
			if cFlags.Branch != "" {
				gitConfig.Branch = cFlags.Branch
			}
			if cFlags.Status != "" {
				gitConfig.Status = cFlags.Status
			}

			if cFlags.AllBranches {
				gitConfig.Branch = helpers.AllBranchesOption
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

			pipelines, err := helpers.GetPipelines(gitConfig)
			if err != nil {
				return err
			}

			s.FinalMSG = fmt.Sprintf(
				"\n[%s] %s [%s]\n",
				color.CyanString(gitConfig.BaseURL),
				color.GreenString(fmt.Sprintf("✔️ Found %d pipelines", len(pipelines))),
				projectWithBranchName,
			)
			s.Stop()

			return helpers.PrintTable(pipelines)
		},
	}
	cmd.Flags().BoolVarP(&cFlags.AllBranches, "all-branches", "A", false, "Get all pipelines on all branches")
	cmd.Flags().StringVarP(&cFlags.Branch, "branch", "b", "", "Get pipelines on a specific branch")
	defaultConfigFilePath := fmt.Sprintf("%s/.config/fortress.yaml", os.Getenv("HOME"))
	cmd.Flags().StringVarP(&cFlags.ConfigFile, "config-file", "c", defaultConfigFilePath, "Fortress config file path")
	cmd.Flags().StringVarP(&cFlags.Status, "status", "s", "", "Get pipelines with specific status")
	cmd.Flags().StringVarP(&cFlags.AccessToken, "token", "t", "", "Git service personal access token")
	return cmd
}

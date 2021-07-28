package cmds

import (
	"github.com/spf13/cobra"
)

func newLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get pipeline logs of the current git repository",
		RunE: func(_ *cobra.Command, _ []string) error {

			return nil
		},
	}
	return cmd
}

//func getGithubLogs() {

//ts := oauth2.StaticTokenSource(
//&oauth2.Token{AccessToken: "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"},
//)
//tc := oauth2.NewClient(cmd.Context(), ts)

//git := github.NewClient(tc)
//logsURL, _, err := git.Actions.GetWorkflowRunLogs(cmd.Context(), "mrjosh", "fortress", 1058166048, false)
//if err != nil {
//return err
//}

//log.Println(fmt.Sprintf("Downloading [%s]", logsURL))

//resp, err := http.Get(logsURL.String())
//if err != nil {
//return err
//}

//byt, err := ioutil.ReadAll(resp.Body)
//if err != nil {
//return err
//}

//num := int64(len(byt))
//zipReader, err := zip.NewReader(bytes.NewReader(byt), num)
//if err != nil {
//return err
//}
//defer resp.Body.Close()

//for _, logFile := range zipReader.File {
//log.Println(logFile.Name)
//}

//}

//func getGitlabLogs() {

//gitBaseURL := fmt.Sprintf("https://%s/api/v4", "git.tradingclub.app")
//git, err := gitlab.NewClient("XsUnsHb8hsBTkTPVEvc2", gitlab.WithBaseURL(gitBaseURL))
//if err != nil {
//log.Fatalf("Failed to create client: %v", err)
//}

//opts := &gitlab.ListJobsOptions{}

//logs, _, _ := git.Jobs.ListProjectJobs("tradingclub/application", opts)

////logs, _, err := git.Jobs.ListPipelineJobs("tradingclub/application", 570, nil)
////if err != nil {
////log.Fatalln(err)
////}

//for _, j := range logs {
//log.Println(j.Pipeline, "\n")
//}

//}

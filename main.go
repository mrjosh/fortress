package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/pterm/pterm"
	gitlab "github.com/xanzy/go-gitlab"
)

type PipelineRow struct {
	Index int
	p     *gitlab.Pipeline
}

func (p *PipelineRow) initTable() {

	pterm.Info.Prefix = pterm.Prefix{
		Text:  "CANCELED",
		Style: pterm.NewStyle(pterm.FgGray),
	}

	pterm.Success.Prefix = pterm.Prefix{
		Text:  "PASSED",
		Style: pterm.NewStyle(pterm.FgGreen),
	}

	pterm.Error.Prefix = pterm.Prefix{
		Text:  "FAILED",
		Style: pterm.NewStyle(pterm.FgRed),
	}

}

func (p *PipelineRow) GetTableRowStrings() []string {

	p.initTable()

	var (
		pipStatus string
		commitID  = []rune(p.p.SHA)
	)

	if p.p.Status == "canceled" {
		pipStatus = pterm.Info.Sprint()
	}

	if p.p.Status == "success" {
		pipStatus = pterm.Success.Sprint()
	}

	if p.p.Status == "failed" {
		pipStatus = pterm.Error.Sprint()
	}

	var (
		branch string
		ref    = []rune(p.p.Ref)
	)

	if len(ref) > 15 {
		branch = fmt.Sprintf("%s...", string(ref[0:15]))
	} else {
		branch = string(ref)
	}

	return []string{
		pipStatus,
		fmt.Sprintf("#%d", p.p.ID),
		branch,
		string(commitID[0:8]),
		fmt.Sprintf("%d", p.p.Duration),
		fmt.Sprintf("@%s", p.p.User.Username),
	}
}

func main() {

	projectID := flag.Int("project_id", 0, "Project ID")
	gitURL := flag.String("base_url", "gitlab.com", "Gitlab base url")
	accessToken := flag.String("access-token", "", "Gitlab personal access token")
	branch := flag.String("branch", "", "Branch name")
	status := flag.String("status", "", "Pipeline status")
	flag.Parse()

	gitBaseURL := fmt.Sprintf("https://%s/api/v4", *gitURL)

	if *projectID == 0 {
		log.Fatalln("ProjectID is required!")
		return
	}

	if *accessToken == "" {
		log.Fatalln("AccessToken is required!")
		return
	}

	spinner, err := pterm.DefaultSpinner.Start("Getting pipelines...")
	if err != nil {
		log.Fatalf("Failed to create progressbar: %v", err)
	}

	git, err := gitlab.NewClient(*accessToken, gitlab.WithBaseURL(gitBaseURL))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	opt := &gitlab.ListProjectPipelinesOptions{}
	if *branch != "" {
		opt.Ref = gitlab.String(*branch)
	}
	if *status != "" {
		opt.Status = gitlab.BuildState(gitlab.BuildStateValue(*status))
	}

	defaultTable := pterm.TablePrinter{
		Style:          &pterm.ThemeDefault.TableStyle,
		HeaderStyle:    &pterm.ThemeDefault.TableHeaderStyle,
		Separator:      " | ",
		SeparatorStyle: &pterm.ThemeDefault.TableSeparatorStyle,
	}

	pipelinesData, _, err := git.Pipelines.ListProjectPipelines(*projectID, opt)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	piplineRows := make([]*PipelineRow, 0)
	for index, pip := range pipelinesData {
		wg.Add(1)
		go func(index, pipID int) {

			pipeline, _, err := git.Pipelines.GetPipeline(*projectID, pipID)
			if err != nil {
				log.Fatal(err)
			}

			piplineRows = append(piplineRows, &PipelineRow{
				Index: index,
				p:     pipeline,
			})

			wg.Done()

		}(index, pip.ID)
	}

	wg.Wait()

	sort.SliceStable(piplineRows, func(i, j int) bool {
		return piplineRows[i].Index < piplineRows[j].Index
	})

	tableData := pterm.TableData{
		{"STATUS", "PIPELINE", "BRANCH", "COMMIT", "DURATION", "TRIGGERER"},
	}

	for _, pipRow := range piplineRows {
		tableData = append(tableData, pipRow.GetTableRowStrings())
	}

	table := defaultTable.
		WithHasHeader(true).
		WithBoxed(true).
		WithData(tableData)

	spinner.Success(fmt.Sprintf("Found %d pipelines", len(piplineRows)))

	if err := table.Render(); err != nil {
		log.Fatalln(err)
		return
	}

}

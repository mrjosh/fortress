package flags

const (
	AllBranchesOption = "*all-branches*"
)

type JobsCommandFlags struct {
	ConfigFile    string
	BaseURL       string
	Branch        string
	Status        string
	AccessToken   string
	Repository    string
	Service       int
	Latest        bool
	PerPage       int
	Wide          bool
	AllBranches   bool
	GitExecutable string
}

package cmds

import (
	"os"

	"github.com/spf13/cobra"
)

func Run() error {

	rootCmd := &cobra.Command{
		Use: "fortress",
		Long: `
   __           _                     
  / _|         | |                    
 | |_ ___  _ __| |_ _ __ ___  ___ ___ 
 |  _/ _ \| '__| __| '__/ _ \/ __/ __|
 | || (_) | |  | |_| | |  __/\__ \__ \
 |_| \___/|_|   \__|_|  \___||___/___/`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.SetArgs(os.Args[1:])
	rootCmd.AddCommand(newVersionCommand())
	rootCmd.AddCommand(newJobsCommand())
	rootCmd.AddCommand(newLogsCommand())

	return rootCmd.Execute()
}

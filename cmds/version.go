package cmds

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "Fortress v0.0.1\n")
			fmt.Fprintf(cmd.OutOrStdout(), "Build type: Release\n")
			fmt.Fprintf(cmd.OutOrStdout(), "Go version: %s\n", runtime.Version())
			return nil
		},
	}
}

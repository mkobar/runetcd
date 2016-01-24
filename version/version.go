package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{
		Use:   "version",
		Short: "version tells runetcd version.",
		Run:   CommandFunc,
	}

	Version = "v0.0.1"
)

func CommandFunc(cmd *cobra.Command, args []string) {
	fmt.Println(Version)
}

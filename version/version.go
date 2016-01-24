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

	Version = "v0.1.0"
)

func CommandFunc(cmd *cobra.Command, args []string) {
	fmt.Println(Version)
}

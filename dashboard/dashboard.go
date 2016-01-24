package dashboard

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{}
)

func CommandFunc(cmd *cobra.Command, args []string) {
	etcdBinary, err := cmd.Flags().GetString("etcd-binary")
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

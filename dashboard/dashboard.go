package dashboard

import "github.com/spf13/cobra"

const (
	webPort = ":8080"
)

var (
	Command = &cobra.Command{
		Use:   "dashboard",
		Short: "dashboard provides etcd dashboard in a web browser.",
		Run:   CommandFunc,
	}
)

func CommandFunc(cmd *cobra.Command, args []string) {

}

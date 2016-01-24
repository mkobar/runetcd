package kill

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gyuho/psn/ss"
	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{
		Use:   "kill",
		Short: "kill kills etcd.",
		Run:   CommandFunc,
	}
)

func CommandFunc(cmd *cobra.Command, args []string) {
	etcdBinary, err := cmd.Flags().GetString("etcd-binary")
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	color.Set(color.FgRed)
	fmt.Fprintf(os.Stdout, "\npsn is killing:\n\n")
	color.Unset()

	filter := &ss.Process{Program: etcdBinary}
	ssr, err := ss.List(filter, ss.TCP, ss.TCP6)
	if err != nil {
		fmt.Fprintln(os.Stdout, "exiting with:", err)
		return
	}

	ss.WriteToTable(os.Stdout, ssr...)
	fmt.Fprintf(os.Stdout, "\n")

	ss.Kill(os.Stdout, ssr...)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()
}

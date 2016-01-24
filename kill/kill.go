package kill

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gyuho/psn/ss"
	"github.com/spf13/cobra"
)

type Flag struct {
	EtcdBinary string
}

var (
	Command = &cobra.Command{
		Use:   "kill",
		Short: "kill kills etcd.",
		Run:   CommandFunc,
	}

	cmdFlag = Flag{}
)

func init() {
	cobra.EnablePrefixMatching = true
}

func init() {
	Command.PersistentFlags().StringVarP(&cmdFlag.EtcdBinary, "etcd-binary", "b", "bin/etcd", "Path of executatble etcd binary.")
}

func CommandFunc(cmd *cobra.Command, args []string) {
	color.Set(color.FgRed)
	fmt.Fprintf(os.Stdout, "\npsn is killing:\n\n")
	color.Unset()

	filter := &ss.Process{Program: cmdFlag.EtcdBinary}
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

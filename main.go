// runetcd runs etcd.
//
//	Usage:
//	  runetcd [command]
//
//	Available Commands:
//	  kill        kill kills etcd.
//	  dashboard   dashboard provides etcd dashboard in a web browser.
//	  demo        demo demos etcd in terminal.
//	  demo-web    demo-web demos etcd in a web browser.
//	  version     version tells runetcd version.
//
//	Flags:
//	  -h, --help[=false]: help for runetcd
//
//	Use "runetcd [command] --help" for more information about a command.
//
package main

import (
	"fmt"
	"os"

	"github.com/gophergala2016/runetcd/dashboard"
	"github.com/gophergala2016/runetcd/demo"
	"github.com/gophergala2016/runetcd/demoweb"
	"github.com/gophergala2016/runetcd/kill"
	"github.com/gophergala2016/runetcd/version"
	"github.com/spf13/cobra"
)

var (
	rootCommand = &cobra.Command{
		Use:        "runetcd",
		Short:      "runetcd runs etcd.",
		SuggestFor: []string{"runetfd", "rnetcd", "runetdc"},
	}
)

func init() {
	cobra.EnablePrefixMatching = true
}

func init() {
	rootCommand.AddCommand(kill.Command)
	rootCommand.AddCommand(dashboard.Command)
	rootCommand.AddCommand(demo.Command)
	rootCommand.AddCommand(demoweb.Command)
	rootCommand.AddCommand(version.Command)
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

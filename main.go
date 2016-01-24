// runetcd runs etcd.
package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gophergala2016/runetcd/kill"
	"github.com/gyuho/psn/ss"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetLevel(log.DebugLevel)
}

const (
	cliName        = "runetcd"
	cliDescription = "runetcd runs etcd."
)

type GlobalFlag struct {
	EtcdBinary          string
	IntervalPortRefresh time.Duration
	Timeout             time.Duration

	ClusterSize int

	ProcSave bool
	ProcPath string

	// demo command
	DemoSimple           bool
	DemoPause            time.Duration
	DemoConnectionNumber int
	DemoClientNumber     int
	DemoStressNumber     int

	DemoIsClientTLS    bool
	DemoIsPeerTLS      bool
	DemoCertPath       string
	DemoPrivateKeyPath string
	DemoCAPath         string

	// demo-web command
	DemoWebSimple           bool
	DemoWebPause            time.Duration
	DemoWebConnectionNumber int
	DemoWebClientNumber     int
	DemoWebStressNumber     int
}

var (
	// need a way to clean this
	globalFlag  = GlobalFlag{}
	globalPorts = ss.NewPorts()

	nameToTerminate = "etcd1"
	nameToStress    = "etcd2"

	stressKeyN = 5
	stressValN = 5

	certPath       = "testcerts/cert.pem" // signed key-pair
	privateKeyPath = "testcerts/key.pem"  // signed key-pair
	caPath         = "testcerts/ca.perm"  // CA certificate

	rootCommand = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"runetfd", "rnetcd", "runetdc"},
	}
	demoCommand = &cobra.Command{
		Use:   "demo",
		Short: "demo demos etcd in terminal.",
		Run:   demoCommandFunc,
	}
	demoWebCommand = &cobra.Command{
		Use:   "demo-web",
		Short: "demo-web demos etcd in a web browser.",
		Run:   demoWebCommandFunc,
	}
)

func init() {
	globalPorts.Refresh()
	go func() {
		for {
			select {
			case <-time.After(globalFlag.IntervalPortRefresh):
				globalPorts.Refresh()
			}
		}
	}()
}

func init() {
	rootCommand.PersistentFlags().StringVarP(&globalFlag.EtcdBinary, "etcd-binary", "b", "bin/etcd", "Path of executatble etcd binary.")
	rootCommand.PersistentFlags().DurationVar(&globalFlag.IntervalPortRefresh, "port-refresh", 10*time.Second, "Interval to refresh free ports.")
	rootCommand.PersistentFlags().DurationVar(&globalFlag.Timeout, "timeout", 5*time.Minute, "After timeout, etcd shuts down itself.")
	rootCommand.PersistentFlags().IntVarP(&globalFlag.ClusterSize, "cluster-size", "n", 3, "Size of cluster to create.")

	rootCommand.PersistentFlags().BoolVarP(&globalFlag.ProcSave, "proc-save", "s", false, "'true' to save the procfile to local disk.")
	rootCommand.PersistentFlags().StringVar(&globalFlag.ProcPath, "proc-path", "Procfile", "Path of Procfile to save.")

	demoCommand.PersistentFlags().BoolVar(&globalFlag.DemoSimple, "simple", false, "'true' to run demo without auto-termination.")
	demoCommand.PersistentFlags().DurationVar(&globalFlag.DemoPause, "pause", 5*time.Second, "Duration to pause between demo operations.")
	demoCommand.PersistentFlags().IntVar(&globalFlag.DemoConnectionNumber, "connection-number", 1, "Number of connections.")
	demoCommand.PersistentFlags().IntVar(&globalFlag.DemoClientNumber, "client-number", 10, "Number of clients.")
	demoCommand.PersistentFlags().IntVar(&globalFlag.DemoStressNumber, "stress-number", 10, "Size of stress requests.")
	demoCommand.PersistentFlags().BoolVar(&globalFlag.DemoIsClientTLS, "client-tls", false, "'true' to set up client-to-server TLS.")
	demoCommand.PersistentFlags().BoolVar(&globalFlag.DemoIsPeerTLS, "peer-tls", false, "'true' to set up peer-to-peer TLS.")
	demoCommand.PersistentFlags().StringVar(&globalFlag.DemoCertPath, "cert-path", certPath, "CERT path.")
	demoCommand.PersistentFlags().StringVar(&globalFlag.DemoPrivateKeyPath, "private-key-path", privateKeyPath, "Private key path.")
	demoCommand.PersistentFlags().StringVar(&globalFlag.DemoCAPath, "ca-path", caPath, "CA path.")

	demoWebCommand.PersistentFlags().BoolVar(&globalFlag.DemoWebSimple, "simple", false, "'true' to run demo without auto-termination.")
	demoWebCommand.PersistentFlags().DurationVar(&globalFlag.DemoWebPause, "pause", 5*time.Second, "Duration to pause between demo operations.")
	demoWebCommand.PersistentFlags().IntVar(&globalFlag.DemoWebConnectionNumber, "connection-number", 1, "Number of connections.")
	demoWebCommand.PersistentFlags().IntVar(&globalFlag.DemoWebClientNumber, "client-number", 10, "Number of clients.")
	demoWebCommand.PersistentFlags().IntVar(&globalFlag.DemoWebStressNumber, "stress-number", 10, "Size of stress requests.")

	rootCommand.AddCommand(demoCommand)
	rootCommand.AddCommand(demoWebCommand)
	rootCommand.AddCommand(kill.Command)
}

func init() {
	cobra.EnablePrefixMatching = true
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

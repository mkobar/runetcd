// runetcd runs etcd.
package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
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
	killCommand = &cobra.Command{
		Use:   "kill",
		Short: "kill kills etcd.",
		Run:   killCommandFunc,
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

// runetcd runs etcd.
package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
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
}

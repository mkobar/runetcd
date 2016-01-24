package demo

import (
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gophergala2016/runetcd/etcdproc"
	"github.com/gyuho/psn/ss"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

type Flag struct {
	EtcdBinary          string
	IntervalPortRefresh time.Duration
	Timeout             time.Duration
	ClusterSize         int
	ProcSave            bool
	ProcPath            string
	IsSimple            bool
	Pause               time.Duration
	ConnectionNumber    int
	ClientNumber        int
	StressNumber        int
	IsClientTLS         bool
	IsPeerTLS           bool
	CertPath            string
	PrivateKeyPath      string
	CAPath              string
}

var (
	Command = &cobra.Command{
		Use:   "demo",
		Short: "demo demos etcd in terminal.",
		Run:   CommandFunc,
	}

	cmdFlag     = Flag{}
	globalPorts = ss.NewPorts()

	nameToTerminate = "etcd1"
	nameToStress    = "etcd2"

	stressKeyN = 5
	stressValN = 5

	certPath       = "testcerts/cert.pem" // signed key-pair
	privateKeyPath = "testcerts/key.pem"  // signed key-pair
	caPath         = "testcerts/ca.perm"  // CA certificate
)

func init() {
	cobra.EnablePrefixMatching = true
}

func init() {
	globalPorts.Refresh()
	go func() {
		for {
			select {
			case <-time.After(cmdFlag.IntervalPortRefresh):
				globalPorts.Refresh()
			}
		}
	}()
}

func init() {
	Command.PersistentFlags().StringVarP(&cmdFlag.EtcdBinary, "etcd-binary", "b", "bin/etcd", "Path of executatble etcd binary.")
	Command.PersistentFlags().DurationVar(&cmdFlag.IntervalPortRefresh, "port-refresh", 15*time.Second, "Interval to refresh free ports.")
	Command.PersistentFlags().DurationVar(&cmdFlag.Timeout, "timeout", 5*time.Minute, "After timeout, etcd shuts down itself.")
	Command.PersistentFlags().IntVarP(&cmdFlag.ClusterSize, "cluster-size", "n", 3, "Size of cluster to create.")
	Command.PersistentFlags().BoolVarP(&cmdFlag.ProcSave, "proc-save", "s", false, "'true' to save the procfile to local disk.")
	Command.PersistentFlags().StringVar(&cmdFlag.ProcPath, "proc-path", "Procfile", "Path of Procfile to save.")
	Command.PersistentFlags().BoolVar(&cmdFlag.IsSimple, "simple", false, "'true' to run demo without auto-termination.")
	Command.PersistentFlags().DurationVar(&cmdFlag.Pause, "pause", 5*time.Second, "Duration to pause between demo operations.")
	Command.PersistentFlags().IntVar(&cmdFlag.ConnectionNumber, "connection-number", 1, "Number of connections.")
	Command.PersistentFlags().IntVar(&cmdFlag.ClientNumber, "client-number", 10, "Number of clients.")
	Command.PersistentFlags().IntVar(&cmdFlag.StressNumber, "stress-number", 10, "Size of stress requests.")
	Command.PersistentFlags().BoolVar(&cmdFlag.IsClientTLS, "client-tls", false, "'true' to set up client-to-server TLS.")
	Command.PersistentFlags().BoolVar(&cmdFlag.IsPeerTLS, "peer-tls", false, "'true' to set up peer-to-peer TLS.")
	Command.PersistentFlags().StringVar(&cmdFlag.CertPath, "cert-path", certPath, "CERT path.")
	Command.PersistentFlags().StringVar(&cmdFlag.PrivateKeyPath, "private-key-path", privateKeyPath, "Private key path.")
	Command.PersistentFlags().StringVar(&cmdFlag.CAPath, "ca-path", caPath, "CA path.")
}

func CommandFunc(cmd *cobra.Command, args []string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stdout, "[demo.CommandFunc - panic]", err)
			os.Exit(0)
		}
	}()

	fs := make([]*etcdproc.Flags, cmdFlag.ClusterSize)
	for i := range fs {
		df, err := etcdproc.NewFlags(fmt.Sprintf("etcd%d", i+1), globalPorts, 11+i, "etcd-cluster-token", "new", uuid.NewV4().String(), cmdFlag.IsClientTLS, cmdFlag.IsPeerTLS, cmdFlag.CertPath, cmdFlag.PrivateKeyPath, cmdFlag.CAPath)
		if err != nil {
			fmt.Fprintln(os.Stdout, "exiting with:", err)
			return
		}
		fs[i] = df
	}

	c, err := etcdproc.CreateCluster(os.Stdout, nil, etcdproc.ToTerminal, cmdFlag.EtcdBinary, fs...)
	if err != nil {
		fmt.Fprintln(os.Stdout, "exiting with:", err)
		return
	}

	if cmdFlag.ProcSave {
		f, err := openToOverwrite(cmdFlag.ProcPath)
		if err != nil {
			fmt.Fprintln(os.Stdout, "exiting with:", err)
			return
		}
		c.WriteProc(f)
		f.Close()
	}

	// this does not run with the program exits with os.Exit(0)
	defer c.RemoveAllDataDirs()

	fmt.Fprintf(os.Stdout, "\n")
	fmt.Fprintln(os.Stdout, "####### Starting all of those 3 nodes in default cluster group")
	clusterDone := make(chan struct{})
	go func() {
		defer func() {
			clusterDone <- struct{}{}
		}()
		if err := c.StartAll(); err != nil {
			fmt.Fprintln(os.Stdout, "exiting with:", err)
			return
		}
	}()

	operationDone := make(chan struct{})
	if !cmdFlag.IsSimple {
		go func() {
			defer func() {
				operationDone <- struct{}{}
			}()

			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### Trying to terminate one of the node")
			if err := c.Terminate(nameToTerminate); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			// Stress here to trigger log compaction
			// (make terminated node fall behind)

			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### Trying to restart that node")
			if err := c.Restart(nameToTerminate); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### Stressing one node")
			if err := c.SimpleStress(os.Stdout, etcdproc.ToTerminal, nameToStress); err != nil {
				// if err := c.Stress(os.Stdout, nameToStress, cmdFlag.ConnectionNumber, cmdFlag.ClientNumber, cmdFlag.StressNumber, stressKeyN, stressValN); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### Watch and Put")
			if err := c.WatchAndPut(os.Stdout, nameToStress, cmdFlag.ConnectionNumber, cmdFlag.ClientNumber, cmdFlag.StressNumber); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			// TODO: not working for now
			if !cmdFlag.IsClientTLS {
				time.Sleep(cmdFlag.Pause)
				fmt.Fprintf(os.Stdout, "\n")
				fmt.Fprintln(os.Stdout, "####### Stats")
				if vm, err := c.GetStats(); err != nil {
					fmt.Fprintln(os.Stdout, "exiting with:", err)
					return
				} else {
					fmt.Fprintf(os.Stdout, "%+v\n", vm)
				}

				time.Sleep(cmdFlag.Pause)
				fmt.Fprintf(os.Stdout, "\n")
				fmt.Fprintln(os.Stdout, "####### Metrics")
				if vm, err := c.GetMetrics(); err != nil {
					fmt.Fprintln(os.Stdout, "exiting with:", err)
					return
				} else {
					for n, mm := range vm {
						var fb uint64
						if fv, ok := mm["etcd_storage_db_total_size_in_bytes"]; ok {
							fb = uint64(fv)
						}
						fmt.Fprintf(os.Stdout, "%s: etcd_storage_keys_total             = %f\n", n, mm["etcd_storage_keys_total"])
						fmt.Fprintf(os.Stdout, "%s: etcd_storage_db_total_size_in_bytes = %s\n", n, humanize.Bytes(fb))
						fmt.Fprintf(os.Stdout, "\n")
					}
				}
			}
		}()
	}

	select {
	case <-clusterDone:
		fmt.Fprintf(os.Stdout, "\n")
		fmt.Fprintln(os.Stdout, "[demo.CommandFunc END] etcd cluster terminated!")
		fmt.Fprintf(os.Stdout, "\n")
		return
	case <-operationDone:
		fmt.Fprintf(os.Stdout, "\n")
		fmt.Fprintln(os.Stdout, "[demo.CommandFunc END] operation terminated!")
		fmt.Fprintf(os.Stdout, "\n")
		return
	case <-time.After(cmdFlag.Timeout):
		fmt.Fprintf(os.Stdout, "\n")
		fmt.Fprintln(os.Stdout, "[demo.CommandFunc END] timed out!")
		fmt.Fprintf(os.Stdout, "\n")
		return
	}
}

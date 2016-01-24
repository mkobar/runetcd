package demo

import (
	"fmt"
	"os"
	"sort"
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
	IsSimpleSimulation  bool
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
	Command.PersistentFlags().BoolVar(&cmdFlag.IsSimpleSimulation, "simple-simulation", false, "'true' to run demo without auto-termination and some stressing. It overrides 'simple' flag.")
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

	if cmdFlag.IsSimpleSimulation {

		go func() {
			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### Stress")
			if err := c.Stress(cmdFlag.ConnectionNumber, cmdFlag.ClientNumber, cmdFlag.StressNumber, 15, 15); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### SimpleStress")
			if err := c.SimpleStress(); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}
		}()

	} else if !cmdFlag.IsSimple {

		go func() {
			defer func() {
				operationDone <- struct{}{}
			}()

			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### Terminate")
			if err := c.Terminate(nameToTerminate); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			// Stress here to trigger log compaction
			// (make terminated node fall behind)

			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### Restart")
			if err := c.Restart(nameToTerminate); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			key, val := []byte("sample_key"), []byte("sample_value")
			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### Put")
			if err := c.Put(key, val); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### Range")
			if err := c.Range(key); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### Stress")
			if err := c.Stress(cmdFlag.ConnectionNumber, cmdFlag.ClientNumber, cmdFlag.StressNumber, 15, 15); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### SimpleStress")
			if err := c.SimpleStress(); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			time.Sleep(cmdFlag.Pause)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintln(os.Stdout, "####### WatchAndPut")
			if err := c.WatchAndPut(cmdFlag.ConnectionNumber, cmdFlag.ClientNumber, cmdFlag.StressNumber); err != nil {
				fmt.Fprintln(os.Stdout, "exiting with:", err)
				return
			}

			if !cmdFlag.IsClientTLS { // TODO: not working for now
				time.Sleep(cmdFlag.Pause)
				fmt.Fprintf(os.Stdout, "\n")
				fmt.Fprintln(os.Stdout, "####### GetStats #1")
				vm, ne, err := c.GetStats()
				if err != nil {
					fmt.Fprintln(os.Stdout, "exiting with:", err)
					return
				} else {
					fmt.Fprintf(os.Stdout, "Endpoint To Stats: %+v\n", vm)
					fmt.Fprintf(os.Stdout, "Name To Endpoint : %+v\n", ne)

					fmt.Fprintf(os.Stdout, "\n")
					fmt.Fprintln(os.Stdout, "####### GetStats #2")
					endpoints := []string{}
					for _, endpoint := range ne {
						endpoints = append(endpoints, endpoint)
					}
					sort.Strings(endpoints)
					vm2, ne2, err := etcdproc.GetStats(endpoints...)
					if err != nil {
						fmt.Fprintln(os.Stdout, "exiting with:", err)
						return
					}
					fmt.Fprintf(os.Stdout, "Endpoint To Stats: %+v\n", vm2)
					fmt.Fprintf(os.Stdout, "Name To Endpoint : %+v\n", ne2)
				}

				time.Sleep(cmdFlag.Pause)
				fmt.Fprintf(os.Stdout, "\n")
				fmt.Fprintln(os.Stdout, "####### GetMetrics #1")
				{
					vm, ne, err := c.GetMetrics()
					if err != nil {
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
						}
						fmt.Fprintf(os.Stdout, "Name To Endpoint: %+v\n", ne)

						fmt.Fprintf(os.Stdout, "\n")
						fmt.Fprintln(os.Stdout, "####### GetMetrics #2")
						endpoints := []string{}
						for _, endpoint := range ne {
							endpoints = append(endpoints, endpoint)
						}
						sort.Strings(endpoints)
						vm2, ne2, err := etcdproc.GetMetrics(endpoints...)
						if err != nil {
							fmt.Fprintln(os.Stdout, "exiting with:", err)
							return
						}
						for n, mm := range vm2 {
							var fb uint64
							if fv, ok := mm["etcd_storage_db_total_size_in_bytes"]; ok {
								fb = uint64(fv)
							}
							fmt.Fprintf(os.Stdout, "%s: etcd_storage_keys_total             = %f\n", n, mm["etcd_storage_keys_total"])
							fmt.Fprintf(os.Stdout, "%s: etcd_storage_db_total_size_in_bytes = %s\n", n, humanize.Bytes(fb))
						}
						fmt.Fprintf(os.Stdout, "Name To Endpoint : %+v\n", ne2)
					}
				}

				fmt.Println()
			}
		}()

	}

	select {
	case <-clusterDone:
		fmt.Fprintln(os.Stdout, "[demo.CommandFunc END] etcd cluster terminated!")
		return
	case <-operationDone:
		fmt.Fprintln(os.Stdout, "[demo.CommandFunc END] operation terminated!")
		return
	case <-time.After(cmdFlag.Timeout):
		fmt.Fprintln(os.Stdout, "[demo.CommandFunc END] timed out!")
		return
	}
}

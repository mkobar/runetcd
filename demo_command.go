package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gophergala2016/runetcd/etcdproc"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

func demoCommandFunc(cmd *cobra.Command, args []string) {
	runDemoTerminal(os.Stdout, globalFlag.DemoSimple, globalFlag.DemoIsClientTLS, globalFlag.DemoIsPeerTLS, globalFlag.DemoCertPath, globalFlag.DemoPrivateKeyPath, globalFlag.DemoCAPath)
}

func runDemoTerminal(writer io.Writer, isSimple, isClientTLS, isPeerTLS bool, certPath, privateKeyPath, caPath string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(writer, "[main - panic]", err)
			os.Exit(0)
		}
	}()

	fs := make([]*etcdproc.Flags, globalFlag.ClusterSize)
	for i := range fs {
		df, err := etcdproc.NewFlags(fmt.Sprintf("etcd%d", i+1), globalPorts, 11+i, "etcd-cluster-token", "new", uuid.NewV4().String(), isClientTLS, isPeerTLS, certPath, privateKeyPath, caPath)
		if err != nil {
			fmt.Fprintln(writer, "exiting with:", err)
			return
		}
		fs[i] = df
	}

	c, err := etcdproc.CreateCluster(writer, nil, etcdproc.ToTerminal, globalFlag.EtcdBinary, fs...)
	if err != nil {
		fmt.Fprintln(writer, "exiting with:", err)
		return
	}

	if globalFlag.ProcSave {
		f, err := openToOverwrite(globalFlag.ProcPath)
		if err != nil {
			fmt.Fprintln(writer, "exiting with:", err)
			return
		}
		c.WriteProc(f)
		f.Close()
	}

	// this does not run with the program exits with os.Exit(0)
	defer c.RemoveAllDataDirs()

	fmt.Fprintf(writer, "\n")
	fmt.Fprintln(writer, "####### Starting all of those 3 nodes in default cluster group")
	clusterDone := make(chan struct{})
	go func() {
		defer func() {
			clusterDone <- struct{}{}
		}()
		if err := c.StartAll(); err != nil {
			fmt.Fprintln(writer, "exiting with:", err)
			return
		}
	}()

	operationDone := make(chan struct{})
	if !isSimple {
		go func() {
			defer func() {
				operationDone <- struct{}{}
			}()

			time.Sleep(globalFlag.DemoPause)
			fmt.Fprintf(writer, "\n")
			fmt.Fprintln(writer, "####### Trying to terminate one of the node")
			if err := c.Terminate(nameToTerminate); err != nil {
				fmt.Fprintln(writer, "exiting with:", err)
				return
			}

			// Stress here to trigger log compaction
			// (make terminated node fall behind)

			time.Sleep(globalFlag.DemoPause)
			fmt.Fprintf(writer, "\n")
			fmt.Fprintln(writer, "####### Trying to restart that node")
			if err := c.Restart(nameToTerminate); err != nil {
				fmt.Fprintln(writer, "exiting with:", err)
				return
			}

			time.Sleep(globalFlag.DemoPause)
			fmt.Fprintf(writer, "\n")
			fmt.Fprintln(writer, "####### Stressing one node")
			if err := c.SimpleStress(writer, etcdproc.ToTerminal, nameToStress); err != nil {
				// if err := c.Stress(writer, nameToStress, globalFlag.DemoConnectionNumber, globalFlag.DemoClientNumber, globalFlag.DemoStressNumber, stressKeyN, stressValN); err != nil {
				fmt.Fprintln(writer, "exiting with:", err)
				return
			}

			time.Sleep(globalFlag.DemoPause)
			fmt.Fprintf(writer, "\n")
			fmt.Fprintln(writer, "####### Watch and Put")
			if err := c.WatchAndPut(writer, nameToStress, globalFlag.DemoConnectionNumber, globalFlag.DemoClientNumber, globalFlag.DemoStressNumber); err != nil {
				fmt.Fprintln(writer, "exiting with:", err)
				return
			}

			// TODO: not working for now
			if !isClientTLS {
				time.Sleep(globalFlag.DemoPause)
				fmt.Fprintf(writer, "\n")
				fmt.Fprintln(writer, "####### Stats")
				if vm, err := c.GetStats(); err != nil {
					fmt.Fprintln(writer, "exiting writerith:", err)
					return
				} else {
					fmt.Fprintf(writer, "%+v\n", vm)
				}

				time.Sleep(globalFlag.DemoPause)
				fmt.Fprintf(writer, "\n")
				fmt.Fprintln(writer, "####### Metrics")
				if vm, err := c.GetMetrics(); err != nil {
					fmt.Fprintln(writer, "exiting with:", err)
					return
				} else {
					for n, mm := range vm {
						var fb uint64
						if fv, ok := mm["etcd_storage_db_total_size_in_bytes"]; ok {
							fb = uint64(fv)
						}
						fmt.Fprintf(writer, "%s: etcd_storage_keys_total             = %f\n", n, mm["etcd_storage_keys_total"])
						fmt.Fprintf(writer, "%s: etcd_storage_db_total_size_in_bytes = %s\n", n, humanize.Bytes(fb))
						fmt.Fprintf(writer, "\n")
					}
				}
			}
		}()
	}

	select {
	case <-clusterDone:
		fmt.Fprintf(writer, "\n")
		fmt.Fprintln(writer, "[runetcd demo END] etcd cluster terminated!")
		fmt.Fprintf(writer, "\n")
		return
	case <-operationDone:
		fmt.Fprintf(writer, "\n")
		fmt.Fprintln(writer, "[runetcd demo END] operation terminated!")
		fmt.Fprintf(writer, "\n")
		return
	case <-time.After(globalFlag.Timeout):
		fmt.Fprintf(writer, "\n")
		fmt.Fprintln(writer, "[runetcd demo END] timed out!")
		fmt.Fprintf(writer, "\n")
		return
	}
}

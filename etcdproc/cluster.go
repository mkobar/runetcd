package etcdproc

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
)

// Cluster groups a set of Node processes.
type Cluster struct {
	// define in pointer type to not copy over cluster
	mu                *sync.Mutex // guards the following
	maxProcNameLength *int
	wg                *sync.WaitGroup
	NameToNode        map[string]*Node
}

// CreateCluster creates a Cluster by parsing the input io.Reader.
// io.Writer and buufferStream channel are shared across cluster.
func CreateCluster(w io.Writer, bufStream chan string, outputOption OutputOption, command string, fs ...*Flags) (*Cluster, error) {
	if len(fs) == 0 {
		return nil, nil
	}

	var maxProcNameLength int
	c := &Cluster{
		mu:                &sync.Mutex{},
		maxProcNameLength: &maxProcNameLength,
		wg:                &sync.WaitGroup{},
		NameToNode:        make(map[string]*Node),
	}

	if err := CombineFlags(fs...); err != nil {
		return nil, err
	}

	colorIdx := 0
	for _, f := range fs {
		if colorIdx >= len(colors) {
			colorIdx = 0
		}

		name := f.Name
		if len(name) > *c.maxProcNameLength {
			*c.maxProcNameLength = len(name)
		}

		nd := &Node{
			pmu:                c.mu,
			pmaxProcNameLength: c.maxProcNameLength,

			outputOption: outputOption,
			colorIdx:     colorIdx,
			w:            w,
			BufferStream: bufStream,

			Command: command,
			Flags:   f,

			PID:        0,
			Terminated: false,
			cmd:        nil,
		}
		c.NameToNode[name] = nd

		colorIdx++
	}

	return c, nil
}

func (c *Cluster) WriteProc(w io.Writer) {
	lines := []string{}
	for name, nd := range c.NameToNode {
		line := strings.TrimSpace(fmt.Sprintf("%s: %s %s", name, nd.Command, nd.Flags))
		lines = append(lines, line)
	}
	sort.Strings(lines)
	fmt.Fprintf(w, "%s", strings.Join(lines, "\n"))
}

func (c *Cluster) DoneSafely(w io.Writer) {
	func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintln(w, err)
			}
		}()
		c.wg.Done()
	}()
}

func (c *Cluster) StartAll() error {
	if len(c.NameToNode) == 0 {
		return nil
	}

	// execute all Nodes at the same time
	c.wg.Add(len(c.NameToNode))
	for name, nd := range c.NameToNode {

		go func(name string, nd *Node) {

			defer c.DoneSafely(nd)

			cs := []string{"/bin/bash", "-c", nd.Command + " " + nd.Flags.String()}
			cmd := exec.Command(cs[0], cs[1:]...)
			cmd.Stdin = nil
			cmd.Stdout = nd
			cmd.Stderr = nd

			fmt.Fprintf(nd, "Starting %s\n", name)
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(nd, "Failed to start %s\n", name)
				return
			}
			nd.cmd = cmd
			nd.PID = cmd.Process.Pid

			cmd.Wait()

			fmt.Fprintf(nd, "Exiting %s\n", name)
			return

		}(name, nd)

	}

	sc := make(chan os.Signal, 10)
	go func() {
		c.wg.Wait()
		sc <- syscall.SIGINT
	}()
	signal.Notify(sc, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	<-sc

	c.TerminateAll()
	return nil
}

// TerminateAll terminates all Nodes in the group.
func (c *Cluster) TerminateAll() {
	defer func() {
		recover()
	}()

	if len(c.NameToNode) == 0 {
		return
	}

	for _, nd := range c.NameToNode {
		if err := nd.Terminate(); err != nil {
			fmt.Fprintln(nd, err)
			continue
		}
		c.DoneSafely(nd)
	}
}

// Terminate terminates the process.
func (c *Cluster) Terminate(name string) error {
	if v, ok := c.NameToNode[name]; ok {
		if err := v.Terminate(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%s does not exist in the Cluster!", name)
	}
	return nil
}

// Restart restarts the etcd Node.
func (c *Cluster) Restart(name string) error {
	if v, ok := c.NameToNode[name]; ok {
		if err := v.Restart(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%s does not exist in the Cluster!", name)
	}
	return nil
}

// RemoveAllDataDirs removes all etcd data directories.
// Used for cleaning up.
func (c *Cluster) RemoveAllDataDirs() {
	for k, nd := range c.NameToNode {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Fprintf(nd, "panic while deleting %s (%v)\n", nd.Flags.DataDir, err)
				}
			}()
			fmt.Fprintf(nd, "Deleting data-dir for %s (%s)\n", k, nd.Flags.DataDir)
			if err := os.RemoveAll(nd.Flags.DataDir); err != nil {
				fmt.Fprintf(nd, "error while deleting %s (%v)\n", nd.Flags.DataDir, err)
			}
		}()
	}
}

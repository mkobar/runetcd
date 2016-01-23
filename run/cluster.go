package run

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

// Cluster groups a set of process Members.
type Cluster struct {
	// define in pointer type to not copy over cluster
	mu                *sync.Mutex // guards the following
	maxProcNameLength *int
	wg                *sync.WaitGroup
	NameToMember      map[string]*Member
}

// CreateCluster creates a Cluster by parsing the input io.Reader.
func CreateCluster(w io.Writer, bufferStream chan string, outputOption OutputOption, command string, fs ...*Flags) (*Cluster, error) {

	if len(fs) == 0 {
		return nil, nil
	}

	var maxProcNameLength int
	c := &Cluster{
		mu:                &sync.Mutex{},
		maxProcNameLength: &maxProcNameLength,
		wg:                &sync.WaitGroup{},
		NameToMember:      make(map[string]*Member),
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

		m := &Member{
			pmu:                c.mu,
			pmaxProcNameLength: c.maxProcNameLength,

			outputOption: outputOption,
			colorIdx:     colorIdx,
			w:            w,
			BufferStream: bufferStream,

			Command: command,
			Flags:   f,

			PID:        0,
			Terminated: false,
			cmd:        nil,
		}
		c.NameToMember[name] = m

		colorIdx++
	}

	return c, nil
}

func (c *Cluster) WriteProc(w io.Writer) {
	lines := []string{}
	for name, m := range c.NameToMember {
		line := strings.TrimSpace(fmt.Sprintf("%s: %s %s", name, m.Command, m.Flags))
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
	if len(c.NameToMember) == 0 {
		return nil
	}

	// execute all members at the same time
	c.wg.Add(len(c.NameToMember))
	for name, m := range c.NameToMember {

		go func(name string, m *Member) {

			defer c.DoneSafely(m)

			cs := []string{"/bin/bash", "-c", m.Command + " " + m.Flags.String()}
			cmd := exec.Command(cs[0], cs[1:]...)
			cmd.Stdin = nil
			cmd.Stdout = m
			cmd.Stderr = m

			fmt.Fprintf(m, "Starting %s\n", name)
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(m, "Failed to start %s\n", name)
				return
			}
			m.cmd = cmd
			m.PID = cmd.Process.Pid

			cmd.Wait()

			fmt.Fprintf(m, "Exiting %s\n", name)
			return

		}(name, m)

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

// TerminateAll terminates all members in the group.
func (c *Cluster) TerminateAll() {
	defer func() {
		recover()
	}()

	if len(c.NameToMember) == 0 {
		return
	}

	for _, m := range c.NameToMember {
		if err := m.Terminate(); err != nil {
			fmt.Fprintln(m, err)
			continue
		}
		c.DoneSafely(m)
	}
}

// Terminate terminates the process.
func (c *Cluster) Terminate(name string) error {
	if v, ok := c.NameToMember[name]; ok {
		if err := v.Terminate(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%s does not exist in the Cluster!", name)
	}
	return nil
}

// Restart restarts the etcd member.
func (c *Cluster) Restart(name string) error {
	if v, ok := c.NameToMember[name]; ok {
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
	for k, m := range c.NameToMember {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Fprintf(m, "panic while deleting %s (%v)\n", m.Flags.DataDir, err)
				}
			}()
			fmt.Fprintf(m, "Deleting data-dir for %s (%s)\n", k, m.Flags.DataDir)
			if err := os.RemoveAll(m.Flags.DataDir); err != nil {
				fmt.Fprintf(m, "error while deleting %s (%v)\n", m.Flags.DataDir, err)
			}
		}()
	}
}

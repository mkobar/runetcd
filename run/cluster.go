package run

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
)

// Cluster groups a set of process Membes.
type Cluster struct {
	// define in pointer type to not copy over cluster
	mu                *sync.Mutex // guards the following
	maxProcNameLength *int
	wg                *sync.WaitGroup
	NameToMember      map[string]*Member
}

// CreateCluster creates a Cluster by parsing the input io.Reader.
func CreateCluster(w io.Writer, command string, fs ...*Flags) (*Cluster, error) {

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

			colorIdx: colorIdx,
			w:        w,

			Command: command,
			Flags:   f,

			PID:        0,
			Terminated: false,
			cmd:        nil,
		}
		c.NameToMember[name] = m
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

package etcdproc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"sync"
	"syscall"

	"github.com/fatih/color"
)

// Node represents an etcd node process.
type Node struct {
	// inherited from Group
	pmu                *sync.Mutex
	pmaxProcNameLength *int

	outputOption OutputOption
	colorIdx     int
	w            io.Writer
	BufferStream chan string

	Command string
	Flags   *Flags

	PID        int
	Terminated bool
	cmd        *exec.Cmd
}

type OutputOption int

const (
	ToTerminal OutputOption = iota
	ToHTML
)

var (
	colors = []color.Attribute{
		color.FgRed,
		color.FgGreen,
		color.FgYellow,
		color.FgBlue,
		color.FgMagenta,
	}
	colorMap = map[color.Attribute]string{
		color.FgRed:     "#ff0000",
		color.FgGreen:   "#008000",
		color.FgYellow:  "#ff9933",
		color.FgBlue:    "#0000ff",
		color.FgMagenta: "#ff00ff",
	}
)

func (nd *Node) Write(p []byte) (int, error) {
	buf := bytes.NewBuffer(p)
	wrote := 0
	for {
		line, err := buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return wrote, err
		}
		if len(line) > 1 {
			format := fmt.Sprintf("%%%ds | ", *(nd.pmaxProcNameLength))

			nd.pmu.Lock()

			switch nd.outputOption {
			case ToTerminal:
				color.Set(colors[nd.colorIdx])
				fmt.Fprintf(nd.w, format, nd.Flags.Name)
				color.Unset()
				fmt.Fprint(nd.w, string(line))
			case ToHTML:
				format = fmt.Sprintf(`<b><font color="%s">`, colorMap[colors[nd.colorIdx]]) + format + "</font>" + "%s</b>"
				nd.BufferStream <- fmt.Sprintf(format, nd.Flags.Name, line)
				if f, ok := nd.w.(http.Flusher); ok {
					if f != nil {
						f.Flush()
					}
				}
			}

			nd.pmu.Unlock()

			wrote += len(line)
		}
	}

	if len(p) > 0 && p[len(p)-1] != '\n' {
		switch nd.outputOption {
		case ToTerminal:
			nd.pmu.Lock()
			fmt.Fprintln(nd.w)
			nd.pmu.Unlock()
		}
	}

	return len(p), nil
}

func (nd *Node) Terminate() error {
	defer func() {
		if err := recover(); err != nil {
			switch nd.outputOption {
			case ToTerminal:
				fmt.Fprintf(nd, "panic while Terminate Node %s (%v)\n", nd.Flags.Name, err)
			case ToHTML:
				nd.BufferStream <- fmt.Sprintf("panic while Terminate Node %s (%v)\n", nd.Flags.Name, err)
				if f, ok := nd.w.(http.Flusher); ok {
					if f != nil {
						f.Flush()
					}
				}
			}
		}
	}()

	if nd.Terminated {
		return fmt.Errorf("%s is already terminated", nd.Flags.Name)
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintln(nd, "Terminate:", nd.Flags.Name)
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("Terminate: %s", nd.Flags.Name)
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	if err := syscall.Kill(nd.PID, syscall.SIGKILL); err != nil {
		return err
	}

	nd.Terminated = true

	return nil
}

// Restart restarts a Node.
func (nd *Node) Restart() error {
	defer func() {
		if err := recover(); err != nil {
			switch nd.outputOption {
			case ToTerminal:
				fmt.Fprintf(nd, "panic while Restart Node %s (%v)\n", nd.Flags.Name, err)
			case ToHTML:
				nd.BufferStream <- fmt.Sprintf("panic while Restart Node %s (%v)\n", nd.Flags.Name, err)
				if f, ok := nd.w.(http.Flusher); ok {
					if f != nil {
						f.Flush()
					}
				}
			}
		}
	}()

	if !nd.Terminated {
		return fmt.Errorf("%s is already running", nd.Flags.Name)
	}

	nd.Flags.InitialClusterState = "existing"

	cs := []string{"/bin/bash", "-c", nd.Command + " " + nd.Flags.String()}
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = nd
	cmd.Stderr = nd

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintln(nd, "Restart:", nd.Flags.Name)
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("Restart: %s", nd.Flags.Name)
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start %s with %v\n", nd.Flags.Name, err)
	}
	nd.cmd = cmd
	nd.PID = cmd.Process.Pid
	nd.Terminated = false

	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Fprintf(nd, "Exiting %s with %v\n", nd.Flags.Name, err)
			return
		}
		fmt.Fprintf(nd, "Exiting %s\n", nd.Flags.Name)
	}()

	return nil
}

package run

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

type OutputOption int

const (
	ToTerminal OutputOption = iota
	ToHTML
)

// Member represents a process or command.
type Member struct {
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

func (m *Member) Write(p []byte) (int, error) {
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
			format := fmt.Sprintf("%%%ds | ", *(m.pmaxProcNameLength))

			m.pmu.Lock()

			switch m.outputOption {
			case ToTerminal:
				color.Set(colors[m.colorIdx])
				fmt.Fprintf(m.w, format, m.Flags.Name)
				color.Unset()
				fmt.Fprint(m.w, string(line))
			case ToHTML:
				format = fmt.Sprintf(`<b><font color="%s">`, colorMap[colors[m.colorIdx]]) + format + "</font>" + "%s</b>"
				m.BufferStream <- fmt.Sprintf(format, m.Flags.Name, line)
				if f, ok := m.w.(http.Flusher); ok {
					if f != nil {
						f.Flush()
					}
				}
			}

			m.pmu.Unlock()

			wrote += len(line)
		}
	}

	if len(p) > 0 && p[len(p)-1] != '\n' {
		switch m.outputOption {
		case ToTerminal:
			m.pmu.Lock()
			fmt.Fprintln(m.w)
			m.pmu.Unlock()
		}
	}

	return len(p), nil
}

func (m *Member) Terminate() error {
	defer func() {
		if err := recover(); err != nil {
			switch m.outputOption {
			case ToTerminal:
				fmt.Fprintf(m, "panic while Terminate member %s (%v)\n", m.Flags.Name, err)
			case ToHTML:
				m.BufferStream <- fmt.Sprintf("panic while Terminate member %s (%v)\n", m.Flags.Name, err)
				if f, ok := m.w.(http.Flusher); ok {
					if f != nil {
						f.Flush()
					}
				}
			}
		}
	}()

	if m.Terminated {
		return fmt.Errorf("%s is already terminated", m.Flags.Name)
	}

	switch m.outputOption {
	case ToTerminal:
		fmt.Fprintln(m, "Terminate:", m.Flags.Name)
	case ToHTML:
		m.BufferStream <- fmt.Sprintf("Terminate: %s", m.Flags.Name)
		if f, ok := m.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	if err := syscall.Kill(m.PID, syscall.SIGKILL); err != nil {
		return err
	}

	m.Terminated = true

	return nil
}

// Restart restarts a member.
func (m *Member) Restart() error {
	defer func() {
		if err := recover(); err != nil {
			switch m.outputOption {
			case ToTerminal:
				fmt.Fprintf(m, "panic while Restart member %s (%v)\n", m.Flags.Name, err)
			case ToHTML:
				m.BufferStream <- fmt.Sprintf("panic while Restart member %s (%v)\n", m.Flags.Name, err)
				if f, ok := m.w.(http.Flusher); ok {
					if f != nil {
						f.Flush()
					}
				}
			}
		}
	}()

	if !m.Terminated {
		return fmt.Errorf("%s is already running", m.Flags.Name)
	}

	m.Flags.InitialClusterState = "existing"

	cs := []string{"/bin/bash", "-c", m.Command + " " + m.Flags.String()}
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = m
	cmd.Stderr = m

	switch m.outputOption {
	case ToTerminal:
		fmt.Fprintln(m, "Restart:", m.Flags.Name)
	case ToHTML:
		m.BufferStream <- fmt.Sprintf("Restart: %s", m.Flags.Name)
		if f, ok := m.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start %s with %v\n", m.Flags.Name, err)
	}
	m.cmd = cmd
	m.PID = cmd.Process.Pid
	m.Terminated = false

	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Fprintf(m, "Exiting %s with %v\n", m.Flags.Name, err)
			return
		}
		fmt.Fprintf(m, "Exiting %s\n", m.Flags.Name)
	}()

	return nil
}

package etcdproc

import (
	"os"
	"testing"
)

func TestNodeMap(t *testing.T) {
	nd := &Node{}
	nd.w = os.Stdout
	nd.PID = 0
	mp := make(map[string]*Node)
	mp["TEST"] = nd
	if v, ok := mp["TEST"]; ok {
		v.Terminated = true
	}
	if !mp["TEST"].Terminated {
		t.Errorf("expect 'terminated' true, but got %v", mp["TEST"].Terminated)
	}
}

package run

import (
	"os"
	"testing"
)

func TestMemberMap(t *testing.T) {
	m := &Member{}
	m.w = os.Stdout
	m.PID = 0
	mp := make(map[string]*Member)
	mp["TEST"] = m
	if v, ok := mp["TEST"]; ok {
		v.Terminated = true
	}
	if !mp["TEST"].Terminated {
		t.Errorf("expect 'terminated' true, but got %v", mp["TEST"].Terminated)
	}
}

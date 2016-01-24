package etcdproc

import (
	"fmt"
	"os"
	"testing"

	"github.com/satori/go.uuid"
)

func TestCreateCluster(t *testing.T) {
	isClientTLS := false
	isPeerTLS := false
	certPath := ""
	privateKeyPath := ""
	caPath := ""
	fs := make([]*Flags, 3)
	for i := range fs {
		df, err := NewFlags(fmt.Sprintf("etcd%d", i), nil, 11+i, "etcd-cluster-token", "new", uuid.NewV4().String(), isClientTLS, isPeerTLS, certPath, privateKeyPath, caPath)
		if err != nil {
			t.Error(err)
		}
		fs[i] = df
	}
	CombineFlags(fs...)
	c, err := CreateCluster(os.Stdout, nil, ToTerminal, "bin/etcd", fs...)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("c.WriteProc()")
	c.WriteProc(os.Stdout)

	fmt.Println()
	for k, v := range c.NameToNode {
		fmt.Printf("%s's ListenClientURLs: %v\n", k, v.Flags.ListenClientURLs)
	}
}

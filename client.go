package run

import (
	"fmt"
	"os"

	"github.com/coreos/etcd/Godeps/_workspace/src/google.golang.org/grpc"
)

func mustCreateConn(endpoint string) *grpc.ClientConn {
	conn, err := grpc.Dial(endpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial error: %v\n", err)
		os.Exit(1)
	}
	return conn
}

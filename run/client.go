package run

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/Godeps/_workspace/src/google.golang.org/grpc"
	"github.com/coreos/etcd/etcdserver/etcdserverpb"
)

func mustCreateConn(endpoint string) *grpc.ClientConn {
	conn, err := grpc.Dial(endpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial error: %v\n", err)
		os.Exit(1)
	}
	return conn
}

// Put puts key-value.
func (c *Cluster) Put(

	w io.Writer,
	name string,

	key []byte,
	val []byte,

) error {

	endpoint := ""
	if m, ok := c.NameToMember[name]; !ok {
		return fmt.Errorf("%s does not exist in the Cluster!", name)
	} else {
		if m.Flags.ExperimentalgRPCAddr == "" {
			return fmt.Errorf("no experimental-gRPC-addr found for %s", name)
		}
		endpoint = m.Flags.ExperimentalgRPCAddr
	}

	conn := mustCreateConn(endpoint)
	client := etcdserverpb.NewKVClient(conn)

	r := &etcdserverpb.PutRequest{
		Key:   key,
		Value: val,
	}

	fmt.Fprintf(w, "[Put] Started! (endpoint: %s)\n", endpoint)

	ts := time.Now()
	if _, err := client.Put(context.Background(), r); err != nil {
		return err
	}

	nk := string(key)
	if len(key) > 5 {
		nk = nk[:5] + "..."
	}
	nv := string(val)
	if len(val) > 5 {
		nv = nv[:5] + "..."
	}
	fmt.Fprintf(w, "[Put] Done! Took %v for %s/%s (endpoint: %s).\n", time.Since(ts), nk, nv, endpoint)
	fmt.Fprintf(w, "\n")

	return nil
}

// Stress generates random data and loads them to
// given member with the name.
func (c *Cluster) Stress(

	w io.Writer,
	name string,

	connsN int,
	clientsN int,

	stressN int,
	stressKeyN int,
	stressValN int,

) error {

	endpoint := ""
	if m, ok := c.NameToMember[name]; !ok {
		return fmt.Errorf("%s does not exist in the Cluster!", name)
	} else {
		if m.Flags.ExperimentalgRPCAddr == "" {
			return fmt.Errorf("no experimental-gRPC-addr found for %s", name)
		}
		endpoint = m.Flags.ExperimentalgRPCAddr
	}

	fmt.Fprintf(w, "[Stress] Started generating %d random data...\n", stressN)
	sr := time.Now()
	keys := make([][]byte, stressN)
	vals := make([][]byte, stressN)
	for i := range keys {
		keys[i] = RandBytes(stressKeyN)
		vals[i] = RandBytes(stressValN)
	}
	fmt.Fprintf(w, "[Stress] Done with generating %d random data! Took %v\n", stressN, time.Since(sr))

	conns := make([]*grpc.ClientConn, connsN)
	for i := range conns {
		conns[i] = mustCreateConn(endpoint)
	}
	clients := make([]etcdserverpb.KVClient, clientsN)
	for i := range clients {
		clients[i] = etcdserverpb.NewKVClient(conns[i%int(connsN)])
	}

	fmt.Fprintf(w, "[Stress] Started stressing with GRPC (endpoint %s).\n", endpoint)

	requests := make(chan *etcdserverpb.PutRequest, stressN)
	done, errChan := make(chan struct{}), make(chan error)

	for i := range clients {
		go func(i int, requests <-chan *etcdserverpb.PutRequest) {
			for r := range requests {
				if _, err := clients[i].Put(context.Background(), r); err != nil {
					errChan <- err
					return
				}
			}
			done <- struct{}{}
		}(i, requests)
	}

	st := time.Now()

	for i := 0; i < stressN; i++ {
		r := &etcdserverpb.PutRequest{
			Key:   keys[i],
			Value: vals[i],
		}
		requests <- r
	}

	close(requests)

	cn := 0
	for cn != len(clients) {
		select {
		case err := <-errChan:
			return err
		case <-done:
			cn++
		}
	}

	tt := time.Since(st)
	pt := tt / time.Duration(stressN)
	fmt.Fprintf(
		w,
		"[Stress] Done! Took %v for %d requests(%v per each) with %d connection(s), %d client(s) (endpoint: %s)\n",
		tt, stressN, pt, connsN, clientsN, endpoint,
	)
	fmt.Fprintf(w, "\n")
	return nil
}

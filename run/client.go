package run

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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

// WatchAndPut watches key and later put that key
// so that the watcher can return.
func (c *Cluster) WatchAndPut(w io.Writer, name string, connsN, streamsN, watchersN int) error {
	keyToWatch := []byte("fo")
	endpoint := ""
	if m, ok := c.NameToMember[name]; !ok {
		return fmt.Errorf("%s does not exist in the Cluster!", name)
	} else {
		if m.Flags.ExperimentalgRPCAddr == "" {
			return fmt.Errorf("no experimental-gRPC-addr found for %s", name)
		}
		endpoint = m.Flags.ExperimentalgRPCAddr
	}
	fmt.Fprintf(w, "[Watch-Request] Started! (endpoint: %s)\n", endpoint)

	conns := make([]*grpc.ClientConn, connsN)
	for i := range conns {
		conns[i] = mustCreateConn(endpoint)
	}

	streams := make([]etcdserverpb.Watch_WatchClient, streamsN)
	for i := range streams {
		watchClient := etcdserverpb.NewWatchClient(conns[i%int(connsN)])
		wStream, err := watchClient.Watch(context.Background())
		if err != nil {
			return err
		}
		streams[i] = wStream
	}

	fmt.Fprintf(w, "[Watch-Request] Launching all watchers! (endpoint: %s)\n", endpoint)
	for i := 0; i < watchersN; i++ {
		go func(i int) {
			wStream := streams[i%int(streamsN)]
			wr := &etcdserverpb.WatchRequest{
				CreateRequest: &etcdserverpb.WatchCreateRequest{Prefix: keyToWatch},
			}
			if err := wStream.Send(wr); err != nil {
				fmt.Fprintf(w, "[Watch-Send] error (%v)\n", err)
			}
		}(i)
	}

	streamsToWatchId := make(map[etcdserverpb.Watch_WatchClient]map[int64]struct{})
	for i := 0; i < watchersN; i++ {
		wStream := streams[i%int(streamsN)]
		wresp, err := wStream.Recv()
		if err != nil {
			return err
		}
		if !wresp.Created {
			fmt.Fprintf(w, "[Watch-Request] wresp.Created is supposed to be true! Something wrong (endpoint: %s)\n", endpoint)
		}
		if _, ok := streamsToWatchId[wStream]; !ok {
			streamsToWatchId[wStream] = make(map[int64]struct{})
		}
		streamsToWatchId[wStream][wresp.WatchId] = struct{}{}
	}

	fmt.Fprintln(w, "[Put-Request] trigger notifications with PUT!")
	kvc := etcdserverpb.NewKVClient(conns[0])
	if _, err := kvc.Put(context.Background(), &etcdserverpb.PutRequest{Key: []byte("foo"), Value: []byte("bar")}); err != nil {
		return err
	}

	ts := time.Now()
	fmt.Fprintf(w, "[Watch] Started! (endpoint: %s)\n", endpoint)
	fmt.Fprintf(w, "\n")

	var wg sync.WaitGroup
	wg.Add(watchersN)
	for i := 0; i < watchersN; i++ {
		go func(i int) {
			defer wg.Done()
			wStream := streams[i%int(streamsN)]
			wresp, err := wStream.Recv()
			if err != nil {
				fmt.Fprintf(w, "[Watch] send error (%v)", err)
				return
			}
			switch {
			case wresp.Created:
				fmt.Fprintf(w, "[revision] %d / watcher created %08x\n", wresp.Header.Revision, wresp.WatchId)
			case wresp.Canceled:
				fmt.Fprintf(w, "[revision] %d / watcher canceled %08x\n", wresp.Header.Revision, wresp.WatchId)
			default:
				fmt.Fprintf(w, "[revision] %d\n", wresp.Header.Revision)
				for _, ev := range wresp.Events {
					fmt.Fprintf(w, "%s: %s %s\n", ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
				}
			}
		}(i)
	}
	wg.Wait()
	fmt.Fprintf(w, "[Watch] Done! Took %v (endpoint: %s)\n", time.Since(ts), endpoint)
	fmt.Fprintf(w, "\n")
	return nil
}

// ServerStats encapsulates various statistics about an EtcdServer and its
// communication with other members of the cluster.
// (https://github.com/coreos/etcd/tree/master/etcdserver/stats)
type ServerStats struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	State     string    `json:"state"`
	StartTime time.Time `json:"startTime"`

	LeaderInfo struct {
		ID        string    `json:"leader"`
		Uptime    string    `json:"uptime"`
		StartTime time.Time `json:"startTime"`
	} `json:"leaderInfo"`

	RecvAppendRequestCnt uint64  `json:"recvAppendRequestCnt,"`
	RecvingPkgRate       float64 `json:"recvPkgRate,omitempty"`
	RecvingBandwidthRate float64 `json:"recvBandwidthRate,omitempty"`

	SendAppendRequestCnt uint64  `json:"sendAppendRequestCnt"`
	SendingPkgRate       float64 `json:"sendPkgRate,omitempty"`
	SendingBandwidthRate float64 `json:"sendBandwidthRate,omitempty"`
}

// GetStats returns the leader of the cluster.
func (c *Cluster) GetStats() (map[string]ServerStats, []string, error) {
	nameToEndpoint := make(map[string][]string)
	for n, m := range c.NameToMember {
		for v := range m.Flags.ListenClientURLs {
			if _, ok := nameToEndpoint[n]; !ok {
				nameToEndpoint[n] = []string{}
			}
			nameToEndpoint[n] = append(nameToEndpoint[n], v)
		}
	}

	rm := make(map[string]ServerStats)
	for name, endpoints := range nameToEndpoint {
		var errMsg error
		for _, endpoint := range endpoints {
			resp, err := http.Get(endpoint + "/v2/stats/self")
			if err != nil {
				errMsg = err
				continue
			}

			sts := ServerStats{}
			if err := json.NewDecoder(resp.Body).Decode(&sts); err != nil {
				errMsg = err
				continue
			}
			resp.Body.Close()
			rm[name] = sts
		}
		if errMsg != nil {
			return nil, nil, errMsg
		}
	}

	leaderNames := []string{}
	for name, stat := range rm {
		if stat.State == "StateLeader" {
			leaderNames = append(leaderNames, name)
		}
	}

	return rm, leaderNames, nil
}

// GetMetrics returns the metrics of the cluster.
//
// Some useful metrics:
// 	- etcd_storage_keys_total
// 	- etcd_storage_db_total_size_in_bytes
func (c *Cluster) GetMetrics() (map[string]map[string]float64, error) {
	nameToEndpoint := make(map[string][]string)
	for n, m := range c.NameToMember {
		for v := range m.Flags.ListenClientURLs {
			if _, ok := nameToEndpoint[n]; !ok {
				nameToEndpoint[n] = []string{}
			}
			nameToEndpoint[n] = append(nameToEndpoint[n], v)
		}
	}

	rm := make(map[string]map[string]float64)
	for name, endpoints := range nameToEndpoint {
		var errMsg error

		for _, endpoint := range endpoints {
			resp, err := http.Get(endpoint + "/metrics")
			if err != nil {
				errMsg = err
				continue
			}

			scanner := bufio.NewScanner(resp.Body)

			mm := make(map[string]float64)
			for scanner.Scan() {
				txt := scanner.Text()
				if strings.HasPrefix(txt, "#") {
					continue
				}
				ts := strings.SplitN(txt, " ", 2)
				fv := 0.0
				if len(ts) == 2 {
					v, err := strconv.ParseFloat(ts[1], 64)
					if err == nil {
						fv = v
					}
				}
				mm[ts[0]] = fv
			}
			if err := scanner.Err(); err != nil {
				errMsg = err
				continue
			}

			resp.Body.Close()
			rm[name] = mm
		}

		if errMsg != nil {
			return nil, errMsg
		}
	}

	return rm, nil
}

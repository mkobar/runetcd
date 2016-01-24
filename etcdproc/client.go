package etcdproc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/Godeps/_workspace/src/google.golang.org/grpc"
	pb "github.com/coreos/etcd/etcdserver/etcdserverpb"
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
func (c *Cluster) Put(key []byte, val []byte) error {
	endpoint := ""
	var nd *Node
	for _, node := range c.NameToNode {
		if node.Flags.ExperimentalgRPCAddr != "" && !node.Terminated {
			endpoint = node.Flags.ExperimentalgRPCAddr
			nd = node
			break
		}
	}
	if endpoint == "" || nd == nil {
		return fmt.Errorf("no experimental-gRPC found")
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, "[Put] Started! (name: %s, endpoint: %s)\n", nd.Flags.Name, endpoint)
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[Put] Started! (name: %s, endpoint: %s)", nd.Flags.Name, endpoint)
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	conn := mustCreateConn(endpoint)
	kvc := pb.NewKVClient(conn)
	req := &pb.PutRequest{
		Key:   key,
		Value: val,
	}

	st := time.Now()
	if _, err := kvc.Put(context.Background(), req); err != nil {
		return err
	}

	nk := string(key)
	printLimit := 15
	if len(key) > printLimit {
		nk = nk[:printLimit] + "..."
	}
	nv := string(val)
	if len(val) > printLimit {
		nv = nv[:printLimit] + "..."
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, "[Put] %s : %s / Took %v (name: %s, endpoint: %s)\n", nk, nv, time.Since(st), nd.Flags.Name, endpoint)
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[Put] %s : %s / Took %v (name: %s, endpoint: %s)", nk, nv, time.Since(st), nd.Flags.Name, endpoint)
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	return nil
}

// Range ranges on the key.
func (c *Cluster) Range(key []byte) error {
	endpoint := ""
	var nd *Node
	for _, node := range c.NameToNode {
		if node.Flags.ExperimentalgRPCAddr != "" && !node.Terminated {
			endpoint = node.Flags.ExperimentalgRPCAddr
			nd = node
			break
		}
	}
	if endpoint == "" || nd == nil {
		return fmt.Errorf("no experimental-gRPC found")
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, "[Range] Started on %s\n", key)
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[Range] Started on %s", key)
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	rangeStart := key
	rangeEnd := []byte{}
	sortByOrder := pb.RangeRequest_ASCEND
	sortByTarget := pb.RangeRequest_KEY
	rangeLimit := 0

	conn := mustCreateConn(endpoint)
	kvc := pb.NewKVClient(conn)
	req := &pb.RangeRequest{
		Key:        rangeStart,
		RangeEnd:   rangeEnd,
		SortOrder:  sortByOrder,
		SortTarget: sortByTarget,
		Limit:      int64(rangeLimit),
	}

	st := time.Now()
	resp, err := kvc.Range(context.Background(), req)
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		switch nd.outputOption {
		case ToTerminal:
			fmt.Fprintf(nd.w, fmt.Sprintf("[Range] %s : %s\n", ev.Key, ev.Value))
		case ToHTML:
			nd.BufferStream <- fmt.Sprintf("[Range] %s : %s", ev.Key, ev.Value)
			if f, ok := nd.w.(http.Flusher); ok {
				if f != nil {
					f.Flush()
				}
			}
		}
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, "[Range] Took %v (name: %s, endpoint: %s)\n", time.Since(st), nd.Flags.Name, endpoint)
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[Range] Took %v (name: %s, endpoint: %s)", time.Since(st), nd.Flags.Name, endpoint)
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	return nil
}

// Stress generates random data to one selected Node.
func (c *Cluster) Stress(connsN int, clientsN int, stressN int, stressKeyN int, stressValN int) error {
	endpoint := ""
	var nd *Node
	for _, node := range c.NameToNode {
		if node.Flags.ExperimentalgRPCAddr != "" && !node.Terminated {
			endpoint = node.Flags.ExperimentalgRPCAddr
			nd = node
			break
		}
	}
	if endpoint == "" || nd == nil {
		return fmt.Errorf("no experimental-gRPC found")
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, "[Stress] Started generating %d random data to %s\n", stressN, nd.Flags.Name)
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[Stress] Started generating %d random data to %s", stressN, nd.Flags.Name)
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	sr := time.Now()
	keys := make([][]byte, stressN)
	vals := make([][]byte, stressN)
	for i := range keys {
		keys[i] = RandBytes(stressKeyN)
		vals[i] = RandBytes(stressValN)
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, "[Stress] Done with generating %d random data! Took %v\n", stressN, time.Since(sr))
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[Stress] Done with generating %d random data! Took %v", stressN, time.Since(sr))
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	conns := make([]*grpc.ClientConn, connsN)
	for i := range conns {
		conns[i] = mustCreateConn(endpoint)
	}
	clients := make([]pb.KVClient, clientsN)
	for i := range clients {
		clients[i] = pb.NewKVClient(conns[i%int(connsN)])
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, "[Stress] Started stressing with GRPC (endpoint %s)\n", endpoint)
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[Stress] Started stressing with GRPC (endpoint %s)", endpoint)
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	requests := make(chan *pb.PutRequest, stressN)
	done, errChan := make(chan struct{}), make(chan error)

	for i := range clients {
		go func(i int, requests <-chan *pb.PutRequest) {
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
		r := &pb.PutRequest{
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

	rMsg := fmt.Sprintf("[Stress] Done! Took %v for %d requests(%v per each) with %d connection(s), %d client(s) (endpoint: %s)", tt, stressN, pt, connsN, clientsN, endpoint)
	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintln(nd.w, rMsg)
	case ToHTML:
		nd.BufferStream <- rMsg
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	return nil
}

// SimpleStress generates random data to one selected Node.
func (c *Cluster) SimpleStress() error {
	endpoint := ""
	var nd *Node
	for _, node := range c.NameToNode {
		if node.Flags.ExperimentalgRPCAddr != "" && !node.Terminated {
			endpoint = node.Flags.ExperimentalgRPCAddr
			nd = node
			break
		}
	}
	if endpoint == "" || nd == nil {
		return fmt.Errorf("no experimental-gRPC found")
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, "[SimpleStress] Started generating simple random data to %s\n", nd.Flags.Name)
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[SimpleStress] Started generating simple random data to %s", nd.Flags.Name)
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	stressN := 10
	connsN := 1
	clientsN := 1

	sr := time.Now()
	keys := make([][]byte, stressN)
	vals := make([][]byte, stressN)
	for i := range keys {
		keys[i] = []byte(fmt.Sprintf("sample_%d_%s", i, RandBytes(5)))
		vals[i] = []byte(fmt.Sprintf(`{"value": "created at %d"}`, time.Now().Nanosecond()))
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, "[SimpleStress] Done with generating %d random data! Took %v\n", stressN, time.Since(sr))
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[SimpleStress] Done with generating %d random data! Took %v", stressN, time.Since(sr))
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	conns := make([]*grpc.ClientConn, connsN)
	for i := range conns {
		conns[i] = mustCreateConn(endpoint)
	}
	clients := make([]pb.KVClient, clientsN)
	for i := range clients {
		clients[i] = pb.NewKVClient(conns[i%int(connsN)])
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, "[SimpleStress] Started stressing with GRPC (endpoint %s)\n", endpoint)
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[SimpleStress] Started stressing with GRPC (endpoint %s)", endpoint)
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	requests := make(chan *pb.PutRequest, stressN)
	done, errChan := make(chan struct{}), make(chan error)

	for i := range clients {
		go func(i int, requests <-chan *pb.PutRequest) {
			for r := range requests {
				if _, err := clients[i].Put(context.Background(), r); err != nil {
					errChan <- err
					return
				}
				switch nd.outputOption {
				case ToTerminal:
					fmt.Fprintf(nd.w, "[SimpleStress PUT] %s : %s\n", r.Key, r.Value)
				case ToHTML:
					nd.BufferStream <- fmt.Sprintf("[SimpleStress PUT] %s : %s", r.Key, r.Value)
					if f, ok := nd.w.(http.Flusher); ok {
						if f != nil {
							f.Flush()
						}
					}
				}
			}
			done <- struct{}{}
		}(i, requests)
	}

	st := time.Now()

	for i := 0; i < stressN; i++ {
		r := &pb.PutRequest{
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
	fMsg := fmt.Sprintf(
		"[SimpleStress] Done! Took %v for %d requests(%v per each) with %d connection(s), %d client(s) (endpoint: %s)",
		tt, stressN, pt, connsN, clientsN, endpoint,
	)
	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintln(nd.w, fMsg)
	case ToHTML:
		nd.BufferStream <- fMsg
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	return nil
}

// WatchAndPut watches key and later put that key
// so that the watcher can return.
func (c *Cluster) WatchAndPut(connsN, streamsN, watchersN int) error {
	endpoint := ""
	var nd *Node
	for _, node := range c.NameToNode {
		if node.Flags.ExperimentalgRPCAddr != "" && !node.Terminated {
			endpoint = node.Flags.ExperimentalgRPCAddr
			nd = node
			break
		}
	}
	if endpoint == "" || nd == nil {
		return fmt.Errorf("no experimental-gRPC found")
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintln(nd.w, "[WatchAndPut] Started!")
	case ToHTML:
		nd.BufferStream <- "[WatchAndPut] Started!"
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	conns := make([]*grpc.ClientConn, connsN)
	for i := range conns {
		conns[i] = mustCreateConn(endpoint)
	}

	streams := make([]pb.Watch_WatchClient, streamsN)
	for i := range streams {
		watchClient := pb.NewWatchClient(conns[i%int(connsN)])
		wStream, err := watchClient.Watch(context.Background())
		if err != nil {
			return err
		}
		streams[i] = wStream
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, fmt.Sprintf("[WatchAndPut] Launching all watchers! (name: %s, endpoint: %s)\n", nd.Flags.Name, endpoint))
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[WatchAndPut] Launching all watchers! (name: %s, endpoint: %s)", nd.Flags.Name, endpoint)
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	keyToWatch := []byte("fo")
	for i := 0; i < watchersN; i++ {
		go func(i int) {
			wStream := streams[i%int(streamsN)]
			wr := &pb.WatchRequest{
				CreateRequest: &pb.WatchCreateRequest{Prefix: keyToWatch},
			}
			if err := wStream.Send(wr); err != nil {
				switch nd.outputOption {
				case ToTerminal:
					fmt.Fprintf(nd.w, fmt.Sprintf("[wStream.Send] error (%v)\n", err))
				case ToHTML:
					nd.BufferStream <- fmt.Sprintf("[wStream.Send] error (%v)", err)
					if f, ok := nd.w.(http.Flusher); ok {
						if f != nil {
							f.Flush()
						}
					}
				}
			}
		}(i)
	}

	streamsToWatchId := make(map[pb.Watch_WatchClient]map[int64]struct{})
	for i := 0; i < watchersN; i++ {
		wStream := streams[i%int(streamsN)]
		wresp, err := wStream.Recv()
		if err != nil {
			return err
		}
		if !wresp.Created {
			switch nd.outputOption {
			case ToTerminal:
				fmt.Fprintln(nd.w, "[WatchResponse] wresp.Created is supposed to be true! Something wrong!")
			case ToHTML:
				nd.BufferStream <- "[WatchResponse] wresp.Created is supposed to be true! Something wrong!"
				if f, ok := nd.w.(http.Flusher); ok {
					if f != nil {
						f.Flush()
					}
				}
			}
		}
		if _, ok := streamsToWatchId[wStream]; !ok {
			streamsToWatchId[wStream] = make(map[int64]struct{})
		}
		streamsToWatchId[wStream][wresp.WatchId] = struct{}{}
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintln(nd.w, "[PutRequest] Triggering notifications with PUT!")
	case ToHTML:
		nd.BufferStream <- "[PutRequest] Triggering notifications with PUT!"
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	kvc := pb.NewKVClient(conns[0])
	if _, err := kvc.Put(context.Background(), &pb.PutRequest{Key: []byte("foo"), Value: []byte("bar")}); err != nil {
		return err
	}

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintln(nd.w, "[Watch] Started!")
	case ToHTML:
		nd.BufferStream <- "[Watch] Started!"
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	st := time.Now()
	var wg sync.WaitGroup
	wg.Add(watchersN)

	for i := 0; i < watchersN; i++ {

		go func(i int) {

			defer wg.Done()

			wStream := streams[i%int(streamsN)]
			wresp, err := wStream.Recv()
			if err != nil {
				switch nd.outputOption {
				case ToTerminal:
					fmt.Fprintf(nd.w, fmt.Sprintf("[Watch] send error (%v)\n", err))
				case ToHTML:
					nd.BufferStream <- fmt.Sprintf("[Watch] send error (%v)", err)
					if f, ok := nd.w.(http.Flusher); ok {
						if f != nil {
							f.Flush()
						}
					}
				}
				return
			}

			switch {

			case wresp.Created:
				switch nd.outputOption {
				case ToTerminal:
					fmt.Fprintf(nd.w, fmt.Sprintf("[Watch revision] %d / watcher created %08x\n", wresp.Header.Revision, wresp.WatchId))
				case ToHTML:
					nd.BufferStream <- fmt.Sprintf("[Watch revision] %d / watcher created %08x", wresp.Header.Revision, wresp.WatchId)
					if f, ok := nd.w.(http.Flusher); ok {
						if f != nil {
							f.Flush()
						}
					}
				}

			case wresp.Canceled:
				switch nd.outputOption {
				case ToTerminal:
					fmt.Fprintf(nd.w, fmt.Sprintf("[Watch revision] %d / watcher canceled %08x\n", wresp.Header.Revision, wresp.WatchId))
				case ToHTML:
					nd.BufferStream <- fmt.Sprintf("[Watch revision] %d / watcher canceled %08x", wresp.Header.Revision, wresp.WatchId)
					if f, ok := nd.w.(http.Flusher); ok {
						if f != nil {
							f.Flush()
						}
					}
				}

			default:
				switch nd.outputOption {
				case ToTerminal:
					fmt.Fprintf(nd.w, fmt.Sprintf("[Watch revision] %d\n", wresp.Header.Revision))
				case ToHTML:
					nd.BufferStream <- fmt.Sprintf("[Watch revision] %d", wresp.Header.Revision)
					if f, ok := nd.w.(http.Flusher); ok {
						if f != nil {
							f.Flush()
						}
					}
				}
				for _, ev := range wresp.Events {
					switch nd.outputOption {
					case ToTerminal:
						fmt.Fprintf(nd.w, fmt.Sprintf("[%s] %s : %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value))
					case ToHTML:
						nd.BufferStream <- fmt.Sprintf("[%s] %s : %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
						if f, ok := nd.w.(http.Flusher); ok {
							if f != nil {
								f.Flush()
							}
						}
					}
				}
			}

		}(i)

	}

	wg.Wait()

	switch nd.outputOption {
	case ToTerminal:
		fmt.Fprintf(nd.w, fmt.Sprintf("[Watch] Done! Took %v!\n", time.Since(st)))
	case ToHTML:
		nd.BufferStream <- fmt.Sprintf("[Watch] Done! Took %v!", time.Since(st))
		if f, ok := nd.w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}
	}

	return nil
}

// ServerStats encapsulates various statistics about an EtcdServer and its
// communication with other nodes of the cluster.
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

// GetStats gets the statistics of the cluster.
func (c *Cluster) GetStats() (map[string]ServerStats, map[string]string, error) {
	nameToEndpoints := make(map[string][]string)
	for n, nd := range c.NameToNode {
		for v := range nd.Flags.ListenClientURLs {
			if _, ok := nameToEndpoints[n]; !ok {
				nameToEndpoints[n] = []string{}
			}
			nameToEndpoints[n] = append(nameToEndpoints[n], v)
		}
	}

	rm := make(map[string]ServerStats)
	nameToEndpoint := make(map[string]string)
	for name, endpoints := range nameToEndpoints {
		for _, endpoint := range endpoints {
			sts := ServerStats{}
			resp, err := http.Get(endpoint + "/v2/stats/self")
			if err != nil {
				sts.Name = name
				sts.ID = ""
				sts.State = "Unreachable"
				rm[endpoint] = sts
				nameToEndpoint[name] = endpoint
				continue
			}

			if err := json.NewDecoder(resp.Body).Decode(&sts); err != nil {
				sts.Name = name
				sts.ID = ""
				sts.State = "Unreachable"
				rm[endpoint] = sts
				nameToEndpoint[name] = endpoint
				continue
			}
			resp.Body.Close()

			rm[endpoint] = sts
			nameToEndpoint[name] = endpoint
			break // use only one client URL
		}
	}

	return rm, nameToEndpoint, nil
}

// GetStats gets the statistics of the cluster.
func GetStats(endpoints ...string) (map[string]ServerStats, map[string]string, error) {
	rm := make(map[string]ServerStats)
	nameToEndpoint := make(map[string]string)
	for _, endpoint := range endpoints {
		sts := ServerStats{}
		resp, err := http.Get(endpoint + "/v2/stats/self")
		if err != nil {
			sts.ID = ""
			sts.State = "Unreachable"
			rm[endpoint] = sts
			continue
		}

		if err := json.NewDecoder(resp.Body).Decode(&sts); err != nil {
			sts.ID = ""
			sts.State = "Unreachable"
			rm[endpoint] = sts
			continue
		}
		resp.Body.Close()

		rm[endpoint] = sts
		nameToEndpoint[sts.Name] = endpoint
	}

	return rm, nameToEndpoint, nil
}

var defaultMetrics = map[string]float64{
	"etcd_storage_keys_total":              0.0,
	"etcd_storage_db_total_size_in_bytes":  0.0,
	"etcd_wal_fsync_durations_seconds_sum": 0.0,
	"go_gc_duration_seconds_sum":           0.0,
	"go_memstats_alloc_bytes":              0.0,
	"go_memstats_heap_alloc_bytes":         0.0,
	"go_memstats_mallocs_total":            0.0,
	"process_cpu_seconds_total":            0.0,
	"go_goroutines":                        0.0,
	"etcd_storage_watcher_total":           0.0,
	"etcd_storage_watch_stream_total":      0.0,
	"etcd_storage_slow_watcher_total":      0.0,
}

// GetMetrics returns the metrics of the cluster.
//
// Some useful metrics:
// 	- etcd_storage_keys_total
// 	- etcd_storage_db_total_size_in_bytes
// 	- etcd_wal_fsync_durations_seconds_sum
// 	- go_gc_duration_seconds_sum
// 	- go_memstats_alloc_bytes
// 	- go_memstats_heap_alloc_bytes
// 	- go_memstats_mallocs_total
// 	- process_cpu_seconds_total
// 	- go_goroutines
// 	- etcd_storage_watcher_total
// 	- etcd_storage_watch_stream_total
// 	- etcd_storage_slow_watcher_total
//
func (c *Cluster) GetMetrics() (map[string]map[string]float64, map[string]string, error) {
	nameToEndpoints := make(map[string][]string)
	for n, m := range c.NameToNode {
		for v := range m.Flags.ListenClientURLs {
			if _, ok := nameToEndpoints[n]; !ok {
				nameToEndpoints[n] = []string{}
			}
			nameToEndpoints[n] = append(nameToEndpoints[n], v)
		}
	}

	rm := make(map[string]map[string]float64)
	nameToEndpoint := make(map[string]string)
	for name, endpoints := range nameToEndpoints {
		for _, endpoint := range endpoints {
			resp, err := http.Get(endpoint + "/metrics")
			if err != nil {
				rm[endpoint] = defaultMetrics
				nameToEndpoint[name] = endpoint
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
				rm[endpoint] = defaultMetrics
				nameToEndpoint[name] = endpoint
				continue
			}
			resp.Body.Close()

			rm[endpoint] = mm
			nameToEndpoint[name] = endpoint
			break // use only one client URL
		}
	}

	return rm, nameToEndpoint, nil
}

// GetMetrics returns the metrics of the cluster.
//
// Some useful metrics:
// 	- etcd_storage_keys_total
// 	- etcd_storage_db_total_size_in_bytes
// 	- etcd_wal_fsync_durations_seconds_sum
// 	- go_gc_duration_seconds_sum
// 	- go_memstats_alloc_bytes
// 	- go_memstats_heap_alloc_bytes
// 	- go_memstats_mallocs_total
// 	- process_cpu_seconds_total
// 	- go_goroutines
// 	- etcd_storage_watcher_total
// 	- etcd_storage_watch_stream_total
// 	- etcd_storage_slow_watcher_total
//
func GetMetrics(endpoints ...string) (map[string]map[string]float64, map[string]string, error) {
	_, nameToEndpoint, err := GetStats(endpoints...)
	if err != nil {
		return nil, nil, err
	}

	rm := make(map[string]map[string]float64)
	for _, endpoint := range endpoints {
		resp, err := http.Get(endpoint + "/metrics")
		if err != nil {
			rm[endpoint] = defaultMetrics
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
			rm[endpoint] = defaultMetrics
			continue
		}
		resp.Body.Close()

		rm[endpoint] = mm
	}

	return rm, nameToEndpoint, nil
}

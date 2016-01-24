package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gophergala2016/runetcd/etcdproc"
	"golang.org/x/net/context"
)

func endpointHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			return err
		}
		cm := ""
		if len(req.Form["endpoint_name"]) != 0 {
			cm = req.Form["endpoint_name"][0]
		}
		if err := req.ParseForm(); err != nil {
			return err
		}
		cm = strings.Replace(cm, " ", "", -1)
		emap := make(map[string]struct{})
		for _, ep := range strings.Split(cm, "\n") {
			ep = cleanUp(ep)
			if len(ep) != 0 {
				u, err := url.Parse(ep)
				if err == nil {
					if strings.HasPrefix(ep, "http://") || strings.HasPrefix(ep, "https://") {
						emap[u.String()] = struct{}{}
					}
				}
			}
		}
		endpoints := []string{}
		for k := range emap {
			endpoints = append(endpoints, k)
		}
		sort.Strings(endpoints)
		if len(endpoints) > 5 {
			endpoints = endpoints[:5:5]
		}
		globalCache.mu.Lock()
		globalCache.endpoints = endpoints
		globalCache.mu.Unlock()

	case "GET":
		globalCache.mu.Lock()
		endpoints := globalCache.endpoints
		globalCache.mu.Unlock()

		success := len(endpoints) != 0
		msg := fmt.Sprintf(`&nbsp;&nbsp;<b><font color="blue">Success</font><b>: %s`, strings.Join(endpoints, ", "))
		if !success {
			msg = `&nbsp;&nbsp;<b><font color="red">Fail</font></b>: no valid endpoints`
		}
		resp := struct {
			RefreshInMillisecond int
			Success              bool
			Message              string
		}{
			int(cmdFlag.RefreshRate.Seconds()) * 1000,
			success,
			msg,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			return err
		}

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

func statsHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case "GET":
		globalCache.mu.Lock()
		endpoints := globalCache.endpoints
		globalCache.mu.Unlock()

		endpointToStats, nameToEndpoint, err := etcdproc.GetStats(endpoints...)
		if err != nil {
			return err
		}

		names := []string{}
		for k := range nameToEndpoint {
			names = append(names, k)
		}
		sort.Strings(names)
		copied := make([]string, 5)
		copy(copied, names)
		names = copied

		name1, endpoint1 := names[0], nameToEndpoint[names[0]]
		name2, endpoint2 := names[1], nameToEndpoint[names[1]]
		name3, endpoint3 := names[2], nameToEndpoint[names[2]]
		name4, endpoint4 := names[3], nameToEndpoint[names[3]]
		name5, endpoint5 := names[4], nameToEndpoint[names[4]]
		etcd1ID := ""
		etcd1State := ""
		if v, ok := endpointToStats[endpoint1]; ok {
			etcd1ID = v.ID
			etcd1State = v.State
		}
		etcd2ID := ""
		etcd2State := ""
		if v, ok := endpointToStats[endpoint2]; ok {
			etcd2ID = v.ID
			etcd2State = v.State
		}
		etcd3ID := ""
		etcd3State := ""
		if v, ok := endpointToStats[endpoint3]; ok {
			etcd3ID = v.ID
			etcd3State = v.State
		}
		etcd4ID := ""
		etcd4State := ""
		if v, ok := endpointToStats[endpoint4]; ok {
			etcd4ID = v.ID
			etcd4State = v.State
		}
		etcd5ID := ""
		etcd5State := ""
		if v, ok := endpointToStats[endpoint5]; ok {
			etcd5ID = v.ID
			etcd5State = v.State
		}
		resp := struct {
			Etcd1Endpoint string
			Etcd1Name     string
			Etcd1ID       string
			Etcd1State    string

			Etcd2Endpoint string
			Etcd2Name     string
			Etcd2ID       string
			Etcd2State    string

			Etcd3Endpoint string
			Etcd3Name     string
			Etcd3ID       string
			Etcd3State    string

			Etcd4Endpoint string
			Etcd4Name     string
			Etcd4ID       string
			Etcd4State    string

			Etcd5Endpoint string
			Etcd5Name     string
			Etcd5ID       string
			Etcd5State    string
		}{
			endpoint1,
			name1,
			etcd1ID,
			etcd1State,

			endpoint2,
			name2,
			etcd2ID,
			etcd2State,

			endpoint3,
			name3,
			etcd3ID,
			etcd3State,

			endpoint4,
			name4,
			etcd4ID,
			etcd4State,

			endpoint5,
			name5,
			etcd5ID,
			etcd5State,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			return err
		}

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

func metricsHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case "GET":
		globalCache.mu.Lock()
		endpoints := globalCache.endpoints
		globalCache.mu.Unlock()

		endpointToMetrics, nameToEndpoint, err := etcdproc.GetMetrics(endpoints...)
		if err != nil {
			return err
		}

		names := []string{}
		for k := range nameToEndpoint {
			names = append(names, k)
		}
		sort.Strings(names)
		copied := make([]string, 5)
		copy(copied, names)
		names = copied

		name1, endpoint1 := names[0], nameToEndpoint[names[0]]
		name2, endpoint2 := names[1], nameToEndpoint[names[1]]
		name3, endpoint3 := names[2], nameToEndpoint[names[2]]
		name4, endpoint4 := names[3], nameToEndpoint[names[3]]
		name5, endpoint5 := names[4], nameToEndpoint[names[4]]
		etcd1StorageKeysTotal := 0.0
		etcd1StorageBytes := 0.0
		etcd1StorageBytesStr := "0 bytes"
		etcd1WalFsyncSecondsSum := 0.0
		etcd1GcSecondsSum := 0.0
		etcd1MemstatsAllocBytes := 0.0
		etcd1MemstatsAllocBytesStr := "0 bytes"
		etcd1MemstatsHeapAllocBytes := 0.0
		etcd1MemstatsHeapAllocBytesStr := "0 bytes"
		etcd1MemstatsMallocsTotal := 0.0
		etcd1ProcessCPUSeconds := 0.0
		etcd1Goroutines := 0.0
		etcd1StorageWatcherTotal := 0.0
		etcd1StorageWatchStreamTotal := 0.0
		etcd1StorageSlowWatcherTotal := 0.0
		if vm, ok := endpointToMetrics[endpoint1]; ok {
			etcd1StorageKeysTotal = vm["etcd_storage_keys_total"]
			etcd1StorageBytes = vm["etcd_storage_db_total_size_in_bytes"]
			etcd1StorageBytesStr = humanize.Bytes(uint64(vm["etcd_storage_db_total_size_in_bytes"]))
			etcd1WalFsyncSecondsSum = vm["etcd_wal_fsync_durations_seconds_sum"]
			etcd1GcSecondsSum = vm["go_gc_duration_seconds_sum"]
			etcd1MemstatsAllocBytes = vm["go_memstats_alloc_bytes"]
			etcd1MemstatsAllocBytesStr = humanize.Bytes(uint64(vm["go_memstats_alloc_bytes"]))
			etcd1MemstatsHeapAllocBytes = vm["go_memstats_heap_alloc_bytes"]
			etcd1MemstatsHeapAllocBytesStr = humanize.Bytes(uint64(vm["go_memstats_heap_alloc_bytes"]))
			etcd1MemstatsMallocsTotal = vm["go_memstats_mallocs_total"]
			etcd1ProcessCPUSeconds = vm["process_cpu_seconds_total"]
			etcd1Goroutines = vm["go_goroutines"]
			etcd1StorageWatcherTotal = vm["etcd_storage_watcher_total"]
			etcd1StorageWatchStreamTotal = vm["etcd_storage_watch_stream_total"]
			etcd1StorageSlowWatcherTotal = vm["etcd_storage_slow_watcher_total"]
		}
		etcd2StorageKeysTotal := 0.0
		etcd2StorageBytes := 0.0
		etcd2StorageBytesStr := "0 bytes"
		etcd2WalFsyncSecondsSum := 0.0
		etcd2GcSecondsSum := 0.0
		etcd2MemstatsAllocBytes := 0.0
		etcd2MemstatsAllocBytesStr := "0 bytes"
		etcd2MemstatsHeapAllocBytes := 0.0
		etcd2MemstatsHeapAllocBytesStr := "0 bytes"
		etcd2MemstatsMallocsTotal := 0.0
		etcd2ProcessCPUSeconds := 0.0
		etcd2Goroutines := 0.0
		etcd2StorageWatcherTotal := 0.0
		etcd2StorageWatchStreamTotal := 0.0
		etcd2StorageSlowWatcherTotal := 0.0
		if vm, ok := endpointToMetrics[endpoint2]; ok {
			etcd2StorageKeysTotal = vm["etcd_storage_keys_total"]
			etcd2StorageBytes = vm["etcd_storage_db_total_size_in_bytes"]
			etcd2StorageBytesStr = humanize.Bytes(uint64(vm["etcd_storage_db_total_size_in_bytes"]))
			etcd2WalFsyncSecondsSum = vm["etcd_wal_fsync_durations_seconds_sum"]
			etcd2GcSecondsSum = vm["go_gc_duration_seconds_sum"]
			etcd2MemstatsAllocBytes = vm["go_memstats_alloc_bytes"]
			etcd2MemstatsAllocBytesStr = humanize.Bytes(uint64(vm["go_memstats_alloc_bytes"]))
			etcd2MemstatsHeapAllocBytes = vm["go_memstats_heap_alloc_bytes"]
			etcd2MemstatsHeapAllocBytesStr = humanize.Bytes(uint64(vm["go_memstats_heap_alloc_bytes"]))
			etcd2MemstatsMallocsTotal = vm["go_memstats_mallocs_total"]
			etcd2ProcessCPUSeconds = vm["process_cpu_seconds_total"]
			etcd2Goroutines = vm["go_goroutines"]
			etcd2StorageWatcherTotal = vm["etcd_storage_watcher_total"]
			etcd2StorageWatchStreamTotal = vm["etcd_storage_watch_stream_total"]
			etcd2StorageSlowWatcherTotal = vm["etcd_storage_slow_watcher_total"]
		}
		etcd3StorageKeysTotal := 0.0
		etcd3StorageBytes := 0.0
		etcd3StorageBytesStr := "0 bytes"
		etcd3WalFsyncSecondsSum := 0.0
		etcd3GcSecondsSum := 0.0
		etcd3MemstatsAllocBytes := 0.0
		etcd3MemstatsAllocBytesStr := "0 bytes"
		etcd3MemstatsHeapAllocBytes := 0.0
		etcd3MemstatsHeapAllocBytesStr := "0 bytes"
		etcd3MemstatsMallocsTotal := 0.0
		etcd3ProcessCPUSeconds := 0.0
		etcd3Goroutines := 0.0
		etcd3StorageWatcherTotal := 0.0
		etcd3StorageWatchStreamTotal := 0.0
		etcd3StorageSlowWatcherTotal := 0.0
		if vm, ok := endpointToMetrics[endpoint3]; ok {
			etcd3StorageKeysTotal = vm["etcd_storage_keys_total"]
			etcd3StorageBytes = vm["etcd_storage_db_total_size_in_bytes"]
			etcd3StorageBytesStr = humanize.Bytes(uint64(vm["etcd_storage_db_total_size_in_bytes"]))
			etcd3WalFsyncSecondsSum = vm["etcd_wal_fsync_durations_seconds_sum"]
			etcd3GcSecondsSum = vm["go_gc_duration_seconds_sum"]
			etcd3MemstatsAllocBytes = vm["go_memstats_alloc_bytes"]
			etcd3MemstatsAllocBytesStr = humanize.Bytes(uint64(vm["go_memstats_alloc_bytes"]))
			etcd3MemstatsHeapAllocBytes = vm["go_memstats_heap_alloc_bytes"]
			etcd3MemstatsHeapAllocBytesStr = humanize.Bytes(uint64(vm["go_memstats_heap_alloc_bytes"]))
			etcd3MemstatsMallocsTotal = vm["go_memstats_mallocs_total"]
			etcd3ProcessCPUSeconds = vm["process_cpu_seconds_total"]
			etcd3Goroutines = vm["go_goroutines"]
			etcd3StorageWatcherTotal = vm["etcd_storage_watcher_total"]
			etcd3StorageWatchStreamTotal = vm["etcd_storage_watch_stream_total"]
			etcd3StorageSlowWatcherTotal = vm["etcd_storage_slow_watcher_total"]
		}
		etcd4StorageKeysTotal := 0.0
		etcd4StorageBytes := 0.0
		etcd4StorageBytesStr := "0 bytes"
		etcd4WalFsyncSecondsSum := 0.0
		etcd4GcSecondsSum := 0.0
		etcd4MemstatsAllocBytes := 0.0
		etcd4MemstatsAllocBytesStr := "0 bytes"
		etcd4MemstatsHeapAllocBytes := 0.0
		etcd4MemstatsHeapAllocBytesStr := "0 bytes"
		etcd4MemstatsMallocsTotal := 0.0
		etcd4ProcessCPUSeconds := 0.0
		etcd4Goroutines := 0.0
		etcd4StorageWatcherTotal := 0.0
		etcd4StorageWatchStreamTotal := 0.0
		etcd4StorageSlowWatcherTotal := 0.0
		if vm, ok := endpointToMetrics[endpoint4]; ok {
			etcd4StorageKeysTotal = vm["etcd_storage_keys_total"]
			etcd4StorageBytes = vm["etcd_storage_db_total_size_in_bytes"]
			etcd4StorageBytesStr = humanize.Bytes(uint64(vm["etcd_storage_db_total_size_in_bytes"]))
			etcd4WalFsyncSecondsSum = vm["etcd_wal_fsync_durations_seconds_sum"]
			etcd4GcSecondsSum = vm["go_gc_duration_seconds_sum"]
			etcd4MemstatsAllocBytes = vm["go_memstats_alloc_bytes"]
			etcd4MemstatsAllocBytesStr = humanize.Bytes(uint64(vm["go_memstats_alloc_bytes"]))
			etcd4MemstatsHeapAllocBytes = vm["go_memstats_heap_alloc_bytes"]
			etcd4MemstatsHeapAllocBytesStr = humanize.Bytes(uint64(vm["go_memstats_heap_alloc_bytes"]))
			etcd4MemstatsMallocsTotal = vm["go_memstats_mallocs_total"]
			etcd4ProcessCPUSeconds = vm["process_cpu_seconds_total"]
			etcd4Goroutines = vm["go_goroutines"]
			etcd4StorageWatcherTotal = vm["etcd_storage_watcher_total"]
			etcd4StorageWatchStreamTotal = vm["etcd_storage_watch_stream_total"]
			etcd4StorageSlowWatcherTotal = vm["etcd_storage_slow_watcher_total"]
		}
		etcd5StorageKeysTotal := 0.0
		etcd5StorageBytes := 0.0
		etcd5StorageBytesStr := "0 bytes"
		etcd5WalFsyncSecondsSum := 0.0
		etcd5GcSecondsSum := 0.0
		etcd5MemstatsAllocBytes := 0.0
		etcd5MemstatsAllocBytesStr := "0 bytes"
		etcd5MemstatsHeapAllocBytes := 0.0
		etcd5MemstatsHeapAllocBytesStr := "0 bytes"
		etcd5MemstatsMallocsTotal := 0.0
		etcd5ProcessCPUSeconds := 0.0
		etcd5Goroutines := 0.0
		etcd5StorageWatcherTotal := 0.0
		etcd5StorageWatchStreamTotal := 0.0
		etcd5StorageSlowWatcherTotal := 0.0
		if vm, ok := endpointToMetrics[endpoint5]; ok {
			etcd5StorageKeysTotal = vm["etcd_storage_keys_total"]
			etcd5StorageBytes = vm["etcd_storage_db_total_size_in_bytes"]
			etcd5StorageBytesStr = humanize.Bytes(uint64(vm["etcd_storage_db_total_size_in_bytes"]))
			etcd5WalFsyncSecondsSum = vm["etcd_wal_fsync_durations_seconds_sum"]
			etcd5GcSecondsSum = vm["go_gc_duration_seconds_sum"]
			etcd5MemstatsAllocBytes = vm["go_memstats_alloc_bytes"]
			etcd5MemstatsAllocBytesStr = humanize.Bytes(uint64(vm["go_memstats_alloc_bytes"]))
			etcd5MemstatsHeapAllocBytes = vm["go_memstats_heap_alloc_bytes"]
			etcd5MemstatsHeapAllocBytesStr = humanize.Bytes(uint64(vm["go_memstats_heap_alloc_bytes"]))
			etcd5MemstatsMallocsTotal = vm["go_memstats_mallocs_total"]
			etcd5ProcessCPUSeconds = vm["process_cpu_seconds_total"]
			etcd5Goroutines = vm["go_goroutines"]
			etcd5StorageWatcherTotal = vm["etcd_storage_watcher_total"]
			etcd5StorageWatchStreamTotal = vm["etcd_storage_watch_stream_total"]
			etcd5StorageSlowWatcherTotal = vm["etcd_storage_slow_watcher_total"]
		}
		resp := struct {
			Etcd1Name                      string
			Etcd1Endpoint                  string
			Etcd1StorageKeysTotal          float64
			Etcd1StorageBytes              float64
			Etcd1StorageBytesStr           string
			Etcd1WalFsyncSecondsSum        float64
			Etcd1GcSecondsSum              float64
			Etcd1MemstatsAllocBytes        float64
			Etcd1MemstatsAllocBytesStr     string
			Etcd1MemstatsHeapAllocBytes    float64
			Etcd1MemstatsHeapAllocBytesStr string
			Etcd1MemstatsMallocsTotal      float64
			Etcd1ProcessCPUSeconds         float64
			Etcd1Goroutines                float64
			Etcd1StorageWatcherTotal       float64
			Etcd1StorageWatchStreamTotal   float64
			Etcd1StorageSlowWatcherTotal   float64

			Etcd2Name                      string
			Etcd2Endpoint                  string
			Etcd2StorageKeysTotal          float64
			Etcd2StorageBytes              float64
			Etcd2StorageBytesStr           string
			Etcd2WalFsyncSecondsSum        float64
			Etcd2GcSecondsSum              float64
			Etcd2MemstatsAllocBytes        float64
			Etcd2MemstatsAllocBytesStr     string
			Etcd2MemstatsHeapAllocBytes    float64
			Etcd2MemstatsHeapAllocBytesStr string
			Etcd2MemstatsMallocsTotal      float64
			Etcd2ProcessCPUSeconds         float64
			Etcd2Goroutines                float64
			Etcd2StorageWatcherTotal       float64
			Etcd2StorageWatchStreamTotal   float64
			Etcd2StorageSlowWatcherTotal   float64

			Etcd3Name                      string
			Etcd3Endpoint                  string
			Etcd3StorageKeysTotal          float64
			Etcd3StorageBytes              float64
			Etcd3StorageBytesStr           string
			Etcd3WalFsyncSecondsSum        float64
			Etcd3GcSecondsSum              float64
			Etcd3MemstatsAllocBytes        float64
			Etcd3MemstatsAllocBytesStr     string
			Etcd3MemstatsHeapAllocBytes    float64
			Etcd3MemstatsHeapAllocBytesStr string
			Etcd3MemstatsMallocsTotal      float64
			Etcd3ProcessCPUSeconds         float64
			Etcd3Goroutines                float64
			Etcd3StorageWatcherTotal       float64
			Etcd3StorageWatchStreamTotal   float64
			Etcd3StorageSlowWatcherTotal   float64

			Etcd4Name                      string
			Etcd4Endpoint                  string
			Etcd4StorageKeysTotal          float64
			Etcd4StorageBytes              float64
			Etcd4StorageBytesStr           string
			Etcd4WalFsyncSecondsSum        float64
			Etcd4GcSecondsSum              float64
			Etcd4MemstatsAllocBytes        float64
			Etcd4MemstatsAllocBytesStr     string
			Etcd4MemstatsHeapAllocBytes    float64
			Etcd4MemstatsHeapAllocBytesStr string
			Etcd4MemstatsMallocsTotal      float64
			Etcd4ProcessCPUSeconds         float64
			Etcd4Goroutines                float64
			Etcd4StorageWatcherTotal       float64
			Etcd4StorageWatchStreamTotal   float64
			Etcd4StorageSlowWatcherTotal   float64

			Etcd5Name                      string
			Etcd5Endpoint                  string
			Etcd5StorageKeysTotal          float64
			Etcd5StorageBytes              float64
			Etcd5StorageBytesStr           string
			Etcd5WalFsyncSecondsSum        float64
			Etcd5GcSecondsSum              float64
			Etcd5MemstatsAllocBytes        float64
			Etcd5MemstatsAllocBytesStr     string
			Etcd5MemstatsHeapAllocBytes    float64
			Etcd5MemstatsHeapAllocBytesStr string
			Etcd5MemstatsMallocsTotal      float64
			Etcd5ProcessCPUSeconds         float64
			Etcd5Goroutines                float64
			Etcd5StorageWatcherTotal       float64
			Etcd5StorageWatchStreamTotal   float64
			Etcd5StorageSlowWatcherTotal   float64
		}{
			name1,
			endpoint1,
			etcd1StorageKeysTotal,
			etcd1StorageBytes,
			etcd1StorageBytesStr,
			etcd1WalFsyncSecondsSum,
			etcd1GcSecondsSum,
			etcd1MemstatsAllocBytes,
			etcd1MemstatsAllocBytesStr,
			etcd1MemstatsHeapAllocBytes,
			etcd1MemstatsHeapAllocBytesStr,
			etcd1MemstatsMallocsTotal,
			etcd1ProcessCPUSeconds,
			etcd1Goroutines,
			etcd1StorageWatcherTotal,
			etcd1StorageWatchStreamTotal,
			etcd1StorageSlowWatcherTotal,

			name2,
			endpoint2,
			etcd2StorageKeysTotal,
			etcd2StorageBytes,
			etcd2StorageBytesStr,
			etcd2WalFsyncSecondsSum,
			etcd2GcSecondsSum,
			etcd2MemstatsAllocBytes,
			etcd2MemstatsAllocBytesStr,
			etcd2MemstatsHeapAllocBytes,
			etcd2MemstatsHeapAllocBytesStr,
			etcd2MemstatsMallocsTotal,
			etcd2ProcessCPUSeconds,
			etcd2Goroutines,
			etcd2StorageWatcherTotal,
			etcd2StorageWatchStreamTotal,
			etcd2StorageSlowWatcherTotal,

			name3,
			endpoint3,
			etcd3StorageKeysTotal,
			etcd3StorageBytes,
			etcd3StorageBytesStr,
			etcd3WalFsyncSecondsSum,
			etcd3GcSecondsSum,
			etcd3MemstatsAllocBytes,
			etcd3MemstatsAllocBytesStr,
			etcd3MemstatsHeapAllocBytes,
			etcd3MemstatsHeapAllocBytesStr,
			etcd3MemstatsMallocsTotal,
			etcd3ProcessCPUSeconds,
			etcd3Goroutines,
			etcd3StorageWatcherTotal,
			etcd3StorageWatchStreamTotal,
			etcd3StorageSlowWatcherTotal,

			name4,
			endpoint4,
			etcd4StorageKeysTotal,
			etcd4StorageBytes,
			etcd4StorageBytesStr,
			etcd4WalFsyncSecondsSum,
			etcd4GcSecondsSum,
			etcd4MemstatsAllocBytes,
			etcd4MemstatsAllocBytesStr,
			etcd4MemstatsHeapAllocBytes,
			etcd4MemstatsHeapAllocBytesStr,
			etcd4MemstatsMallocsTotal,
			etcd4ProcessCPUSeconds,
			etcd4Goroutines,
			etcd4StorageWatcherTotal,
			etcd4StorageWatchStreamTotal,
			etcd4StorageSlowWatcherTotal,

			name5,
			endpoint5,
			etcd5StorageKeysTotal,
			etcd5StorageBytes,
			etcd5StorageBytesStr,
			etcd5WalFsyncSecondsSum,
			etcd5GcSecondsSum,
			etcd5MemstatsAllocBytes,
			etcd5MemstatsAllocBytesStr,
			etcd5MemstatsHeapAllocBytes,
			etcd5MemstatsHeapAllocBytesStr,
			etcd5MemstatsMallocsTotal,
			etcd5ProcessCPUSeconds,
			etcd5Goroutines,
			etcd5StorageWatcherTotal,
			etcd5StorageWatchStreamTotal,
			etcd5StorageSlowWatcherTotal,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			return err
		}

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

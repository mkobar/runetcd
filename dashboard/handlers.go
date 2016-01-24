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
		if vm, ok := endpointToMetrics[endpoint1]; ok {
			etcd1StorageKeysTotal = vm["etcd_storage_keys_total"]
			etcd1StorageBytes = vm["etcd_storage_db_total_size_in_bytes"]
			etcd1StorageBytesStr = humanize.Bytes(uint64(vm["etcd_storage_db_total_size_in_bytes"]))
		}
		etcd2StorageKeysTotal := 0.0
		etcd2StorageBytes := 0.0
		etcd2StorageBytesStr := "0 bytes"
		if vm, ok := endpointToMetrics[endpoint2]; ok {
			etcd2StorageKeysTotal = vm["etcd_storage_keys_total"]
			etcd2StorageBytes = vm["etcd_storage_db_total_size_in_bytes"]
			etcd2StorageBytesStr = humanize.Bytes(uint64(vm["etcd_storage_db_total_size_in_bytes"]))
		}
		etcd3StorageKeysTotal := 0.0
		etcd3StorageBytes := 0.0
		etcd3StorageBytesStr := "0 bytes"
		if vm, ok := endpointToMetrics[endpoint3]; ok {
			etcd3StorageKeysTotal = vm["etcd_storage_keys_total"]
			etcd3StorageBytes = vm["etcd_storage_db_total_size_in_bytes"]
			etcd3StorageBytesStr = humanize.Bytes(uint64(vm["etcd_storage_db_total_size_in_bytes"]))
		}
		etcd4StorageKeysTotal := 0.0
		etcd4StorageBytes := 0.0
		etcd4StorageBytesStr := "0 bytes"
		if vm, ok := endpointToMetrics[endpoint4]; ok {
			etcd4StorageKeysTotal = vm["etcd_storage_keys_total"]
			etcd4StorageBytes = vm["etcd_storage_db_total_size_in_bytes"]
			etcd4StorageBytesStr = humanize.Bytes(uint64(vm["etcd_storage_db_total_size_in_bytes"]))
		}
		etcd5StorageKeysTotal := 0.0
		etcd5StorageBytes := 0.0
		etcd5StorageBytesStr := "0 bytes"
		if vm, ok := endpointToMetrics[endpoint5]; ok {
			etcd5StorageKeysTotal = vm["etcd_storage_keys_total"]
			etcd5StorageBytes = vm["etcd_storage_db_total_size_in_bytes"]
			etcd5StorageBytesStr = humanize.Bytes(uint64(vm["etcd_storage_db_total_size_in_bytes"]))
		}
		resp := struct {
			Etcd1Name             string
			Etcd1StorageKeysTotal float64
			Etcd1StorageBytes     float64
			Etcd1StorageBytesStr  string

			Etcd2Name             string
			Etcd2StorageKeysTotal float64
			Etcd2StorageBytes     float64
			Etcd2StorageBytesStr  string

			Etcd3Name             string
			Etcd3StorageKeysTotal float64
			Etcd3StorageBytes     float64
			Etcd3StorageBytesStr  string

			Etcd4Name             string
			Etcd4StorageKeysTotal float64
			Etcd4StorageBytes     float64
			Etcd4StorageBytesStr  string

			Etcd5Name             string
			Etcd5StorageKeysTotal float64
			Etcd5StorageBytes     float64
			Etcd5StorageBytesStr  string
		}{
			name1,
			etcd1StorageKeysTotal,
			etcd1StorageBytes,
			etcd1StorageBytesStr,

			name2,
			etcd2StorageKeysTotal,
			etcd2StorageBytes,
			etcd2StorageBytesStr,

			name3,
			etcd3StorageKeysTotal,
			etcd3StorageBytes,
			etcd3StorageBytesStr,

			name4,
			etcd4StorageKeysTotal,
			etcd4StorageBytes,
			etcd4StorageBytesStr,

			name5,
			etcd5StorageKeysTotal,
			etcd5StorageBytes,
			etcd5StorageBytesStr,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			return err
		}

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

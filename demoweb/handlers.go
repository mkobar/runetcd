package demoweb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gophergala2016/runetcd/etcdproc"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

func wsHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	user := ctx.Value(userKey).(*string)
	userID := *user
	globalCache.mu.Lock()
	upgrader := globalCache.perUserID[userID].upgrader
	globalCache.mu.Unlock()

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		globalCache.perUserID[userID].donec <- struct{}{}
		return err
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			globalCache.perUserID[userID].donec <- struct{}{}
			return err
		}
		if err := c.WriteMessage(mt, message); err != nil {
			globalCache.perUserID[userID].donec <- struct{}{}
			return err
		}
	}
}

func streamHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	user := ctx.Value(userKey).(*string)
	userID := *user

	switch req.Method {
	case "GET":
		// no need Lock because it's channel
		//
		// globalCache.mu.Lock()
		// globalCache.mu.Unlock()
		streams := []string{}

	escape:
		for {
			select {
			case s, ok := <-globalCache.perUserID[userID].bufStream:
				if !ok {
					break // when the channel is closed
				}
				streams = append(streams, s)
			case <-time.After(time.Second):
				break escape
			}
		}

		if len(streams) > 0 {
			// When used with new EventSource('/stream') in Javascript
			//
			// w.Header().Set("Content-Type", "text/event-stream")
			// fmt.Fprintf(w, fmt.Sprintf("id: %s\nevent: %s\ndata: %s\n\n", userID, "stream_log", strings.Join(streams, "\n")))

			// When used with setInterval in Javascript
			//
			// fmt.Fprintln(w, strings.Join(streams, "<br>"))
			//
			resp := struct {
				Logs string
				Size int
			}{
				strings.Join(streams, "<br>"),
				len(streams),
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				return err
			}
		}

		if f, ok := w.(http.Flusher); ok {
			if f != nil {
				f.Flush()
			}
		}

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

func startClusterHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	user := ctx.Value(userKey).(*string)
	userID := *user

	switch req.Method {
	case "GET":
		now := time.Now()
		toRun := false
		var sub time.Duration
		globalCache.mu.Lock()
		if globalCache.perUserID[userID].clusterStarted.IsZero() {
			globalCache.perUserID[userID].clusterStarted = now
			toRun = true
		} else {
			sub = now.Sub(globalCache.perUserID[userID].clusterStarted)
			if sub > startClusterMinInterval {
				globalCache.perUserID[userID].clusterStarted = now
				toRun = true
			}
		}
		globalCache.mu.Unlock()
		if !toRun {
			// globalCache.perUserID[userID].bufStream <- boldHTMLMsg(fmt.Sprintf("Limit exceeded (%v since last)! Please wait a few minutes or run locally!", sub))
			fmt.Fprintln(w, boldHTMLMsg(fmt.Sprintf("Limit exceeded (%v since last)! Please wait %v or run locally!", sub, startClusterMinInterval)))
			return nil
		}
		go spawnCluster(userID)
		fmt.Fprintln(w, boldHTMLMsg("Start cluster successfully requested!!!"))

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

func spawnCluster(userID string) {
	defer func() {
		if err := recover(); err != nil {
			globalCache.perUserID[userID].bufStream <- fmt.Sprintf("[cluster - panic] %+v", err)
			return
		}
	}()

	portPrefix := atomic.LoadInt32(&portStart)
	fs := make([]*etcdproc.Flags, 3)
	for i := range fs {
		df, err := etcdproc.NewFlags(fmt.Sprintf("etcd%d", i+1), globalPorts, int(portPrefix)+i, "etcd-cluster-token", "new", uuid.NewV4().String(), false, false, "", "", "")
		if err != nil {
			globalCache.perUserID[userID].bufStream <- fmt.Sprintf("[cluster - error] %+v", err)
			return
		}
		fs[i] = df
	}
	atomic.AddInt32(&portStart, 7)

	cs, err := etcdproc.CreateCluster(os.Stdout, globalCache.perUserID[userID].bufStream, etcdproc.ToHTML, cmdFlag.EtcdBinary, fs...)
	if err != nil {
		globalCache.perUserID[userID].bufStream <- fmt.Sprintf("[cluster - error] %+v", err)
		return
	}
	globalCache.mu.Lock()
	globalCache.perUserID[userID].cluster = cs
	globalCache.mu.Unlock()

	// this does not run with the program exits with os.Exit(0)
	defer func() {
		cs.RemoveAllDataDirs()
		globalCache.mu.Lock()
		globalCache.perUserID[userID].cluster = nil
		globalCache.mu.Unlock()
	}()

	errChan, done := make(chan error), make(chan struct{})
	go func() {
		globalCache.perUserID[userID].bufStream <- boldHTMLMsg("Starting all of those 3 nodes in default cluster group")
		if err := cs.StartAll(); err != nil {
			errChan <- err
			return
		}
		done <- struct{}{}
	}()

	select {
	case err := <-errChan:
		globalCache.perUserID[userID].bufStream <- fmt.Sprintf("[cluster - error] %+v", err)
	case <-done:
		globalCache.perUserID[userID].bufStream <- boldHTMLMsg("Cluster done!")
	case <-globalCache.perUserID[userID].donec:
		globalCache.perUserID[userID].bufStream <- boldHTMLMsg("Cluster done!")
	case <-time.After(cmdFlag.Timeout):
		globalCache.perUserID[userID].bufStream <- boldHTMLMsg("Cluster time out!")
	}
}

func startStressHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	user := ctx.Value(userKey).(*string)
	userID := *user

	switch req.Method {
	case "GET":

		toRun := false
		globalCache.mu.Lock()
		toRun = globalCache.perUserID[userID].cluster != nil
		globalCache.mu.Unlock()

		if !toRun {
			fmt.Fprintln(w, boldHTMLMsg("Cluster is not ready to receive requests!!!"))
			return nil
		}

		globalCache.mu.Lock()
		cs := globalCache.perUserID[userID].cluster
		globalCache.mu.Unlock()

		// TODO: print out readable k-v
		if err := cs.SimpleStress(os.Stdout, etcdproc.ToHTML, nameToStress); err != nil {
			fmt.Fprintln(w, boldHTMLMsg(fmt.Sprintf("exiting with: %v", err)))
			return err
		}

		fmt.Fprintln(w, boldHTMLMsg("Stress cluster successfully requested!!!"))

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

func statsHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	user := ctx.Value(userKey).(*string)
	userID := *user

	switch req.Method {
	case "GET":
		toRun := false
		globalCache.mu.Lock()
		toRun = globalCache.perUserID[userID].cluster != nil
		globalCache.mu.Unlock()

		if !toRun {
			globalCache.perUserID[userID].bufStream <- boldHTMLMsg("Cluster is not ready to provide stats!!!")
			return nil
		}

		globalCache.mu.Lock()
		cs := globalCache.perUserID[userID].cluster
		nameToStats, err := cs.GetStats()
		globalCache.mu.Unlock()
		if err != nil {
			globalCache.perUserID[userID].bufStream <- boldHTMLMsg(fmt.Sprintf("exiting with: %v", err))
			return err
		}

		names := []string{}
		for n := range nameToStats {
			names = append(names, n)
		}
		sort.Strings(names)
		if len(names) != 3 {
			return fmt.Errorf("expected 3 nodes but got %d nodes", len(names))
		}

		resp := struct {
			Etcd1Name  string
			Etcd1ID    string
			Etcd1State string

			Etcd2Name  string
			Etcd2ID    string
			Etcd2State string

			Etcd3Name  string
			Etcd3ID    string
			Etcd3State string
		}{
			names[0],
			nameToStats[names[0]].ID,
			nameToStats[names[0]].State,

			names[1],
			nameToStats[names[1]].ID,
			nameToStats[names[1]].State,

			names[2],
			nameToStats[names[2]].ID,
			nameToStats[names[2]].State,
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
	user := ctx.Value(userKey).(*string)
	userID := *user

	emptyResp := struct {
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
	}{
		"etcd1",
		0.0,
		0.0,
		"0 bytes",

		"etcd2",
		0.0,
		0.0,
		"0 bytes",

		"etcd3",
		0.0,
		0.0,
		"0 bytes",
	}

	switch req.Method {
	case "GET":
		toRun := false
		globalCache.mu.Lock()
		toRun = globalCache.perUserID[userID].cluster != nil
		globalCache.mu.Unlock()

		if !toRun {
			if err := json.NewEncoder(w).Encode(emptyResp); err != nil {
				return err
			}
			return nil
		}

		globalCache.mu.Lock()
		cs := globalCache.perUserID[userID].cluster
		nameToMetrics, err := cs.GetMetrics()
		globalCache.mu.Unlock()
		if err != nil {
			globalCache.perUserID[userID].bufStream <- boldHTMLMsg(fmt.Sprintf("exiting with: %v", err))
			if errJ := json.NewEncoder(w).Encode(emptyResp); errJ != nil {
				return errJ
			}
			return err
		}

		names := []string{}
		for n := range nameToMetrics {
			names = append(names, n)
		}
		sort.Strings(names)
		if len(names) != 3 {
			return fmt.Errorf("expected 3 nodes but got %d nodes", len(names))
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
		}{
			names[0],
			nameToMetrics[names[0]]["etcd_storage_keys_total"],
			nameToMetrics[names[0]]["etcd_storage_db_total_size_in_bytes"],
			humanize.Bytes(uint64(nameToMetrics[names[0]]["etcd_storage_db_total_size_in_bytes"])),

			names[1],
			nameToMetrics[names[1]]["etcd_storage_keys_total"],
			nameToMetrics[names[1]]["etcd_storage_db_total_size_in_bytes"],
			humanize.Bytes(uint64(nameToMetrics[names[1]]["etcd_storage_db_total_size_in_bytes"])),

			names[2],
			nameToMetrics[names[2]]["etcd_storage_keys_total"],
			nameToMetrics[names[2]]["etcd_storage_db_total_size_in_bytes"],
			humanize.Bytes(uint64(nameToMetrics[names[2]]["etcd_storage_db_total_size_in_bytes"])),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			return err
		}

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

func listCtlHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	user := ctx.Value(userKey).(*string)
	userID := *user

	switch req.Method {
	case "GET":
		globalCache.mu.Lock()
		ctlHistory := globalCache.perUserID[userID].ctlHistory
		globalCache.mu.Unlock()

		resp := struct {
			Values []string
		}{
			ctlHistory,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			return err
		}

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

func ctlHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	user := ctx.Value(userKey).(*string)
	userID := *user

	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			return err
		}
		cm := ""
		if len(req.Form["ctl_name"]) != 0 {
			cm = req.Form["ctl_name"][0]
		}
		globalCache.mu.Lock()
		globalCache.perUserID[userID].ctlCmd = cm
		globalCache.perUserID[userID].ctlHistory = append(globalCache.perUserID[userID].ctlHistory, cm)
		globalCache.mu.Unlock()

	case "GET":

		fmt.Println("ctl get")

		toRun := false
		globalCache.mu.Lock()
		toRun = globalCache.perUserID[userID].cluster != nil
		ctlCmd := globalCache.perUserID[userID].ctlCmd
		globalCache.mu.Unlock()

		if !toRun {
			fmt.Fprintln(w, boldHTMLMsg("Cluster is not ready to receive requests!!!"))
			return nil
		}
		if ctlCmd == "" {
			fmt.Fprintln(w, boldHTMLMsg("Invalid command received!!!"))
			return nil
		}

		// TODO: run command here and return the results
		fmt.Fprintln(w, "<br><b>"+ctlCmd+"</b><br>[result] a b<br><br>")

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

func killHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	user := ctx.Value(userKey).(*string)
	userID := *user

	switch req.Method {
	case "GET":

		toRun := false
		globalCache.mu.Lock()
		toRun = globalCache.perUserID[userID].cluster != nil
		globalCache.mu.Unlock()

		if !toRun {
			fmt.Fprintln(w, boldHTMLMsg("Cluster is not ready to receive requests!!! "+urlToName(req.URL.String())))
			return nil
		}

		globalCache.mu.Lock()
		cs := globalCache.perUserID[userID].cluster
		globalCache.mu.Unlock()

		if err := cs.Terminate(urlToName(req.URL.String())); err != nil {
			fmt.Fprintln(w, boldHTMLMsg("Terminate error "+urlToName(req.URL.String())))
			return err
		}
		fmt.Fprintln(w, boldHTMLMsg("Kill successfully requested for "+urlToName(req.URL.String())))

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

func restartHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	user := ctx.Value(userKey).(*string)
	userID := *user

	switch req.Method {
	case "GET":

		toRun := false
		globalCache.mu.Lock()
		toRun = globalCache.perUserID[userID].cluster != nil
		globalCache.mu.Unlock()

		if !toRun {
			fmt.Fprintln(w, boldHTMLMsg("Cluster is not ready to receive requests!!! "+urlToName(req.URL.String())))
			return nil
		}

		globalCache.mu.Lock()
		cs := globalCache.perUserID[userID].cluster
		globalCache.mu.Unlock()

		if err := cs.Restart(urlToName(req.URL.String())); err != nil {
			fmt.Fprintln(w, boldHTMLMsg("Restart error "+urlToName(req.URL.String())))
			return err
		}
		fmt.Fprintln(w, boldHTMLMsg("Restart successfully requested for "+urlToName(req.URL.String())))

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

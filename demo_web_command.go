package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gophergala2016/runetcd/run"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type (
	key int

	cache struct {
		mu        sync.Mutex
		perUserID map[string]*userData
	}
	userData struct {
		upgrader *websocket.Upgrader

		clusterStarted time.Time
		cluster        *run.Cluster
		donec          chan struct{}

		bufStream chan string

		ctlCmd     string
		ctlHistory []string
	}
)

const (
	demoWebPort     = ":8000"
	userKey     key = 0
)

var (
	globalCache             cache
	portStart               int32 = 11
	startClusterMinInterval       = 15 * time.Minute
)

func wsHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	user := ctx.Value(userKey).(*string)
	userID := *user
	globalCache.mu.Lock()
	upgrader := globalCache.perUserID[userID].upgrader
	globalCache.mu.Unlock()

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return err
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			return err
		}
		if err := c.WriteMessage(mt, message); err != nil {
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

func demoWebCommandFunc(cmd *cobra.Command, args []string) {
	rootContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/", http.FileServer(http.Dir("./static")))
	mainRouter.Handle("/ws", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(wsHandler)),
	})

	mainRouter.Handle("/stream", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(streamHandler)),
	})

	mainRouter.Handle("/start_cluster", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(startClusterHandler)),
	})
	mainRouter.Handle("/start_stress", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(startStressHandler)),
	})
	mainRouter.Handle("/stats", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(statsHandler)),
	})

	mainRouter.Handle("/list_ctl", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(listCtlHandler)),
	})
	mainRouter.Handle("/ctl", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(ctlHandler)),
	})

	mainRouter.Handle("/kill_1", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(killHandler)),
	})
	mainRouter.Handle("/kill_2", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(killHandler)),
	})
	mainRouter.Handle("/kill_3", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(killHandler)),
	})

	mainRouter.Handle("/restart_1", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(restartHandler)),
	})
	mainRouter.Handle("/restart_2", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(restartHandler)),
	})
	mainRouter.Handle("/restart_3", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(restartHandler)),
	})

	fmt.Fprintln(os.Stdout, "Serving http://localhost"+demoWebPort)
	if err := http.ListenAndServe(demoWebPort, mainRouter); err != nil {
		fmt.Fprintln(os.Stdout, "[runDemoWeb - error]", err)
		os.Exit(0)
	}
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
			// globalCache.perUserID[userID].bufStream <- boldHTMLMsg(fmt.Sprintf("Limit exceeded (%v since last): Please wait a few minutes or run locally", sub))
			fmt.Fprintln(w, boldHTMLMsg(fmt.Sprintf("Limit exceeded (%v since last): Please wait a few minutes or run locally", sub)))
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
		globalCache.mu.Lock()
		globalCache.perUserID[userID].cluster = nil
		globalCache.mu.Unlock()
		if err := recover(); err != nil {
			globalCache.perUserID[userID].bufStream <- fmt.Sprintf("[cluster - panic] %+v", err)
			return
		}
	}()

	portPrefix := atomic.LoadInt32(&portStart)
	fs := make([]*run.Flags, globalFlag.ClusterSize)
	for i := range fs {
		df, err := run.NewFlags(fmt.Sprintf("etcd%d", i+1), globalPorts, int(portPrefix)+i, "etcd-cluster-token", "new", uuid.NewV4().String(), false, false, "", "", "")
		if err != nil {
			globalCache.perUserID[userID].bufStream <- fmt.Sprintf("[cluster - error] %+v", err)
			return
		}
		fs[i] = df
	}
	atomic.AddInt32(&portStart, 7)

	cs, err := run.CreateCluster(os.Stdout, globalCache.perUserID[userID].bufStream, run.ToHTML, globalFlag.EtcdBinary, fs...)
	if err != nil {
		globalCache.perUserID[userID].bufStream <- fmt.Sprintf("[cluster - error] %+v", err)
		return
	}
	globalCache.mu.Lock()
	globalCache.perUserID[userID].cluster = cs
	globalCache.mu.Unlock()

	// this does not run with the program exits with os.Exit(0)
	defer cs.RemoveAllDataDirs()

	globalCache.perUserID[userID].bufStream <- boldHTMLMsg("Starting all of those 3 members in default cluster group")
	if err := cs.StartAll(); err != nil {
		globalCache.perUserID[userID].bufStream <- fmt.Sprintf("[cluster - error] %+v", err)
		return
	}

	select {
	case <-globalCache.perUserID[userID].donec:
		globalCache.perUserID[userID].bufStream <- boldHTMLMsg("Cluster done!")
		return
	case <-time.After(globalFlag.Timeout):
		globalCache.perUserID[userID].bufStream <- boldHTMLMsg("Cluster time out!")
		return
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
		if err := cs.Stress(w, nameToStress, globalFlag.DemoWebConnectionNumber, globalFlag.DemoWebClientNumber, globalFlag.DemoWebStressNumber, stressKeyN, stressValN); err != nil {
			fmt.Printf("exiting with: %+v\n", err)
			return err
		}

		fmt.Fprintln(w, boldHTMLMsg("Stress cluster successfully requested!!!"))

	default:
		http.Error(w, "Method Not Allowed", 405)
	}

	return nil
}

// TODO: graph this
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
		nameToStats, leaderNames, err := cs.GetStats()
		var errMsg error
		if err != nil {
			errMsg = err
		}
		nameToMetrics, err := cs.GetMetrics()
		if err != nil {
			errMsg = err
		}
		globalCache.mu.Unlock()
		if errMsg != nil {
			globalCache.perUserID[userID].bufStream <- boldHTMLMsg(fmt.Sprintf("exiting with: %v", errMsg))
			return errMsg
		}

		names := []string{}
		for n := range nameToStats {
			names = append(names, n)
		}
		sort.Strings(names)
		if len(names) != 3 {
			return fmt.Errorf("expected 3 members but got %d members", len(names))
		}

		resp := struct {
			LeaderNames []string

			Etcd1Name             string
			Etcd1ID               string
			Etcd1State            string
			Etcd1StorageKeysTotal float64
			Etcd1StorageBytes     float64
			Etcd1StorageBytesStr  string

			Etcd2Name             string
			Etcd2ID               string
			Etcd2State            string
			Etcd2StorageKeysTotal float64
			Etcd2StorageBytes     float64
			Etcd2StorageBytesStr  string

			Etcd3Name             string
			Etcd3ID               string
			Etcd3State            string
			Etcd3StorageKeysTotal float64
			Etcd3StorageBytes     float64
			Etcd3StorageBytesStr  string
		}{
			leaderNames,

			names[0],
			nameToStats[names[0]].ID,
			nameToStats[names[0]].State,
			nameToMetrics[names[0]]["etcd_storage_keys_total"],
			nameToMetrics[names[0]]["etcd_storage_db_total_size_in_bytes"],
			humanize.Bytes(uint64(nameToMetrics[names[0]]["etcd_storage_db_total_size_in_bytes"])),

			names[1],
			nameToStats[names[1]].ID,
			nameToStats[names[1]].State,
			nameToMetrics[names[1]]["etcd_storage_keys_total"],
			nameToMetrics[names[1]]["etcd_storage_db_total_size_in_bytes"],
			humanize.Bytes(uint64(nameToMetrics[names[1]]["etcd_storage_db_total_size_in_bytes"])),

			names[2],
			nameToStats[names[2]].ID,
			nameToStats[names[2]].State,
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

func urlToName(s string) string {
	ss := strings.Split(s, "_")
	suffix := ss[len(ss)-1]
	switch suffix {
	case "1":
		return "etcd1"
	case "2":
		return "etcd2"
	case "3":
		return "etcd3"
	default:
		return "unknown"
	}
}

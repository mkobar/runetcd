package demoweb

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/gophergala2016/runetcd/etcdproc"
	"github.com/gorilla/websocket"
	"github.com/gyuho/psn/ss"
	"github.com/spf13/cobra"
	"github.com/tylerb/graceful"
)

type Flag struct {
	EtcdBinary          string
	IntervalPortRefresh time.Duration
	Timeout             time.Duration
}

type (
	key int

	cache struct {
		mu        sync.Mutex
		perUserID map[string]*userData
	}
	userData struct {
		upgrader       *websocket.Upgrader
		clusterStarted time.Time

		cluster   *etcdproc.Cluster
		donec     chan struct{}
		bufStream chan string

		ctlCmd     string
		ctlHistory []string
	}
)

const (
	webPort     = ":8000"
	userKey key = 0
)

var (
	Command = &cobra.Command{
		Use:   "demo-web",
		Short: "demo-web demos etcd in a web browser.",
		Run:   CommandFunc,
	}

	cmdFlag     = Flag{}
	globalPorts = ss.NewPorts()

	globalCache             cache
	portStart               int32 = 11
	startClusterMinInterval       = 15 * time.Minute
)

func init() {
	cobra.EnablePrefixMatching = true
}

func init() {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetLevel(log.DebugLevel)

	// logInterval = 3 * time.Second
	// setLog := func() {
	// 	lp := strings.Replace(nowPacific().String()[:19], "-", "", -1)
	// 	lp = strings.Replace(lp, ":", "", -1)
	// 	lp = strings.Replace(lp, ".", "", -1)
	// 	lp = strings.Replace(lp, " ", "", -1)
	// 	fpath := "runetcd_" + lp + ".log"
	// 	f, err := openToAppend(fpath)
	// 	if err != nil {
	// 		fmt.Fprintln(os.Stdout, err)
	// 		return
	// 	}
	// 	log.SetOutput(f)
	// 	time.Sleep(logInterval)
	// 	f.Close()
	// }
	// setLog()
	// go func() {
	// 	for {
	// 		select {
	// 		case <-time.After(logInterval):
	// 			setLog()
	// 		}
	// 	}
	// }()
}

func init() {
	globalPorts.Refresh()
	go func() {
		for {
			select {
			case <-time.After(cmdFlag.IntervalPortRefresh):
				globalPorts.Refresh()
			}
		}
	}()
}

func init() {
	Command.PersistentFlags().StringVarP(&cmdFlag.EtcdBinary, "etcd-binary", "b", "bin/etcd", "Path of executatble etcd binary.")
	Command.PersistentFlags().DurationVar(&cmdFlag.IntervalPortRefresh, "port-refresh", 10*time.Second, "Interval to refresh free ports.")
	Command.PersistentFlags().DurationVar(&cmdFlag.Timeout, "timeout", 5*time.Minute, "After timeout, etcd shuts down itself.")
}

func CommandFunc(cmd *cobra.Command, args []string) {
	lp := strings.Replace(nowPacific().String()[:19], "-", "", -1)
	lp = strings.Replace(lp, ":", "", -1)
	lp = strings.Replace(lp, ".", "", -1)
	lp = strings.Replace(lp, " ", "", -1)

	usr, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
	}

	f, err := openToAppend(filepath.Join(usr.HomeDir, "runetcd_"+lp+".log"))
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		return
	}
	defer f.Close()
	log.SetOutput(f)

	rootContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mainRouter := http.NewServeMux()

	// mainRouter.Handle("/", http.FileServer(http.Dir("./demoweb_frontend")))
	mainRouter.Handle("/", &ContextAdapter{
		ctx:     rootContext,
		handler: ContextHandlerFunc(staticHandler),
	})

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
	mainRouter.Handle("/metrics", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(metricsHandler)),
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

	fmt.Fprintln(os.Stdout, "Serving http://localhost"+webPort)
	// if err := http.ListenAndServe(webPort, mainRouter); err != nil {
	// 	fmt.Fprintln(os.Stdout, "[runetcd demo-web error]", err)
	// 	os.Exit(0)
	// }
	// graceful.Run(webPort, 10*time.Second, withLogrus(mainRouter))
	graceful.Run(webPort, 10*time.Second, mainRouter)
}

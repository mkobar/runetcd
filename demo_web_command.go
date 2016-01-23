package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gophergala2016/runetcd/run"
	"github.com/gorilla/websocket"
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

func demoWebCommandFunc(cmd *cobra.Command, args []string) {
	rootContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/", http.FileServer(http.Dir("./static")))
	mainRouter.Handle("/ws", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(wsHandler)),
	})

	mainRouter.Handle("/start", &ContextAdapter{
		ctx:     rootContext,
		handler: withUserCache(ContextHandlerFunc(startClusterHandler)),
	})

	if err := http.ListenAndServe(demoWebPort, mainRouter); err != nil {
		fmt.Fprintln(os.Stdout, "[runDemoWeb - error]", err)
		os.Exit(0)
	}
}

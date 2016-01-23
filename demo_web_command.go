package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gophergala2016/runetcd/run"
	"github.com/gorilla/websocket"
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

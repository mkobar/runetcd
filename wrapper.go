package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
)

type ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request) error
}

// ContextHandlerFunc wraps func(context.Context, ResponseWriter, *Request)
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request) error

func (f ContextHandlerFunc) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	return f(ctx, w, req)
}

type ContextAdapter struct {
	ctx     context.Context
	handler ContextHandler
}

func (ca *ContextAdapter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := ca.handler.ServeHTTPContext(ca.ctx, w, req); err != nil {
		log.WithFields(log.Fields{
			"event_type": "error",
			"method":     req.Method,
			"path":       req.URL.Path,
			"error":      err,
		}).Errorln("ServeHTTP error")
	}
}

func withUserCache(h ContextHandler) ContextHandler {
	return ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, req *http.Request) error {

		userID := getUserID(req)
		ctx = context.WithValue(ctx, userKey, &userID)

		globalCache.mu.Lock()
		// (X) this will deadlock
		// defer globalCache.mu.Unlock()
		if globalCache.perUserID == nil {
			globalCache.perUserID = make(map[string]*userData)
		}
		if _, ok := globalCache.perUserID[userID]; !ok {
			globalCache.perUserID[userID] = &userData{
				upgrader: &websocket.Upgrader{},
				cluster:  nil,
				donec:    make(chan struct{}),

				bufStream: make(chan string, 5000),

				ctlHistory: []string{
					`etcdctlv3 put YOUR_KEY YOUR_VALUE`,
					`etcdctlv3 range YOUR_KEY`,
				},
			}
		}
		globalCache.mu.Unlock()

		// TODO: reset this periodically

		return h.ServeHTTPContext(ctx, w, req)
	})
}

package demoweb

import (
	"net/http"
	"time"

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

		// TODO: reset this periodically
		globalCache.mu.Lock()
		// (X) this will deadlock
		// defer globalCache.mu.Unlock()
		if globalCache.perUserID == nil {
			globalCache.perUserID = make(map[string]*userData)
		}
		if _, ok := globalCache.perUserID[userID]; !ok {
			globalCache.perUserID[userID] = &userData{
				upgrader:  &websocket.Upgrader{},
				cluster:   nil,
				donec:     make(chan struct{}, 10),
				bufStream: make(chan string, 5000),
				ctlHistory: []string{
					`etcdctlv3 put YOUR_KEY YOUR_VALUE`,
					`etcdctlv3 range YOUR_KEY`,
				},
			}
		}
		globalCache.mu.Unlock()

		return h.ServeHTTPContext(ctx, w, req)
	})
}

func withLogrus(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"event_type": "error",
					"referrer":   req.Referer(),
					"ua":         req.UserAgent(),
					"path":       req.URL.Path,
					"real_ip":    getRealIP(req),
					"error":      err,
				}).Errorln("withLogrus error")
			}
		}()

		start := nowPacific()
		h.ServeHTTP(w, req)
		took := time.Since(start)

		log.WithFields(log.Fields{
			"event_type": "ok",
			"referrer":   req.Referer(),
			"ua":         req.UserAgent(),
			"path":       req.URL.Path,
			"real_ip":    getRealIP(req),
		}).Debugf("took %s", took)
	}
}

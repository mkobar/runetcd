package dashboard

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type Flag struct {
	WebPort     string
	RefreshRate time.Duration
}

type (
	key int

	cache struct {
		mu        sync.Mutex
		endpoints []string
	}
)

const (
	userKey key = 0
)

var (
	Command = &cobra.Command{
		Use:   "dashboard",
		Short: "dashboard provides etcd dashboard in a web browser.",
		Run:   CommandFunc,
	}

	cmdFlag     = Flag{}
	globalCache cache
)

func init() {
	cobra.EnablePrefixMatching = true
}

func init() {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetLevel(log.DebugLevel)
}

func init() {
	Command.PersistentFlags().StringVarP(&cmdFlag.WebPort, "port", "p", ":8080", "Port to serve the dashboard.")
	Command.PersistentFlags().DurationVarP(&cmdFlag.RefreshRate, "refresh-rate", "r", 3*time.Second, "Refresh interval to get stats and metrics.")
}

func CommandFunc(cmd *cobra.Command, args []string) {
	rootContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/", http.FileServer(http.Dir("./dashboard_frontend")))

	mainRouter.Handle("/endpoint", &ContextAdapter{
		ctx:     rootContext,
		handler: ContextHandlerFunc(endpointHandler),
	})
	mainRouter.Handle("/stats", &ContextAdapter{
		ctx:     rootContext,
		handler: ContextHandlerFunc(statsHandler),
	})
	mainRouter.Handle("/metrics", &ContextAdapter{
		ctx:     rootContext,
		handler: ContextHandlerFunc(metricsHandler),
	})

	fmt.Fprintln(os.Stdout, "Serving http://localhost"+cmdFlag.WebPort)
	if err := http.ListenAndServe(cmdFlag.WebPort, mainRouter); err != nil {
		fmt.Fprintln(os.Stdout, "[runetcd dashboard error]", err)
		os.Exit(0)
	}
}

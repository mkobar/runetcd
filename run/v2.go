package run

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

// PutV2 puts key-value using V2 API.
func (c *Cluster) PutV2(

	w io.Writer,
	key []byte,
	val []byte,

) error {

	endpoints := []string{}
	for _, m := range c.NameToMember {
		for v := range m.Flags.ListenClientURLs {
			endpoints = append(endpoints, v)
		}
	}
	sort.Strings(endpoints)

	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
		// SelectionMode:           client.EndpointSelectionPrioritizeLeader,
	}
	ct, err := client.New(cfg)
	if err != nil {
		return err
	}
	kapi := client.NewKeysAPI(ct)

	ts := time.Now()
	if _, err := kapi.Set(context.Background(), string(key), string(val), nil); err != nil {
		return err
	}

	fmt.Fprintf(w, "[PutV2] Done! Took %v for %s/%s.\n", time.Since(ts), key, val)
	return nil
}

// GetV2 gets the value to the key using V2 API.
func (c *Cluster) GetV2(

	w io.Writer,
	key []byte,

) error {

	endpoints := []string{}
	for _, m := range c.NameToMember {
		for v := range m.Flags.ListenClientURLs {
			endpoints = append(endpoints, v)
		}
	}

	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
		// SelectionMode:           client.EndpointSelectionPrioritizeLeader,
	}
	ct, err := client.New(cfg)
	if err != nil {
		return err
	}
	kapi := client.NewKeysAPI(ct)

	ts := time.Now()
	resp, err := kapi.Get(context.Background(), string(key), nil)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "[GetV2] Done! Took %v for %s/%s.\n", time.Since(ts), key, resp.Node.Value)
	return nil
}

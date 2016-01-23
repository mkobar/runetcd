package run

// Flags is a set of etcd flags.
type Flags struct {
	// Name is a name for an etcd member.
	Name string `name`

	// ExperimentalV3Demo is either 'true' or 'false'.
	ExperimentalV3Demo bool `experimental-v3demo`

	// ExperimentalgRPCAddr is used as an endpoint for gRPC.
	// It is usually composed of host, and port '*78'.
	// Default values are '127.0.0.1:2378'.
	ExperimentalgRPCAddr string `experimental-gRPC-addr`

	// ListenClientURLs is a list of URLs to listen for clients.
	// It is usually composed of scheme, host, and port '*79'.
	// Default values are
	// 'http://localhost:2379,http://localhost:4001'.
	ListenClientURLs map[string]struct{} `listen-client-urls`

	// AdvertiseClientURLs is a list of this member's client URLs
	// to advertise to the public. The client URLs advertised should
	// be accessible to machines that talk to etcd cluster. etcd client
	// libraries parse these URLs to connect to the cluster.
	// It is usually composed of scheme, host, and port '*79'.
	// Default values are
	// 'http://localhost:2379,http://localhost:4001'.
	AdvertiseClientURLs map[string]struct{} `advertise-client-urls`

	// ListenPeerURLs is a list of URLs to listen on for peer traffic.
	// It is usually composed of scheme, host, and port '*80'.
	// Default values are
	// 'http://localhost:2380,http://localhost:7001'.
	ListenPeerURLs map[string]struct{} `listen-peer-urls`

	// InitialAdvertisePeerURLs is URL to advertise to other members
	// in the cluster, used to communicate between members.
	// It is usually composed of scheme, host, and port '*80'.
	// Default values are
	// 'http://localhost:2380,http://localhost:7001'.
	InitialAdvertisePeerURLs map[string]struct{} `initial-advertise-peer-urls`

	// InitialCluster is a map of each member name to its
	// InitialAdvertisePeerURLs.
	InitialCluster map[string]string `initial-cluster`

	// InitialClusterToken is a token specific to cluster.
	// Specifying this can protect you from unintended cross-cluster
	// interaction when running multiple clusters.
	InitialClusterToken string `initial-cluster-token`

	// InitialClusterState is either 'new' or 'existing'.
	InitialClusterState string `initial-cluster-state`

	// DataDir is a directory to store its database.
	// It should be suffixed with '.etcd'.
	DataDir string `data-dir`

	// Proxy is either 'on' or 'off'.
	Proxy bool `proxy`

	ClientCertFile      string `cert-file`        // Path to the client server TLS cert file.
	ClientKeyFile       string `key-file`         // Path to the client server TLS key file.
	ClientCertAuth      bool   `client-cert-auth` // Enable client cert authentication.
	ClientTrustedCAFile string `trusted-ca-file`  // Path to the client server TLS trusted CA key file.

	PeerCertFile       string `peer-cert-file`        // Path to the peer server TLS cert file.
	PeerKeyFile        string `peer-key-file`         // Path to the peer server TLS key file.
	PeerClientCertAuth bool   `peer-client-cert-auth` // Enable peer client cert authentication.
	PeerTrustedCAFile  string `peer-trusted-ca-file`  // Path to the peer server TLS trusted CA file.
}

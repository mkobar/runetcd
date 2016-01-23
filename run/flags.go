package run

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/gyuho/psn/ss"
)

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

var (
	DefaultExperimentalV3Demo   = true
	DefaultExperimentalgRPCAddr = "localhost:2378"
	DefaultListenClientURLs     = map[string]struct{}{
		"http://localhost:2379": struct{}{},
	}
	DefaultAdvertiseClientURLs = map[string]struct{}{
		"http://localhost:2379": struct{}{},
	}
	DefaultListenPeerURLs = map[string]struct{}{
		"http://localhost:2380": struct{}{},
	}
	DefaultAdvertisePeerURLs = map[string]struct{}{
		"http://localhost:2380": struct{}{},
	}
	DefaultInitialClusterState = "new"
)

func NewDefaultFlags() *Flags {
	fs := &Flags{}
	fs.ExperimentalV3Demo = DefaultExperimentalV3Demo
	fs.ExperimentalgRPCAddr = DefaultExperimentalgRPCAddr
	fs.ListenClientURLs = DefaultListenClientURLs
	fs.AdvertiseClientURLs = DefaultAdvertiseClientURLs
	fs.ListenPeerURLs = DefaultListenPeerURLs
	fs.InitialAdvertisePeerURLs = DefaultAdvertisePeerURLs
	fs.InitialCluster = make(map[string]string)
	fs.InitialClusterState = DefaultInitialClusterState
	return fs
}

// NewFlags returns default flags for an etcd member.
func NewFlags(
	name string,
	usedPorts *ss.Ports,
	portPrefix int,
	initialClusterToken string,
	initialClusterState string,
	dataDir string,
	isClientTLS bool,
	isPeerTLS bool,
	certPath string,
	privateKeyPath string,
	caPath string,
) (*Flags, error) {

	// To allow ports between 1178 ~ 65480.
	// Therefore, prefix must be between 11 and 654.
	if portPrefix < 11 || portPrefix > 654 {
		return nil, fmt.Errorf("portPrefix ':%d78,*79,*80' is out of port range!", portPrefix)
	}

	gRPCPort := fmt.Sprintf(":%d78", portPrefix)
	clientURLPort := fmt.Sprintf(":%d79", portPrefix)
	peerURLPort := fmt.Sprintf(":%d80", portPrefix)
	if usedPorts != nil {
		pts, err := ss.GetFreePorts(3, ss.TCP, ss.TCP6)
		if err != nil {
			return nil, err
		}
		if usedPorts.Exist(gRPCPort) {
			gRPCPort = pts[0]
		}
		if usedPorts.Exist(clientURLPort) {
			clientURLPort = pts[1]
		}
		if usedPorts.Exist(peerURLPort) {
			peerURLPort = pts[2]
		}
	}
	gRPCAddr := "localhost" + gRPCPort
	clientURL := "http://localhost" + clientURLPort
	peerURL := "http://localhost" + peerURLPort
	if isClientTLS {
		clientURL = strings.Replace(clientURL, "http://", "https://", -1)
	}
	if isPeerTLS {
		peerURL = strings.Replace(peerURL, "http://", "https://", -1)
	}

	fs := NewDefaultFlags()

	fs.Name = name

	fs.ExperimentalgRPCAddr = gRPCAddr

	fs.ListenClientURLs = map[string]struct{}{
		clientURL: struct{}{},
	}
	fs.AdvertiseClientURLs = map[string]struct{}{
		clientURL: struct{}{},
	}

	fs.ListenPeerURLs = map[string]struct{}{
		peerURL: struct{}{},
	}
	fs.InitialAdvertisePeerURLs = map[string]struct{}{
		peerURL: struct{}{},
	}

	fs.InitialClusterToken = initialClusterToken
	fs.InitialClusterState = initialClusterState

	dataDir = strings.TrimSpace(dataDir)
	if !strings.HasSuffix(dataDir, ".etcd") {
		dataDir = strings.Replace(dataDir, ".", "", -1)
		dataDir = dataDir + ".etcd"
	}
	fs.DataDir = dataDir

	if isClientTLS {
		fs.ClientCertFile = certPath
		fs.ClientKeyFile = privateKeyPath
		fs.ClientCertAuth = true
		fs.ClientTrustedCAFile = caPath
	}
	if isPeerTLS {
		fs.PeerCertFile = certPath
		fs.PeerKeyFile = privateKeyPath
		fs.PeerClientCertAuth = true
		fs.PeerTrustedCAFile = caPath
	}

	return fs, nil
}

// CombineFlags combine flags under a same cluster.
func CombineFlags(cs ...*Flags) error {
	m := make(map[string]string)
	portCheck := ""
	for _, f := range cs {
		if _, ok := m[f.Name]; ok {
			return fmt.Errorf("%s is duplicate!", f.Name)
		}
		tp := strings.Join(f.getAllPorts(), "___")
		if portCheck == "" {
			portCheck = tp
		} else if portCheck == tp {
			return fmt.Errorf("%q has duplicate ports in another member!", f.getAllPorts())
		}
		m[f.Name] = mapToCommaString(f.InitialAdvertisePeerURLs)
	}
	for _, f := range cs {
		f.InitialCluster = m
	}
	return nil
}

func (f *Flags) IsValid() (bool, error) {
	if len(f.Name) == 0 {
		return false, errors.New("Name must be specified!")
	}
	if f.InitialClusterState != "new" && f.InitialClusterState != "existing" {
		return false, errors.New("InitialClusterState must be either 'new' or 'existing'.")
	}
	return true, nil
}

func (f *Flags) String() string {

	valid, err := f.IsValid()
	if !valid || err != nil {
		return ""
	}

	pairs := [][]string{}

	nameTag, err := f.getTag("Name")
	if err != nil {
		return ""
	}
	pairs = append(pairs, []string{nameTag, strings.TrimSpace(f.Name)})

	experimentV3DemoTag, err := f.getTag("ExperimentalV3Demo")
	if err != nil {
		return ""
	}
	if f.ExperimentalV3Demo {
		pairs = append(pairs, []string{experimentV3DemoTag, "true"})
	}

	experimentalgRPCAddrTag, err := f.getTag("ExperimentalgRPCAddr")
	if err != nil {
		return ""
	}
	pairs = append(pairs, []string{experimentalgRPCAddrTag, strings.TrimSpace(f.ExperimentalgRPCAddr)})

	listenClientURLsTag, err := f.getTag("ListenClientURLs")
	if err != nil {
		return ""
	}
	pairs = append(pairs, []string{listenClientURLsTag, mapToCommaString(f.ListenClientURLs)})

	advertiseClientURLsTag, err := f.getTag("AdvertiseClientURLs")
	if err != nil {
		return ""
	}
	pairs = append(pairs, []string{advertiseClientURLsTag, mapToCommaString(f.AdvertiseClientURLs)})

	listenPeerURLsTag, err := f.getTag("ListenPeerURLs")
	if err != nil {
		return ""
	}
	pairs = append(pairs, []string{listenPeerURLsTag, mapToCommaString(f.ListenPeerURLs)})

	initialAdvertisePeerURLsTag, err := f.getTag("InitialAdvertisePeerURLs")
	if err != nil {
		return ""
	}
	pairs = append(pairs, []string{initialAdvertisePeerURLsTag, mapToCommaString(f.InitialAdvertisePeerURLs)})

	initialClusterTag, err := f.getTag("InitialCluster")
	if err != nil {
		return ""
	}
	pairs = append(pairs, []string{initialClusterTag, mapToMapString(f.InitialCluster)})

	initialClusterTokenTag, err := f.getTag("InitialClusterToken")
	if err != nil {
		return ""
	}
	pairs = append(pairs, []string{initialClusterTokenTag, f.InitialClusterToken})

	initialClusterStateTag, err := f.getTag("InitialClusterState")
	if err != nil {
		return ""
	}
	if f.InitialClusterState == "new" {
		pairs = append(pairs, []string{initialClusterStateTag, "new"})
	} else {
		pairs = append(pairs, []string{initialClusterStateTag, "existing"})
	}

	dataDirTag, err := f.getTag("DataDir")
	if err != nil {
		return ""
	}
	pairs = append(pairs, []string{dataDirTag, strings.TrimSpace(f.DataDir)})

	proxyTag, err := f.getTag("Proxy")
	if err != nil {
		return ""
	}
	if f.Proxy {
		pairs = append(pairs, []string{proxyTag, "on"})
	}

	if f.ClientCertFile != "" && f.ClientKeyFile != "" {
		clientCertTag, err := f.getTag("ClientCertFile")
		if err != nil {
			return ""
		}
		pairs = append(pairs, []string{clientCertTag, f.ClientCertFile})

		clientKeyTag, err := f.getTag("ClientKeyFile")
		if err != nil {
			return ""
		}
		pairs = append(pairs, []string{clientKeyTag, f.ClientKeyFile})

		clientClientCertAuthTag, err := f.getTag("ClientCertAuth")
		if err != nil {
			return ""
		}
		pairs = append(pairs, []string{clientClientCertAuthTag, "true"})

		clientTrustedCAFileTag, err := f.getTag("ClientTrustedCAFile")
		if err != nil {
			return ""
		}
		pairs = append(pairs, []string{clientTrustedCAFileTag, f.PeerTrustedCAFile})
	}

	if f.PeerCertFile != "" && f.PeerKeyFile != "" {
		peerCertTag, err := f.getTag("PeerCertFile")
		if err != nil {
			return ""
		}
		pairs = append(pairs, []string{peerCertTag, f.PeerCertFile})

		peerKeyTag, err := f.getTag("PeerKeyFile")
		if err != nil {
			return ""
		}
		pairs = append(pairs, []string{peerKeyTag, f.PeerKeyFile})

		peerClientCertAuthTag, err := f.getTag("PeerClientCertAuth")
		if err != nil {
			return ""
		}
		pairs = append(pairs, []string{peerClientCertAuthTag, "true"})

		peerTrustedCAFileTag, err := f.getTag("PeerTrustedCAFile")
		if err != nil {
			return ""
		}
		pairs = append(pairs, []string{peerTrustedCAFileTag, f.PeerTrustedCAFile})
	}

	sb := new(bytes.Buffer)
	for _, pair := range pairs {
		sb.WriteString(fmt.Sprintf("--%s='%s'", pair[0], pair[1]))
		sb.WriteString(" ")
	}

	return strings.TrimSpace(sb.String())
}

func (f *Flags) getAllPorts() []string {
	tm := make(map[string]struct{})
	if f.ExperimentalgRPCAddr != "" {
		ss := strings.Split(f.ExperimentalgRPCAddr, ":")
		tm[strings.TrimSpace(ss[len(ss)-1])] = struct{}{}
	}
	for k := range f.ListenClientURLs {
		ss := strings.Split(k, ":")
		tm[strings.TrimSpace(ss[len(ss)-1])] = struct{}{}
	}
	for k := range f.AdvertiseClientURLs {
		ss := strings.Split(k, ":")
		tm[strings.TrimSpace(ss[len(ss)-1])] = struct{}{}
	}
	for k := range f.ListenPeerURLs {
		ss := strings.Split(k, ":")
		tm[strings.TrimSpace(ss[len(ss)-1])] = struct{}{}
	}
	for k := range f.InitialAdvertisePeerURLs {
		ss := strings.Split(k, ":")
		tm[strings.TrimSpace(ss[len(ss)-1])] = struct{}{}
	}
	sl := make([]string, len(tm))
	i := 0
	for k := range tm {
		sl[i] = k
		i++
	}
	sort.Strings(sl)
	return sl
}

func (f *Flags) getTag(fn string) (string, error) {
	field, ok := reflect.TypeOf(f).Elem().FieldByName(fn)
	if !ok {
		return "", fmt.Errorf("Field %s is not found!", fn)
	}
	return string(field.Tag), nil
}

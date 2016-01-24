package etcdproc

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/satori/go.uuid"
)

func TestFlagsTag(t *testing.T) {
	df := NewDefaultFlags()
	field, ok := reflect.TypeOf(df).Elem().FieldByName("Name")
	if !ok {
		t.Error("Name field not found")
	}
	if string(field.Tag.Get("flag")) != "name" {
		t.Errorf("expected 'name' but got %s", string(field.Tag))
	}
}

func TestNewDefaultFlags(t *testing.T) {
	isClientTLS := false
	isPeerTLS := false
	certPath := ""
	privateKeyPath := ""
	caPath := ""
	df, err := NewFlags("etcd1", nil, 20, "etcd-cluster-token", "new", uuid.NewV4().String(), isClientTLS, isPeerTLS, certPath, privateKeyPath, caPath)
	if err != nil {
		t.Error(err)
	}
	if strings.HasPrefix(df.String(), `--name='etcd1' --experimental-v3demo='true' --experimental-gRPC-addr='localhost:2078' --listen-client-urls='http://localhost:2079' --advertise-client-urls='http://localhost:2079' --listen-peer-urls='http://localhost:2080' --initial-advertise-peer-urls='http://localhost:2080' --initial-cluster='' --initial-cluster-token='etcd-cluster-token' --initial-cluster-state='new' --data-dir=''`) {
		t.Errorf("wrong FlagString from %s", df)
	}
}

func TestCombineFlags(t *testing.T) {
	isClientTLS := false
	isPeerTLS := false
	certPath := ""
	privateKeyPath := ""
	caPath := ""
	fs := make([]*Flags, 5)
	for i := range fs {
		df, err := NewFlags(fmt.Sprintf("etcd%d", i), nil, 20+i, "etcd-cluster-token", "new", uuid.NewV4().String(), isClientTLS, isPeerTLS, certPath, privateKeyPath, caPath)
		if err != nil {
			t.Error(err)
		}
		fs[i] = df
	}
	if err := CombineFlags(fs...); err != nil {
		t.Error(err)
	}
	fa := []*Flags{fs[0], fs[0]}
	if err := CombineFlags(fa...); err == nil {
		t.Error(err)
	}
}

func TestGetAllPorts(t *testing.T) {
	isClientTLS := false
	isPeerTLS := false
	certPath := ""
	privateKeyPath := ""
	caPath := ""
	df, err := NewFlags("etcd1", nil, 20, "etcd-cluster-token", "new", uuid.NewV4().String(), isClientTLS, isPeerTLS, certPath, privateKeyPath, caPath)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(df.getAllPorts())
}

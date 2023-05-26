package manager

import (
	"context"
	"path"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type etcdLoader struct {
	dir string
	cli *clientv3.Client
}

func NewEtcdLoader(endPoints []string, dir string) *etcdLoader {
	e := new(etcdLoader)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	e.cli = cli
	e.dir = strings.Trim(dir, "/")
	return e
}

func (e *etcdLoader) Read(configSlice []string) (map[string][]byte, error) {
	//获取所有etcd的keys，设置value
	data := make(map[string][]byte)
	var err error
	res, err := e.cli.Get(context.Background(), e.dir+"/", clientv3.WithPrefix())
	if err != nil {
		return data, err
	}
	for _, kv := range res.Kvs {
		if contains(strings.Replace(path.Ext(string(kv.Key)), ".", "", -1), configSlice) {
			data[strings.Replace(string(kv.Key), e.dir+"/", "", -1)] = kv.Value
		}
	}
	return data, err
}
func (e *etcdLoader) Watch(onChange func(data map[string][]byte), configSlice []string) {
	rch := e.cli.Watch(context.Background(), e.dir+"/", clientv3.WithPrefix())
	for ev := range rch {
		if len(ev.Events) > 0 {
			data, err := e.Read(configSlice)
			if err != nil {
				panic(err)
			} else {
				onChange(data)
			}
		}
	}
}

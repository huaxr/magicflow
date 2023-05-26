// Author: XinRui Hua
// Time:   2022/3/22 下午3:38
// Git:    huaxr

package client

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/spf13/cast"
	"google.golang.org/grpc"
)

type conn struct {
	conn *grpc.ClientConn
	addr string
}

type services struct {
	srv map[int]conn
	// crc replace real host addr.
	keys []int
}

var (
	rpcServices = &services{
		srv:  make(map[int]conn),
		keys: make([]int, 0),
	}
	lock sync.RWMutex
)

func (s *services) String() string {
	var addrs = make([]string, 0)
	for _, v := range s.srv {
		addrs = append(addrs, v.addr)
	}
	return "[" + strings.Join(addrs, ", ") + "]"
}

func (s *services) hosts() (addrs []string) {
	lock.RLock()
	defer lock.RUnlock()
	for _, x := range s.srv {
		ip := strings.Split(x.addr, ":")[0]
		// port should be pull port
		port := confutil.GetProm().PullPort
		addrs = append(addrs, fmt.Sprintf("%s:%s", ip, port))
	}
	return
}

func (s *services) register(addr string) {
	lock.Lock()
	defer lock.Unlock()

	con, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logx.L().Errorf("register err when dial:", addr)
		return
	}

	crc := toolutil.CRC(toolutil.String2Byte(addr))

	if _, ok := s.srv[crc]; ok {
		logx.L().Warnf("register addr %v already exist", addr)
		return
	}

	s.srv[crc] = conn{
		conn: con,
		addr: addr,
	}
	s.keys = append(s.keys, crc)
	logx.L().Warnf("register remote addr %v, current services: %+v", addr, s)
}

func (s *services) offline(addr string) {
	lock.Lock()
	defer lock.Unlock()

	crc := toolutil.CRC(toolutil.String2Byte(addr))

	if _, ok := s.srv[crc]; !ok {
		logx.L().Warnf("delete addr %v not exist", addr)
		return
	}

	var index, add int

	for index, add = range s.keys {
		if add == crc {
			break
		}
	}

	s1 := s.keys[:index]
	s2 := s.keys[index+1:]
	s.keys = append(s1, s2...)
	delete(s.srv, crc)
	logx.L().Warnf("delete remote addr %v, current services: %+v", addr, s)
}

// if the method binding with service rather than *service.
// it will be shadowed cause. (*s) option is required here
// cause offline changed the struct.
func (s *services) get(mods interface{}) (*grpc.ClientConn, error) {
	lock.RLock()
	defer lock.RUnlock()
	switch mods.(type) {
	case string:
		if conn, ok := s.srv[cast.ToInt(mods)]; ok {
			return conn.conn, nil
		}
		return nil, fmt.Errorf("not exist addr:%v", mods)
	case int:
		c := len(s.srv)
		if c == 0 {
			logx.L().Warnf("no avaliable services currently")
			return nil, fmt.Errorf("no avaliable services currently")
		}
		mod := mods.(int) % c
		return s.srv[s.keys[mod]].conn, nil
	default:
		return nil, fmt.Errorf("unexpected mod param")
	}
}

func GetSrvAndMod() (string, uint32) {
	rand.Seed(time.Now().UnixNano())
	if c := len(rpcServices.srv); c != 0 {
		// a non-negative pseudo-random number in [0,n)
		mod := rand.Intn(core.Slot)
		return cast.ToString(rpcServices.keys[rand.Intn(c)]), uint32(mod)
	}
	return "", 0
}

func GetConn(mod interface{}) (*grpc.ClientConn, error) {
	return rpcServices.get(mod)
}

func GetRandomConn() (*grpc.ClientConn, error) {
	return rpcServices.get(time.Now().Second())

}

func GetHosts() (addrs []string) {
	return rpcServices.hosts()
}

func watcher(watchChan clientv3.WatchChan) {
	for i := range watchChan {
		// metrx here
		for _, e := range i.Events {
			k := string(e.Kv.Key)
			v := string(e.Kv.Value)

			switch e.Type {
			case mvccpb.PUT:
				logx.L().Warnf("found put %v %v", k, v)
				rpcServices.register(v)
			case mvccpb.DELETE:
				// delete option dose not contains v.
				logx.L().Warnf("found del %v %v", k, v)
				vals := strings.Split(k, "/")
				rpcServices.offline(vals[len(vals)-1])
			}
		}
	}
}

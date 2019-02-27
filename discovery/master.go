package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	client "go.etcd.io/etcd/clientv3"
)

// Master -
type Master struct {
	Name   string
	Path   string
	Nodes  map[string]*Node
	Client *client.Client
}

// Node -
type Node struct {
	State bool
	Key   string
	Info  ServiceInfo
}

// NewMaster -
func NewMaster(endpoints []string, watchPath, name string) (*Master, error) {
	cli, err := client.New(client.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	master := &Master{
		Path:   watchPath,
		Name:   name,
		Nodes:  make(map[string]*Node),
		Client: cli,
	}

	go master.WatchNodes()
	return master, err
}

// AddNode -
func (m *Master) AddNode(key string, info *ServiceInfo) {
	node := &Node{
		State: true,
		Key:   key,
		Info:  *info,
	}

	m.Nodes[node.Key] = node
}

// GetServiceInfo -
func GetServiceInfo(ev *client.Event) *ServiceInfo {
	info := &ServiceInfo{}
	err := json.Unmarshal([]byte(ev.Kv.Value), info)
	if err != nil {
		log.Println(err)
	}
	return info
}

// WatchNodes -
func (m *Master) WatchNodes() {
	t := time.Now().Unix()
	var count int
	rch := m.Client.Watch(context.Background(), m.Path, client.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case client.EventTypePut:
				count++
				nt := time.Now().Unix()

				fmt.Printf("[%s][%s] %q : %q\n", m.Name, ev.Type, ev.Kv.Key, ev.Kv.Value)
				info := GetServiceInfo(ev)
				m.AddNode(string(ev.Kv.Key), info)
				if nt-t >= 2 {
					fmt.Printf("[%s] count:%d \n", m.Name, count)
					return
				}
			case client.EventTypeDelete:
				fmt.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				delete(m.Nodes, string(ev.Kv.Key))
			}
		}
	}
}

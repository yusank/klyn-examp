package etcd

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"sync"

// 	"go.etcd.io/etcd/client"
// )

// type Server struct {
// 	info   Serverinfo // 服务器信息
// 	Status int        // 服务器状态
// }

// // all server
// type serverPool struct {
// 	services map[string]*Server // server 列表
// 	client   client.Client
// 	mu       sync.RWMutex
// }

// var (
// 	DefaultPool serverPool
// 	once        sync.Once
// )

// func Init(endpoints []string) {
// 	once.Do(func() { DefaultPool.init(endpoints) })
// }

// func (p *serverPool) init(hosts []string) {
// 	// init etcd client
// 	cfg := client.Config{
// 		Endpoints: hosts, //
// 		Transport: client.DefaultTransport,
// 	}
// 	etcdcli, err := client.New(cfg)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	p.client = etcdcli
// 	// init
// 	p.services = make(map[string]*Server)

// 	go p.watcher()
// }

// // watcher for data change in etcd directory
// func (p *serverPool) watcher() error {
// 	kAPI := client.NewKeysAPI(p.client)
// 	w := kAPI.Watcher("lc_server/", &client.WatcherOptions{Recursive: true})
// 	// 监听 "lc_server/" 当 "lc_server/" 子目录改变时能收到通知
// 	for {
// 		resp, err := w.Next(context.Background())
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		if resp.Node.Dir {
// 			continue
// 		}
// 		switch resp.Action {
// 		case "set", "create", "update", "compareAndSwap":
// 			fmt.Println(resp.Action)
// 			// 添加或更新 server 信息
// 		case "delete", "compareAndDelete", "expire":
// 			fmt.Println(resp.Action)
// 			// 当过期之后删除 server
// 		}
// 	}
// }

package etcd

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"os"
// 	"strconv"
// 	"time"

// 	"go.etcd.io/etcd/client"
// )

// // Service -
// type Service struct {
// 	ProcessId int            // 进程ID ， 单机调试时用来标志每一个服务
// 	info      Serverinfo     // 服务端信息
// 	KeysAPI   client.KeysAPI // API client, 此处用的是 V2 版本的API，是基于 http 的。 V3版本的是基于grpc的API
// }

// // Serverinfo workerInfo is the service register information to etcd
// type Serverinfo struct {
// 	Id   int32  `json:"id"`   // 服务器ID
// 	IP   string `json:"ip"`   // 对外连接服务的 IP
// 	Port int32  `json:"port"` // 对外服务端口，本机或者端口映射后得到的
// }

// // RegisterClient - 注册服务
// func RegisterClient(id int32, ip string, endpoints []string) {
// 	log.Println("register")
// 	cfg := client.Config{
// 		Endpoints:               endpoints,
// 		Transport:               client.DefaultTransport,
// 		HeaderTimeoutPerRequest: time.Second * 3,
// 	}

// 	etcdClient, err := client.New(cfg)
// 	if err != nil {
// 		panic(err)
// 	}

// 	s := Service{
// 		ProcessId: os.Getgid(),
// 		info:      Serverinfo{IP: ip, Id: id, Port: 7077},
// 		KeysAPI:   client.NewKeysAPI(etcdClient),
// 	}

// 	go s.HeartBeat()
// }

// // HeartBeat -
// func (s *Service) HeartBeat() {
// 	api := s.KeysAPI
// 	for {
// 		key := "lc_server/p_" + strconv.Itoa(s.ProcessId) // 先用 pid 来标识每一个服务， 通常应该用 IP 等来标识。
// 		// etcd 之所以适合用来做服务发现，是因为它是带目录结构的。 注册一类服务，
// 		// 只需要 key 在同一个目录下，此处 lc_sercer 目录下，p_{pid}
// 		value, _ := json.Marshal(s.info)
// 		log.Println(string(value))

// 		_, err := api.Set(context.Background(), key, string(value), &client.SetOptions{
// 			TTL: time.Second * 20,
// 		}) // 调用 API， 设置该 key TTL 为20秒。

// 		if err != nil {
// 			log.Println("Error update workerInfo:", err)
// 		}

// 		resp, err := api.Get(context.Background(), key, &client.GetOptions{Sort: true})
// 		if err != nil {
// 			log.Println("Get error:", err)
// 		}

// 		log.Printf("[%d]resp:%s \n", s.info.Id, resp.Action)

// 		time.Sleep(time.Second * 10)
// 	}
// }

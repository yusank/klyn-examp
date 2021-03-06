package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"github.com/yusank/klyn"
	"github.com/yusank/klyn-examp/discovery"
	klynlog "github.com/yusank/klyn-log"
	klconsts "github.com/yusank/klyn-log/consts"
)

// Logger - global logger
var Logger klynlog.Logger

func main() {
	// log.Println(0%2, 1%2, 2%2, 3%2, -1%2)
	// log.SetFlags(log.LstdFlags)
	// etcdMain()

	// go func() {
	// 	log.Println("start")
	// 	http.Handle("/metrics", promhttp.Handler())
	// 	log.Println(http.ListenAndServe(":8080", nil))
	// }()

	// etcd.Init([]string{"http://127.0.0.1:2379"})
	// etcd.RegisterClient(1, "127.0.0.1", []string{"http://127.0.0.1:2379"})
	// etcd.RegisterClient(2, "127.0.0.1", []string{"http://127.0.0.1:22379"})
	// etcd.RegisterClient(3, "127.0.0.1", []string{"http://127.0.0.1:32379"})

	core := klyn.Default()
	core.UseMiddleware(middleBefore, middleAfter)
	root := core.Group("")
	healthCheck(root)

	group := core.Group("/klyn")
	router(group)

	// go monitorOSSignal()

	go func() {
		time.Sleep(time.Second * 2)
	}()

	Logger = klynlog.NewLogger(&klynlog.LoggerConfig{
		Prefix:    "klyn-examp",
		IsDebug:   true,
		FlushMode: klconsts.FlushModeBySize,
	})
	core.Service(":8081")
}

func monitorOSSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL,
		syscall.SIGUSR1, syscall.SIGUSR2)

}

func setMemory() {
	file, err := os.Create("mem.prof")
	if err != nil {
		return
	}
	defer file.Close()

	var s []string

	n := time.Now().Unix()
	for {
		s = append(s, "nice", "good")

		l := len(s)
		if l%500000 == 0 {
			log.Println(l)
		}

		if time.Now().Unix() > n {
			if err = pprof.WriteHeapProfile(file); err != nil {
				log.Println("write heap profile err:", err.Error())
			}
		}
	}
}

func etcdMaster() {
	ns := rand.NewSource(time.Now().UnixNano())
	r := rand.New(ns)
	watcherName := fmt.Sprintf("m-test-%d", r.Intn(10))

	m, err := discovery.NewMaster([]string{
		"http://127.0.0.1:2379",
		"http://127.0.0.1:22379",
		"http://127.0.0.1:32379",
	}, "services/", watcherName)

	if err != nil {
		log.Fatal(err)
	}

	for {
		for k, v := range m.Nodes {
			fmt.Printf("[%s]node:%s, ip=%s\n", watcherName, k, v.Info.IP)
		}
		time.Sleep(time.Second * 10)
	}
}

func etcdService() {
	// etcd-v3
	ns := rand.NewSource(time.Now().UnixNano())
	r := rand.New(ns)
	serverName := fmt.Sprintf("s-%d", r.Intn(10))

	serviceInfo := discovery.ServiceInfo{IP: "127.0.0.1"}

	s, err := discovery.NewService(serverName, serviceInfo, []string{
		"http://127.0.0.1:2379",
		"http://127.0.0.1:22379",
		"http://127.0.0.1:32379",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name:%s, ip:%s\n", s.Name, s.Info.IP)

	if err = s.Start(); err != nil {
		s.Stop(err)
	}

	// --------------- etcd v3 end here -------------

}

func etcdMain() {
	wg := new(WaitGroupWrapper)

	// watcher
	// 监听目录的服务
	wg.Wrap(etcdMaster)
	wg.Wrap(etcdMaster)

	// service
	// 对目录进行操作的服务
	wg.Wrap(etcdService)
	wg.Wrap(etcdService)

	wg.Wait()
}

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

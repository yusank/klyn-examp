package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"git.yusank.cn/yusank/klyn-log"
	"git.yusank.space/yusank/klyn"
)

// Logger - global logger
var Logger klynlog.Logger

func main() {
	log.SetFlags(log.LstdFlags)
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

	go monitorOSSignal()

	go setMemory()

	Logger = klynlog.NewLogger(&klynlog.LoggerConfig{
		Prefix:    "klyn-examp",
		IsDebug:   true,
		FlushMode: klynlog.FlushModeBySize,
	})
	core.Service(":8081")
}

func monitorOSSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL,
		syscall.SIGUSR1, syscall.SIGUSR2)

	for {
		// 如捕捉到监听的信号，将内存中的日志写入文件
		s := <-c
		log.Println("main catch signal:", s.String())
		switch s {
		// 如果为退出信号 则安全退出
		case syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
			os.Exit(0)
		// 可以通过给进程发送 syscall.SIGUSR1, syscall.SIGUSR2 信号来，强制将缓存中的日志写入文件
		default:
		}
	}
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

// func etcdMaster() {
// 	m, err := discovery.NewMaster([]string{
// 		"http://127.0.0.1:2379",
// 		"http://127.0.0.1:22379",
// 		"http://127.0.0.1:32379",
// 	}, "services/")

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for {
// 		for k, v := range m.Nodes {
// 			fmt.Printf("node:%s, ip=%s\n", k, v.Info.IP)
// 		}
// 		fmt.Printf("nodes num = %d\n", len(m.Nodes))
// 		time.Sleep(time.Second * 5)
// 	}
// }

// func etcdService() {
// 	// etcd-v3

// 	ns := rand.NewSource(time.Now().UnixNano())
// 	r := rand.New(ns)
// 	serviceName := fmt.Sprintf("s-test-%d", r.Intn(10))
// 	serviceInfo := discovery.ServiceInfo{IP: "127.0.0.1"}

// 	s, err := discovery.NewService(serviceName, serviceInfo, []string{
// 		"http://127.0.0.1:2379",
// 		"http://127.0.0.1:22379",
// 		"http://127.0.0.1:32379",
// 	})

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("name:%s, ip:%s\n", s.Name, s.Info.IP)

// 	if err = s.Start(); err != nil {
// 		s.Stop(err)
// 	}

// 	// --------------- etcd v3 end here -------------

// }

// func etcdMain() {
// 	wg := new(WaitGroupWrapper)
// 	wg.Wrap(etcdMaster)

// 	// service
// 	wg.Wrap(etcdService)
// 	wg.Wrap(etcdService)
// 	wg.Wrap(etcdService)

// 	wg.Wait()
// }

// type WaitGroupWrapper struct {
// 	sync.WaitGroup
// }

// func (w *WaitGroupWrapper) Wrap(cb func()) {
// 	w.Add(1)
// 	go func() {
// 		cb()
// 		w.Done()
// 	}()
// }

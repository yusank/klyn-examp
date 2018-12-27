package main

import (
	"log"
	"net/http"

	"git.yusank.space/yusank/klyn"
	"git.yusank.space/yusank/klyn-log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Logger - global logger
var Logger klynlog.Logger

func main() {
	go func() {
		log.Println("start")
		http.Handle("/metrics", promhttp.Handler())
		log.Println(http.ListenAndServe(":8080", nil))
	}()

	core := klyn.Default()
	core.UseMiddleware(middleBefore, middleAfter)
	root := core.Group("")
	healthCheck(root)

	group := core.Group("/klyn")
	router(group)

	Logger = klynlog.NewLogger(&klynlog.LoggerConfig{
		Prefix:    "klyn-examp",
		IsDebug:   true,
		FlushMode: klynlog.FlushModeEveryLog,
	})
	core.Service(":8081")
}

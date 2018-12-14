package main

import (
	"git.yusank.space/yusank/klyn"
	"git.yusank.space/yusank/klyn-log"
)

// Logger - global logger
var Logger klynlog.Logger

func main() {
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

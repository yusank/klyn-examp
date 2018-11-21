package main

import (
	"log"

	"git.yusank.space/yusank/klyn"
)

func main() {
	log.SetFlags(log.Lshortfile)
	core := klyn.Default()
	core.UseMiddleware(middleBefore, middleAfter)
	root := core.Group("")
	healthCheck(root)

	group := core.Group("/klyn")
	router(group)

	core.Service(":8081")
}

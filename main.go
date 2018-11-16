package main

import (
	"git.yusank.space/yusank/klyn"
)

func main() {
	core := klyn.Default()
	core.UseMiddleware(middleBefore, middleAfter)
	group := core.Group("/klyn")
	router(group)

	core.Service(":8081")
}

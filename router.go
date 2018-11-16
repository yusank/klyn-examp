package main

import (
	"fmt"

	"git.yusank.space/yusank/klyn"
)

func router(r *klyn.RouterGroup) {
	r.GET("", testHandler, test2Handler)
	r.GET("/test", setIntHandler, getIntHandler)
}

func testHandler(c *klyn.Context) {
	c.JSON(200, "ok")
	c.Set("abc", "nice")
}

func test2Handler(c *klyn.Context) {
	fmt.Println("nice")
	fmt.Println(c.Get("abc"))
}

func setIntHandler(c *klyn.Context) {
	c.Set("int", 200)
	fmt.Println("in handler")
}

func getIntHandler(c *klyn.Context) {
	i := c.GetInt("int")
	c.JSON(200, i)
}

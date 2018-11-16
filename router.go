package main

import (
	"fmt"
	"net/http"

	"git.yusank.space/yusank/klyn"
)

func router(r *klyn.RouterGroup) {
	r.GET("", testHandler, test2Handler)
	r.GET("/test", setIntHandler, getIntHandler)
	r.GET("/healthz", healthHandler)
	r.GET("/readyz", readyHandler)
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

func healthHandler(c *klyn.Context) {
	c.JSON(http.StatusOK, nil)
}

func readyHandler(c *klyn.Context) {
	c.JSON(http.StatusOK, nil)
}

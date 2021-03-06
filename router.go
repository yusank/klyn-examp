package main

import (
	"log"
	"net/http"
	"time"

	"github.com/yusank/klyn"
)

func router(r *klyn.RouterGroup) {
	r.GET("", testHandler, test2Handler)
	r.POST("", test2Handler)
	r.GET("/test", setIntHandler, getIntHandler)

	r.GET("/healthz", healthHandler)
	r.GET("/readyz", readyHandler)
	r.GET("/ping", ping)

}

func testHandler(c *klyn.Context) {
	c.JSON(200, "ok")
	c.Set("abc", "nice")
}

func test2Handler(c *klyn.Context) {
	log.Println("nice")
	log.Println(c.Get("abc"))
}

func setIntHandler(c *klyn.Context) {
	c.Set("int", 200)
	log.Println("in handler")
}

func getIntHandler(c *klyn.Context) {
	i := c.GetInt("int")
	c.JSON(200, klyn.K{"errCode": 0, "errMsg": i})
}

func healthCheck(r *klyn.RouterGroup) {
	r.GET("/healthz", healthHandler)
	r.GET("/readyz", readyHandler)
	r.GET("/ping", ping)
}

func healthHandler(c *klyn.Context) {
	c.JSON(http.StatusOK, nil)
}

func readyHandler(c *klyn.Context) {
	c.JSON(http.StatusOK, nil)
}

func ping(c *klyn.Context) {
	Logger.Info(map[string]interface{}{
		"event":    "ping check",
		"time":     time.Now().Unix(),
		"clientIP": c.ClientIP(),
	})

	c.JSON(http.StatusOK, "pong")
}

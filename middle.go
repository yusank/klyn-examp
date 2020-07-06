package main

import (
	"log"

	"github.com/yusank/klyn"
)

func middleBefore(c *klyn.Context) {
	log.Println("before handler")
	c.Next()
}

func middleAfter(c *klyn.Context) {
	c.Next()

	log.Println("after handler")
}

package main

import (
	"fmt"
	"git.yusank.space/yusank/klyn"
)

func middleBefore(c *klyn.Context) {
	fmt.Println("before handler")
	c.Next()
}

func middleAfter(c *klyn.Context) {
	c.Next()

	fmt.Println("after handler")
}

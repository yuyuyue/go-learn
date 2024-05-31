package main

import (
	"fmt"
	"time"

	"github.com/yuyuyue/http_framework/gee"
)

func main() {
	g := gee.New()
	g.Use(gee.Recovery())
	v1 := g.Append("/v1")

	{
		v1.Use(middleware_test())
		v1.GET("/api/:id", v1IndexHandler)
	}

	// g.GET("/api/:id", indexHandler)
	g.Run(":9999")
}

func indexHandler(c *gee.Context) {
	fmt.Fprintf(c.Writer, "URL.Path = %q\n", c.Path)
}

func v1IndexHandler(c *gee.Context) {
	fmt.Fprintf(c.Writer, "V1 Group URL.Path = %q\n", c.Path)
}

func middleware_test() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		// c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		fmt.Println("V1 middle test", t)
		c.Next()
		fmt.Println("V1 middle test next", t)
	}
}

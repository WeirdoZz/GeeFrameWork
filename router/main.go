package main

import (
	"gee/gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "hello,%s", "weirdo")
	})

	r.GET("/hello/:name", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "hello,%s", ctx.Param("name"))
	})

	r.GET("/weirdo/*peach", func(ctx *gee.Context) {
		ctx.JSON(http.StatusOK, gee.H{
			"filepath": ctx.Param("peach"),
		})
	})

	r.Run(":9999")
}

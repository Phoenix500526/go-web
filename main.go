package main

import (
	"net/http"
	"webpack"
)

func main() {
	w := webpack.New()
	w.GET("/", func(c *webpack.Context) {
		c.HTML(http.StatusOK, "<h1>Hello World</h1>")
	})
	w.GET("/hello", func(c *webpack.Context) {
		c.FormatResponse(http.StatusOK, "hello, %s, welcome to %s\n", c.Query("name"), c.Path)
	})
	w.GET("/hello/:name", func(c *webpack.Context) {
		c.FormatResponse(http.StatusOK, "hello, %s, welcome to %s\n", c.Param("name"), c.Path)
	})
	w.GET("/assets/*filepath", func(c *webpack.Context) {
		c.JSON(http.StatusOK, webpack.H{"filepath": c.Param("filepath")})
	})
	w.POST("/login", func(c *webpack.Context) {
		c.JSON(http.StatusOK, webpack.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	w.Run(":2020")
}

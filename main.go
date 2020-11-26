package main

import (
	"net/http"
	"webpack"
)

func main() {
	w := webpack.New()
	w.GET("/index", func(c *webpack.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	g1 := w.Group("/v1")
	{
		g1.GET("/", func(c *webpack.Context) {
			c.HTML(http.StatusOK, "<h1>Hello World</h1>")
		})
		g1.GET("/hello", func(c *webpack.Context) {
			c.FormatResponse(http.StatusOK, "hello, welcome to %s\n", c.Query("name"), c.Path)
		})

	}

	g2 := w.Group("/v2")
	{
		g2.GET("/hello/:name", func(c *webpack.Context) {
			c.FormatResponse(http.StatusOK, "hello, %s, welcome to %s\n", c.Param("name"), c.Path)
		})
		g2.GET("/assets/*filepath", func(c *webpack.Context) {
			c.JSON(http.StatusOK, webpack.H{"filepath": c.Param("filepath")})
		})
		g2.POST("/login", func(c *webpack.Context) {
			c.JSON(http.StatusOK, webpack.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	w.Run(":2020")
}

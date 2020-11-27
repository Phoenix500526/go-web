package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"webpack"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func onlyForV2() webpack.HandlerFunc {
	return func(c *webpack.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internel Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

func main() {
	w := webpack.New()
	w.Use(webpack.Logger())
	w.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	w.LoadHTMLGlob("templates/*")
	w.Static("/assets", "./static")

	stu1 := &student{Name: "Phoenix", Age: 26}
	stu2 := &student{Name: "Jack", Age: 20}

	w.GET("/", func(c *webpack.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	w.GET("/students", func(c *webpack.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", webpack.H{
			"title":  "webpack",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	w.GET("/date", func(c *webpack.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", webpack.H{
			"title": "webpack",
			"now":   time.Date(2020, 11, 27, 0, 0, 0, 0, time.UTC),
		})
	})
	w.Run(":2020")
}

package main

import (
	"fmt"
	"net/http"
	"webpack"
)

func main() {
	w := webpack.New()
	w.GET("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "URL.Path = %q\n", req.URL.Path)
	})
	w.GET("/hello", func(res http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(res, "Header[%q] = %q\n", k, v)
		}
	})
	w.Run(":2020")
}

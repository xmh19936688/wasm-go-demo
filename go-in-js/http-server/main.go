package main

import (
	"log"
	"net/http"
	"strings"
)

// go run ./go-in-js/http-server/main.go
func main() {
	startServer()
}

func startServer() {
	fs := http.FileServer(http.Dir("static"))
	log.Println("serving html on http://localhost:8080")
	http.ListenAndServe(":8080", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		// only for debug
		resp.Header().Add("Cache-Control", "no-cache")

		// add header for wasm file
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}

		fs.ServeHTTP(resp, req)
	}))
}

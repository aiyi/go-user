package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"

	"github.com/aiyi/go-user/frontend"
)

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println(http.ListenAndServe(":8080", frontend.Engine))
}

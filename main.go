package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/aiyi/go-user/frontend"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println(http.ListenAndServe(":8080", frontend.Engine))
}

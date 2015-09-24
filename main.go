package main

import (
	"log"
	"net/http"

	"github.com/aiyi/go-user/frontend"
	_ "github.com/aiyi/go-user/frontend/user"
)

func main() {
	log.Println(http.ListenAndServe(":8080", frontend.Router))
}

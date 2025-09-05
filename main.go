package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	Info string `json:"info"`
}

var (
	Memo = make(map[int]User)
	smux = sync.RWMutex{}
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", RootHandler)
	mux.HandleFunc("/user/", userHandler)
	mux.HandleFunc("/users/", usersHandler)
	fmt.Println("Server is Running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

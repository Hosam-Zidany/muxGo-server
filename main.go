package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	mux.HandleFunc("/user/{id}", userHandler)
	fmt.Println("Server is Running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello From Root")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUser(w, r)
	case http.MethodPost:
		addUser(w, r)
		// case http.MethodDelete:
		// deleteUser(w, r)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	smux.RLock()
	user, ok := Memo[id]
	smux.RUnlock()
	if !ok {
		http.Error(w, "No Such a User With This id", http.StatusBadRequest)
		return
	}
	ret, err := json.Marshal(user)
	if err != nil {
		fmt.Println("json fail")
		return
	}

	w.Write(ret)
	w.WriteHeader(http.StatusOK)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user User
	err2 := json.NewDecoder(r.Body).Decode(&user)
	if err2 != nil {
		fmt.Println("json fail")
		return
	}
	if user.Name == "" || user.ID == 0 || user.Info == "" {
		http.Error(w, "Empty User", http.StatusBadRequest)
		return
	}
	smux.Lock()
	Memo[id] = user
	smux.Unlock()
}

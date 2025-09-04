package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello From Root")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/user/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getUser(w, r, id)
	case http.MethodPost:
		addUser(w, r, id)
	case http.MethodDelete:
		deleteUser(w, r, id)
	}
}

func getUser(w http.ResponseWriter, r *http.Request, id int) {
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

	w.WriteHeader(http.StatusOK)
	w.Write(ret)
}

func addUser(w http.ResponseWriter, r *http.Request, id int) {
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

func deleteUser(w http.ResponseWriter, r *http.Request, id int) {
	smux.RLock()
	_, ok := Memo[id]
	smux.RUnlock()
	if ok {
		delete(Memo, id)
	}
	return
}

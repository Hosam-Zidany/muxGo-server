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
	case http.MethodPut:
		updateUser(w, r, id)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ret)
}

func addUser(w http.ResponseWriter, r *http.Request, id int) {
	var user User
	err2 := json.NewDecoder(r.Body).Decode(&user)
	if err2 != nil {
		http.Error(w, "wrong json format", http.StatusBadRequest)
		return
	}
	if user.Name == "" || user.ID == 0 || user.Info == "" {
		http.Error(w, "Empty User", http.StatusBadRequest)
		return
	}
	user.ID = id
	smux.Lock()
	Memo[id] = user
	smux.Unlock()
	ret, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "wrong json format", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(ret)

}

func deleteUser(w http.ResponseWriter, r *http.Request, id int) {
	smux.Lock()
	_, ok := Memo[id]
	defer smux.Unlock()
	if ok {
		delete(Memo, id)
		w.WriteHeader(http.StatusAccepted)
	} else {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request, id int) {
	smux.Lock()
	_, ok := Memo[id]
	defer smux.Unlock()
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	var newuser User
	err := json.NewDecoder(r.Body).Decode(&newuser)
	if err != nil {
		http.Error(w, "wrong json format", http.StatusBadRequest)
		return
	}
	if newuser.Name == "" || newuser.ID == 0 || newuser.Info == "" {
		http.Error(w, "Empty User", http.StatusBadRequest)
		return
	}
	newuser.ID = id
	Memo[id] = newuser
	ret, err := json.Marshal(newuser)
	if err != nil {
		http.Error(w, "wrong user format", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(ret)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User
	for _, user := range Memo {
		users = append(users, user)
	}
	ret, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		http.Error(w, "bad json data", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ret)
}

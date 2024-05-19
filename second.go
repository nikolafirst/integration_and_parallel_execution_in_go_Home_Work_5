package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID      int
	Name    string
	Age     string
	Friends []int
}

var users = make(map[int]User)
var userID = 0

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	users[userID] = newUser
	userID++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": userID - 1})
}

func makeFriends(w http.ResponseWriter, r *http.Request) {
	var friendRequest struct {
		SourceID int json:"source_id"
		TargetID int json:"target_id"
	}
	err := json.NewDecoder(r.Body).Decode(&friendRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users[friendRequest.SourceID].Friends = append(users[friendRequest.SourceID].Friends, friendRequest.TargetID)
	users[friendRequest.TargetID].Friends = append(users[friendRequest.TargetID].Friends, friendRequest.SourceID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s and %s are now friends", users[friendRequest.SourceID].Name, users[friendRequest.TargetID].Name)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	var target struct {
		TargetID int json:"target_id"
	}
	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := users[target.TargetID].Name
	delete(users, target.TargetID)
	for _, user := range users {
		for i, friend := range user.Friends {
			if friend == target.TargetID {
				users[user.ID].Friends = append(user.Friends[:i], user.Friends[i+1:]...)
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user: %s", name)
}

func getFriends(w http.ResponseWriter, r *http.Request) {
	userID := 1
	friendList := users[userID].Friends
	json.NewEncoder(w).Encode(friendList)
}

func updateUserAge(w http.ResponseWriter, r *http.Request) {
	userID := 1
	var newAge struct {
		NewAge string json:"new_age"
	}
	err := json.NewDecoder(r.Body).Decode(&newAge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	users[userID].Age = newAge.NewAge
    
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User age successfully updated")
}

func main() {
	http.HandleFunc("/create", createUser)
	http.HandleFunc("/make_friends", makeFriends)
	http.HandleFunc("/user", deleteUser)
	http.HandleFunc("/friends/", getFriends)
	http.HandleFunc("/update_age/", updateUserAge)
	http.ListenAndServe("localhost:8081", nil)
}
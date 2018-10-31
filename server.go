package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Player is a Player
type Player struct {
	Name string
	Wins int
}

// PlayerStore is an interface allowing injection of dependencies
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}

// PlayerServer has data access, and methods for handling requests
type PlayerServer struct {
	Store PlayerStore
	http.Handler
}

// NewPlayerServer performs one time configuration and setup functionality
func NewPlayerServer(store PlayerStore) *PlayerServer {
	server := new(PlayerServer)

	server.Store = store

	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(server.handleLeague))
	router.Handle("/players/", http.HandlerFunc(server.handlePlayer))

	server.Handler = router

	return server
}
func (server *PlayerServer) handleLeague(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	json.NewEncoder(response).Encode(server.Store.GetLeague())
}

func (server *PlayerServer) handlePlayer(response http.ResponseWriter, request *http.Request) {
	player := request.URL.Path[len("/players/"):]

	switch request.Method {
	case http.MethodPost:
		server.processWin(response, player)
	case http.MethodGet:
		server.showScore(response, player)
	}
}

func (server *PlayerServer) showScore(response http.ResponseWriter, player string) {
	score := server.Store.GetPlayerScore(player)
	if score == 0 {
		response.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(response, score)
}
func (server *PlayerServer) processWin(response http.ResponseWriter, player string) {
	fmt.Println("This is nothing!")
	server.Store.RecordWin(player)
	response.WriteHeader(http.StatusAccepted)
}

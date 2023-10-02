package main

import (
	"log"
	"net/http"
	"webrtc/server"
)

func main() {
	// Initialize Rooms
	server.AllRooms.Init()

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)

	log.Println("Starting server in Port: 5003")
	err := http.ListenAndServe(":5003", nil)

	if err != nil {
		log.Fatal(err)
	}
}

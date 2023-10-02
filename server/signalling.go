package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// CreateRoomRequestHandler Create a room and Return Room ID
func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		RoomID string `json:"room_id"`
	}

	roomId := AllRooms.CreateRoom()

	json.NewEncoder(w).Encode(resp{RoomID: roomId})
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomId  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func broadcaster() {
	for {
		msg := <-broadcast
		for _, client := range AllRooms.Map[msg.RoomId] {
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)
				if err != nil {
					log.Fatal(err)
					client.Conn.Close()
				}
			}
		}
	}
}

// JoinRoomRequestHandler Join Room Request Handler
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomId, ok := r.URL.Query()["roomID"]
	if !ok {
		log.Println("roomId missing in URL Parameters")
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Web Socket Upgrade Error", err)
	}

	AllRooms.InsertIntoRoom(roomId[0], false, ws)
	go broadcaster()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Fatal("Read Error:", err)
		}
		msg.Client = ws
		msg.RoomId = roomId[0]
	}
}

package server

import (
	"encoding/json"
	"fmt"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
		fmt.Println("Message received in channel", msg)
		for _, client := range AllRooms.Map[msg.RoomId] {
			if client.Conn != msg.Client {
				AllRooms.Mutex.Lock()
				err := client.Conn.WriteJSON(msg.Message)
				if err != nil {
					fmt.Print("GOT ERROR")
					log.Fatal(err)
					client.Conn.Close()
				}
				AllRooms.Mutex.Unlock()
			}
		}
	}
}

// JoinRoomRequestHandler Join Room Request Handler
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomId, ok := r.URL.Query()["roomID"]
	if !ok {
		log.Println("roomId missing in URL Parameters")
		return
	}
	fmt.Println("Upgrading conn")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Web Socket Upgrade Error", err)
	}
	fmt.Println("ROOM ID JOIN", roomId[0])
	AllRooms.InsertIntoRoom(roomId[0], false, ws)
	fmt.Println("Connected", AllRooms.Map)
	go broadcaster()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Fatal("Read Error:", err)
		}
		log.Print("MESSAGE", msg.Message)
		msg.Client = ws
		msg.RoomId = roomId[0]

		broadcast <- msg
	}
}

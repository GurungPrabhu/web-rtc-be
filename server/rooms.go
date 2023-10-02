package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Participant struct {
	Host bool
	Conn *websocket.Conn
}

type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[string][]Participant
}


func (r *RoomMap) Init() {

}

func (r *RoomMap) GetParticipants() Participant[] {

}
func (r *RoomMap) CreateRoo() {

}

func (r *RoomMap) DeleteRoom() {

}
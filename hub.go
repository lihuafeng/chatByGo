// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "encoding/json"
  "strconv"
  "time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	//Identity of room
	roomId string

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub(roomId string) *Hub {
	return &Hub{
		roomId: 	roomId,
		broadcast:  make(chan []byte, 10),
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
		clients:    make(map[*Client]bool),
	}
}
// 组装消息体
func packageMessage(user string,msgType string, text string, mine bool) string {
  var packMess = make(map[string]string)
  packMess["ip"] = user
  packMess["type"] = msgType
  packMess["time"] = time.Now().Format("04:05")
  packMess["text"] = text
  packMess["mine"] = strconv.FormatBool(mine)

  message,_ := json.Marshal(packMess)
  return string(message)
}

func (h *Hub) run() {
  defer func() {
    close(h.unregister)
    close(h.broadcast)
  }()
  for {
    select {
    case client := <-h.unregister:
      roomMutex := roomMutexes[h.roomId]
      roomMutex.Lock()
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
        if len(h.clients) == 0 {
          house.Delete(h.roomId)
          roomMutex.Unlock()
          mutexForRoomMutexes.Lock()
          if roomMutex.TryLock() {
            if len(h.clients) == 0 {
              delete(roomMutexes, h.roomId)
            }
            roomMutex.Unlock()
          }
          mutexForRoomMutexes.Unlock()
          return
        }
      }
      roomMutex.Unlock()
    case message := <-h.broadcast:
      for client := range h.clients {
        select {
        case client.send <- message:
        default:
          close(client.send)
          delete(h.clients, client)
        }
      }
    }
  }
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")
var house = make(map[string]*Hub)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	//if r.URL.Path != "/" {
	//	http.Error(w, "Not found", http.StatusNotFound)
	//	return
	//}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
/**
多房间聊天室
房间号通过url链接传递
 */
func main() {
	flag.Parse()
	rtr := mux.NewRouter()
	rtr.HandleFunc("/{room}", serveHome)
	rtr.HandleFunc("/ws/{room}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomId := vars["room"]
		room,ok := house[roomId]
		var hub *Hub
		hub = new(Hub)
		if ok{
			hub = room
		}else{
			hub := newHub()
			house[roomId] = hub
			go hub.run()
		}
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, rtr)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

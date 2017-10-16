// Copyright 2017 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 {the "License"};
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Max time allowed to write a message to the peer.
const maxWriteWait = 10 * time.Second

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		log.Println("sending ping")
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage,
				[]byte{}, time.Now().Add(maxWriteWait)); err != nil {
				log.Println("error while sending a ping:", err)
			}
		case <-done:
			return
		}
	}
}

func showDistance(w http.ResponseWriter, r *http.Request) {
	log.Println("have distance request")

	//serve static wab page
	http.ServeFile(w, r, "/var/www/distance.html")
}

func handleWS(w http.ResponseWriter, r *http.Request, statusC <-chan bool) {
	var upgrader = websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error calling upgrade: %s\n", err)
		return
	}
	defer ws.Close()

	stdoutDone := make(chan struct{})
	//go ping(ws, stdoutDone)

	for {
		status, ok := <-statusC
		if !ok {
			log.Println("status channel cloased")
			<-stdoutDone
			return
		}
		ws.SetWriteDeadline(time.Now().Add(maxWriteWait))
		message := "in"
		if status {
			message = "out"
		}
		if err = ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			fmt.Printf("have ws error: %s\n", err)
			<-stdoutDone
			return
		}
	}
}

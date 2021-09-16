package main

import (
	"log"
	"time"

	"github.com/sacOO7/gowebsocket"
)

func reconnectAfterTime(seconds int, socket *gowebsocket.Socket) {
	timer := time.NewTimer(30 * time.Second)
	<-timer.C
	log.Println("trying to reconnect...")
	socket.Connect()
}

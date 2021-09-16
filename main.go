package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/sacOO7/gowebsocket"
)

func main() {
	deviceID := flag.String("deviceID", "", "Your Pushover device ID")
	secret := flag.String("secret", "", "Your Pushover secret")
	apiURI := flag.String("apiURI", "", "The Uri of the API, which gets called on new Message")
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	closeProgramm := false
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New("wss://client.pushover.net/push")

	socket.OnConnected = func(socket gowebsocket.Socket) {
		socket.SendText("login:" + *deviceID + ":" + *secret + "\n")
		log.Println("Connected to server")
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
		reconnectAfterTime(30, &socket)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message " + message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		respToken := string(data)
		switch respToken {
		case "!":
			log.Println("Got new Message!")
			resp := getNewMessages(secret, deviceID)
			status := deleteLastMessage(resp.Message[len(resp.Message)-1].IDStr, secret, deviceID).Status
			if status == 1 {
				resp := callAPI(*apiURI)
				if resp.StatusCode == 200 {
					log.Println("API successful called")
				}
			}
		case "#":
			log.Println("KeepAlive!")
		default:
			log.Println("Recieved binary data ", respToken)
		}
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved pong " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		if !closeProgramm {
			reconnectAfterTime(30, &socket)
		}
		return
	}

	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			closeProgramm = true
			socket.Close()
			return
		}
	}
}

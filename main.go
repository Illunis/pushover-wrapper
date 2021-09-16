package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sacOO7/gowebsocket"
)

type messages struct {
	ID    int    `json:"id"`
	IDStr string `json:"id_str"`
	/*
		"message": "This is a test alert",
		"app": "LibreNMS Work",
		"aid": 435443636951759403,
		"aid_str": "435443636951759403",
		"icon": "mnsw2ykc6qa5sbn",
		"date": 1631773887,
		"priority": 0,
		"acked": 0,
		"umid": 502550228871720147,
		"umid_str": "502550228871720147",
		"title": "Testing transport from LibreNMS",
		"url": "mailto:<mail>",
		"url_title": "Reply to <mail>",
		"queued_date": 1631773893,
		"dispatched_date": 1631773893
	*/
}

type respJSON struct {
	Message []messages `json:"messages"`
	Status  int        `json:"status"`
	Request string     `json:"request"`
	/*
		"user": {
			"quiet_hours": false,
			"is_android_licensed": true,
			"is_ios_licensed": false,
			"is_desktop_licensed": true,
			"email": "<mail>",
			"show_team_ad": "1"
		},
		"device": {
			"name": "raspberrypi",
			"encryption_enabled": false,
			"default_sound": "po",
			"always_use_default_sound": false,
			"default_high_priority_sound": "po",
			"always_use_default_high_priority_sound": false,
			"dismissal_sync_enabled": false
		},
	*/
}

func main() {
	deviceID := flag.String("deviceID", "", "Your Pushover device ID")
	secret := flag.String("secret", "", "Your Pushover secret")
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
		timer := time.NewTimer(30 * time.Second)
		<-timer.C
		log.Println("trying to reconnect...")
		socket.Connect()
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message " + message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		respToken := string(data)
		switch respToken {
		case "!":
			log.Println("Got new Message!")
			resp, err := http.Get("https://api.pushover.net/1/messages.json?secret=" + *secret + "&device_id=" + *deviceID)
			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()

			//We Read the response body on the line below.
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			//Convert the body to type string
			jsonBody := string(body)
			var respJSON respJSON
			json.Unmarshal([]byte(jsonBody), &respJSON)
			fmt.Println(respJSON.Message[len(respJSON.Message)-1])
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
			timer := time.NewTimer(30 * time.Second)
			<-timer.C
			log.Println("trying to reconnect...")
			socket.Connect()
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

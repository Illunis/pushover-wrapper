package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getNewMessages(secret *string, deviceID *string) {
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
}

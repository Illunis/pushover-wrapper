package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func bodyParser(respBody io.ReadCloser) respJSON {
	defer respBody.Close()
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		log.Fatalln(err)
	}

	//Convert the body to type string
	jsonBody := string(body)
	var respJSON respJSON
	json.Unmarshal([]byte(jsonBody), &respJSON)
	return respJSON
}

func getNewMessages(secret *string, deviceID *string) respJSON {
	resp, err := http.Get("https://api.pushover.net/1/messages.json?secret=" + *secret + "&device_id=" + *deviceID)
	if err != nil {
		log.Fatalln(err)
	}
	respJSON := bodyParser(resp.Body)
	return respJSON
}

func deleteLastMessage(messageID string, secret *string, deviceID *string) respJSON {
	values := map[string]string{"secret": *secret, "message": messageID}
	jsonData, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("https://api.pushover.net/1/devices/"+*deviceID+"/update_highest_message.json", "application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatal(err)
	}
	respJSON := bodyParser(resp.Body)
	return respJSON
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Participant struct {
	id        int
	name      string
	phone     string
	recipient *Participant
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	// List of Secret Santa participants
  // TODO: Add or remove from this list to modify participants
	participants := []Participant{
		{name: "Rudolph", phone: "+10000000000"},
		{name: "Dasher", phone: "+10000000000"},
	}

	if len(participants) < 3 {
		fmt.Println("Too few participants. You can figure this out yourself, c'mon!")
	}

	// Retain a copy of the original slice ordering for shuffling and matching purposes
	original := make([]Participant, len(participants))
	for i := 0; i < len(participants); i++ {
		original[i] = participants[i]
	}

	// Set each participant's ID to that participant's initial position in the slice
	for i := 0; i < len(participants); i++ {
		participants[i].id = i
	}

	// Shuffle the list of participants until the ordering is valid.
	// A participant whose ID (in other words, starting position) does not match that
	// participant's index after shuffling is considered valid.
	for {
		rand.Shuffle(len(participants), func(i, j int) { participants[i], participants[j] = participants[j], participants[i] })
		shuffleIsValid := true
		for i := 0; i < len(participants); i++ {
			if participants[i].id == i {
				shuffleIsValid = false
				break
			}
		}
		if shuffleIsValid {
			break
		}
	}

	// Relay an SMS to each person with their assigned Secret Santa recipient
	for i := 0; i < len(original); i++ {
		original[i].recipient = &participants[i]
		// fmt.Printf("%s -> %s\n", original[i].name, original[i].recipient.name)
		sendText(original[i].phone, "Hey "+original[i].name+"! You are assigned to "+original[i].recipient.name+" for Secret Santa!")
	}
}

// Send an text message to `phone` with body `body` using Twilio's API.
// Expects TWILIO_SID, TWILIO_TOKEN, and TWILIO_NUMBER environment variables to be set.
// See https://www.twilio.com/docs/iam/credentials/api
func sendText(phone string, body string) bool {
	accountSid, defined := os.LookupEnv("TWILIO_SID")
	if !defined {
		log.Fatalln("TWILIO_SID environment value not provided")
	}
	authToken, defined := os.LookupEnv("TWILIO_TOKEN")
	if !defined {
		log.Fatalln("TWILIO_TOKEN environment value not provided")
	}
	twilioNumber, defined := os.LookupEnv("TWILIO_NUMBER")
	if !defined {
		log.Fatalln("TWILIO_NUMBER environment value not provided")
	}

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", phone)
	msgData.Set("From", twilioNumber)
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			// https://support.twilio.com/hc/en-us/articles/223134387-What-is-a-Message-SID-
			// The Message SID is the unique ID for any message successfully created by Twilio’s API.
			// It is a 34 character string that starts with “SM…” for text messages and “MM…” for media messages.
			log.Print("sent text message (sid: " + data["sid"].(string) + ")")
			return true
		}
	} else {
		log.Print("failed to send text message (status: " + resp.Status + ")")
	}

	return false
}

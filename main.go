package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var config Config

func handleNewMessages() {
	fmt.Println("🔔 Detected unread SMS. Fetching messages...")
	readUnreadMessages()
}

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Printf("Error opening config file: %v\n", err)
		return
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Printf("Error unmarshalling config: %v\n", err)
		return
	}

	number, err := getPhoneNumber()
	if err != nil {
		fmt.Printf("Error unmarshalling config: %v\n", err)
		return
	}

	fmt.Printf("✨ Listening for SMS messages to %s!\n", number)
	go pollSMSCount(handleNewMessages)
	select {} // Block forever
}

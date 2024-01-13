package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

var ws *websocket.Conn

func ConnectToGateWay(token string, intents discordgo.Intent) error {
	var err error
	ws, _, err = websocket.DefaultDialer.Dial("wss://gateway.discord.gg", nil)
	if err != nil {
		return err
	}

	// Send IDENTIFY payload to authenticate with the gateway
	identifyPayload := fmt.Sprintf(`{
		"op": 2,
		"d": {
			"token": "%s",
			"intents": %x,  // Replace with the necessary intents for your bot
			"properties": {
				"$os": "linux",
				"$browser": "my-bot",
				"$dev		dg.Client.Get()
				ice": "my-bot"
			}
		}
	}`, token, intents)

	err = ws.WriteMessage(websocket.TextMessage, []byte(identifyPayload))
	if err != nil {
		return err
	}
	defer ws.Close()

	return nil
}

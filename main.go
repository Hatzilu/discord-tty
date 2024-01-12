package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var ws *websocket.Conn

func main() {

	// Set up Discord session
	godotenv.Load()

	// Grab your own discord user token
	// TODO: find a way to get it automatically
	token := os.Getenv("USER_TOKEN")
	if token == "" {
		panic("No token provided")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}
	wsErr := connectToGateWay(&token)
	if err != nil {
		panic(wsErr)
	}

	fmt.Println(dg.UserAgent)

	// Set up event handlers
	dg.AddHandler(messageCreate)

	// Open Discord session
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return
	}
	defer dg.Close()

	select {}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Handle incoming messages

	channel, err := s.State.Channel((m.ChannelID))
	if err != nil {
		panic(err)
	}

	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		panic(err)

	}

	fmt.Printf("[%s] #%s >> %s: %s\n", guild.Name, channel.Name, m.Author.Username, m.Message.Content)
}

func connectToGateWay(token *string) error {
	var err error
	ws, _, err = websocket.DefaultDialer.Dial("wss://gateway.discord.gg", nil)
	if err != nil {
		return err
	}

	intents := discordgo.IntentMessageContent

	// Send IDENTIFY payload to authenticate with the gateway
	identifyPayload := fmt.Sprintf(`{
		"op": 2,
		"d": {
			"token": "%s",
			"intents": %n,  // Replace with the necessary intents for your bot
			"properties": {
				"$os": "linux",
				"$browser": "my-bot",
				"$device": "my-bot"
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

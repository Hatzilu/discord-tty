package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var ws *websocket.Conn

func main() {

	// Set up Discord session
	godotenv.Load()

	token := os.Getenv("BEARER_TOKEN")
	if token == "" {
		log.Fatal("No token provided")
		return
	}

	dg, err := discordgo.New("Bearer " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	fmt.Printf(dg.UserAgent)

	// Set up WebSocket connection
	var errWs error
	ws, _, errWs = websocket.DefaultDialer.Dial("wss://gateway.discord.gg", nil)
	if errWs != nil {
		fmt.Println("Error connecting to Discord WebSocket:", errWs)
		return
	}
	defer ws.Close()

	// Set up event handlers
	dg.AddHandler(messageCreate)

	// Open Discord session
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return
	}
	defer dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Handle incoming messages
	fmt.Printf("[%s] %s: %s\n", m.GuildID, m.Author.Username, m.Content)
}

package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

var ws *websocket.Conn

func main() {
	tview.NewBox().SetBorder(true).SetTitle("Discord2")

	// Set up Discord session
	godotenv.Load()

	// Grab your own discord user token
	// TODO: find a way to get it automatically
	token := os.Getenv("USER_TOKEN")
	if token == "" {
		panic("No token provided")
	}

	dg, err := initializeDiscordClient(token)

	messagesList := initializeUi(dg)

	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}
	wsErr := connectToGateWay(token, dg.Identify.Intents)
	if err != nil {
		panic(wsErr)
	}

	fmt.Println(dg.UserAgent)
	// Set up event handlers
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m, messagesList)
	})

	// guilds := tview.NewList()
	// for i, guild := range dg.State.Guilds {
	// 	guilds.AddItem(guild.Name, "", rune(i), nil)
	// }

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate, list *tview.List) {
	// Handle incoming messages

	// channel, err := s.State.Channel((m.ChannelID))
	// if err != nil {
	// 	panic(err)
	// }

	// guild, err := s.State.Guild(m.GuildID)
	// if err != nil {
	// 	panic(err)

	// }

	var formattedMessage string
	if len(m.Message.Embeds) > 0 {
		formattedMessage = "<Embed>"
	} else if len(m.Message.Attachments) > 0 {
		formattedMessage = "<Attachment>"
	} else {
		formattedMessage = m.Message.Content
	}

	// content := fmt.printf("[%s] #%s >> %s: %s\n", guild.Name, channel.Name, m.Author.Username, formattedMessage)
	// fmt.Printf("[%s] #%s >> %s: %s\n", guild.Name, channel.Name, m.Author.Username, formattedMessage)
	// channelIdRune := []rune(m.ChannelID)
	// list.SetTitle(channel.Name)

	list.AddItem(m.Author.Username+": "+formattedMessage, "", rune(list.GetItemCount()), nil)
	// list.Draw(list.get)

}

func initializeDiscordClient(token string) (*discordgo.Session, error) {

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return dg, err
	}

	dg.Identify.Intents = discordgo.IntentsAll

	// Open Discord session
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return dg, err
	}
	defer dg.Close()

	return dg, nil
}

func initializeUi(dg *discordgo.Session) *tview.List {

	// bx := tview.NewBox().SetBorder(true).SetTitle("Discord")
	serversbx := tview.NewBox().SetBorder(true).SetTitle("Servers")
	messagesbx := tview.NewBox().SetBorder(true).SetTitle("Messages")
	inputbx := tview.NewBox().SetBorder(true)

	app := tview.NewApplication()

	messagesList := tview.NewList()

	messageInput := tview.NewInputField()
	// messagesList.SetTitle()

	guildList := tview.NewList()
	fmt.Println(dg.State.Guilds)

	for i, guild := range dg.State.Guilds {
		guildList.AddItem(guild.Name, "", rune(i+1), nil)
	}
	guildList.Box = serversbx
	messagesList.Box = messagesbx
	messageInput.Box = inputbx
	appGrid := tview.NewGrid().
		SetColumns(-1, 24, 16, -1).
		SetRows(-1, 2, 3, -1).
		// AddItem(bx, 0, 0, 6, 6, 6, 0, false).          // Top - 1 row
		AddItem(guildList, 0, 0, 6, 1, 0, 0, true).     // Left - 3 rows
		AddItem(messagesList, 0, 1, 5, 3, 0, 0, false). // Left - 3 rows
		AddItem(messageInput, 5, 1, 1, 3, 0, 0, false)  // Left - 3 rows
		// AddItem(bx, 0, 3, 3, 3, 0, 0, false) // Right - 3 rows
		// AddItem(bx, 3, 1, 1, 1, 0, 0, false) // Bottom - 1 row
		// AddItem(label, 1, 1, 1, 1, 0, 0, false).
		// AddItem(input, 1, 2, 1, 1, 0, 0, false).
		// AddItem(btn, 2, 1, 1, 2, 0, 0, false)

	if err := app.SetRoot(appGrid, true).Run(); err != nil {
		panic(err)
	}

	return messagesList
}

func connectToGateWay(token string, intents discordgo.Intent) error {
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

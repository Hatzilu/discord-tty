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
	serversBox := tview.NewBox().SetBorder(true).SetTitle("Guilds")
	textChannelsBox := tview.NewBox().SetBorder(true)
	messagesBox := tview.NewBox().SetBorder(true).SetTitle("Messages")
	inputBox := tview.NewBox().SetBorder(true)

	app := tview.NewApplication()

	messagesList := tview.NewList()
	channelsList := tview.NewList()
	messageInput := tview.NewInputField()

	guildList := tview.NewList()

	appGrid := tview.NewGrid().
		SetColumns(20).
		SetRows(20).
		// AddItem(bx, 0, 0, 6, 6, 6, 0, false).          // Top - 1 row
		AddItem(guildList, 0, 0, 6, 1, 1, 1, true).     // Left - 6 rows
		AddItem(channelsList, 0, 1, 6, 1, 0, 0, false). // Left - 6 rows
		AddItem(messagesList, 0, 2, 1, 2, 0, 0, false). // Left - 5 rows
		AddItem(messageInput, 3, 1, 1, 3, 0, 0, false)  // Left - 3 rows
	// AddItem(bx, 0, 3, 3, 3, 0, 0, false) // Right - 3 rows
	// AddItem(bx, 3, 1, 1, 1, 0, 0, false) // Bottom - 1 row
	// AddItem(label, 1, 1, 1, 1, 0, 0, false).
	// AddItem(input, 1, 2, 1, 1, 0, 0, false).
	// AddItem(btn, 2, 1, 1, 2, 0, 0, false)

	for i, guild := range dg.State.Guilds {
		guildList.AddItem(guild.Name, guild.ID, rune(i), nil)
	}

	guildList.SetSelectedFunc(func(i int, guildName string, guildId string, r rune) {
		messagesList.Clear()
		channelsList.Clear()
		app.SetFocus(textChannelsBox)
		appGrid.RemoveItem(guildList)
		channelsList.SetTitle(guildName)
		appGrid.AddItem(channelsList, 0, 0, 6, 1, 1, 1, true) // Left - 6 rows
		guild, err := dg.State.Guild(guildId)
		if err != nil {
			fmt.Printf("Failed to get guild by id \"%s\"", guildId)
			panic(err)
		}
		fmt.Println(guild.Name)

		for j, channel := range guild.Channels {
			if channel.Type == discordgo.ChannelTypeGuildText {
				channelsList.AddItem("#"+channel.Name, channel.ID, rune(j), nil)
			}
		}
	})

	channelsList.SetSelectedFunc(func(i int, channelName string, channelId string, r rune) {
		messagesList.Clear()
		app.SetFocus(messageInput)
		channel, err := dg.State.Channel(channelId)
		if err != nil {
			fmt.Printf("Failed to get guild by id \"%s\"", channelId)
			panic(err)
		}

		for j, m := range channel.Messages {
			if j < 20 {
				var formattedMessage string
				if len(m.Embeds) > 0 {
					formattedMessage = "<Embed>"
				} else if len(m.Attachments) > 0 {
					formattedMessage = "<Attachment>"
				} else {
					formattedMessage = m.Content
				}
				messagesList.AddItem(m.Author.Username+": "+formattedMessage, "", rune(j), nil)
			}
		}
	})

	guildList.Box = serversBox
	messagesList.Box = messagesBox
	messageInput.Box = inputBox
	channelsList.Box = textChannelsBox

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

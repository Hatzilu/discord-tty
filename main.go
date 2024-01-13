package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

func main() {
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
	wsErr := ConnectToGateWay(token, dg.Identify.Intents)
	if err != nil {
		panic(wsErr)
	}

	// Set up event handlers
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m, messagesList)
	})
}

func initializeUi(dg *discordgo.Session) *tview.List {

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

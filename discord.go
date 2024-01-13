package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rivo/tview"
)

func initializeDiscordClient(token string) (*discordgo.Session, error) {

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return dg, err
	}

	dg.Identify.Intents = discordgo.IntentsAll

	dg.State.MaxMessageCount = 100

	// Open Discord session
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return dg, err
	}
	defer dg.Close()

	return dg, nil
}

func formatDiscordMessage(m *discordgo.Message) string {
	var formattedMessage string
	if len(m.Embeds) > 0 {
		formattedMessage = "<Embed>"
	} else if len(m.Attachments) > 0 {
		formattedMessage = "<Attachment>"
	} else {
		formattedMessage = m.Content
	}
	return formattedMessage
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

	formattedMessage := formatDiscordMessage(m.Message)

	// content := fmt.printf("[%s] #%s >> %s: %s\n", guild.Name, channel.Name, m.Author.Username, formattedMessage)
	// fmt.Printf("[%s] #%s >> %s: %s\n", guild.Name, channel.Name, m.Author.Username, formattedMessage)
	// channelIdRune := []rune(m.ChannelID)
	// list.SetTitle(chInitializeDiscordClientannel.Name)

	list.AddItem(m.Author.Username+": "+formattedMessage, "", rune(list.GetItemCount()), nil)
	// list.Draw(list.get)

}

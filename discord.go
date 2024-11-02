package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rivo/tview"
)

type Guild struct {
	Banner      string   `json:"banner"`
	Features    []string `json:"features"`
	Icon        string   `json:"icon"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Owner       bool     `json:"owner"`
	Permissions string   `json:"permissions"`
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    uint32 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type DiscordResponse interface{ []Guild }

func discordApiRequest[T DiscordResponse](method string, endpoint string, body io.Reader, token string) (data T, err error) {
	baseUrl := os.Getenv("DISCORD_API_ENDPOINT")
	if baseUrl == "" {
		return nil, errors.New("missing discord API endpoint env var")
	}

	url := strings.Join([]string{baseUrl, endpoint}, "")

	r, err := http.NewRequest(method, url, body)

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		return nil, err
	}
	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Failed to parse token response body")
		panic(err)
	}
	data = nil
	if err := json.Unmarshal(bodyBytes, &data); err != nil { // Parse []byte to go struct pointer
		fmt.Printf("Failed to unmarshal JSON: %s", err)
		return nil, err
	}
	return data, nil
}

func initializeDiscordClient(token string) (*discordgo.Session, error) {

	dg, err := discordgo.New("Bearer " + token)
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

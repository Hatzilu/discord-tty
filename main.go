package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

func main() {

	// Set up Discord session
	godotenv.Load()

	token := TokenResponse{}
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    os.Getenv("APP_ADDR"),
		Handler: mux,
	}

	fmt.Println("No valid token found, please authenticate again.")

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			fmt.Println("Invalid code value")
			io.WriteString(w, "Invalid code value\n")
			return
		}

		// auth with discord using the code from query params
		u, err := url.Parse("https://discord.com/api/v10/oauth2/token")
		if err != nil {
			fmt.Println("Failed to parse discord API url")
			panic(err)
		}

		u.OmitHost = false

		data := u.Query()
		data.Set("client_id", os.Getenv("APP_CLIENT_ID"))
		data.Set("client_secret", os.Getenv("APP_CLIENT_SECRET"))
		data.Set("grant_type", "authorization_code")
		data.Set("redirect_uri", fmt.Sprintf("%s/auth", os.Getenv("APP_REDIRECT_URL")))
		data.Set("scope", "identify+guilds+guilds.channels.read+messages.read+guilds.members.read")
		data.Set("code", code)

		u.RawQuery = data.Encode()
		r, err2 := http.NewRequest("POST", "https://discord.com/api/v10/oauth2/token", strings.NewReader(u.RawQuery))

		if err2 != nil {
			fmt.Println("Failed to init POST request")
			panic(err2)
		}

		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("Content-Length", fmt.Sprint(len(u.RawQuery)))

		// fmt.Printf("posting to %s \n\n", r.URL)
		// fmt.Printf("request:  %+v\n", r)
		// fmt.Printf("Body:  %+v\n", strings.NewReader(u.RawQuery))

		client := http.Client{}
		res, err := client.Do(r)
		if err != nil {
			fmt.Println("Failed to request token from discord")
			panic(err)
		}

		defer res.Body.Close()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Failed to parse token response body")
			panic(err)
		}

		if err := json.Unmarshal(bodyBytes, &token); err != nil { // Parse []byte to go struct pointer
			fmt.Printf("Failed to unmarshal JSON: %s", err)
		}

		fmt.Printf("Authenticated successfully, shutting down local server... token: %s", token.AccessToken)
		go srv.Shutdown(context.Background())
	})

	err := srv.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	// now that we authenticated the user, we can use the token to auth with discord and make requests to the discord API
	_guilds, err := discordApiRequest[[]discordgo.UserGuild]("GET", "/users/@me/guilds", nil, &token.AccessToken)
	if err != nil {
		panic(err)
	}

	// initializeUi(guilds, &token.AccessToken)

	fmt.Println(_guilds)
	// // Set up event handlers
	// dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 	messageCreate(s, m, messagesList)
	// })
}

func initializeUi(guilds []discordgo.UserGuild, tokenPtr *string) *tview.List {

	guildsBox := tview.NewBox().SetBorder(true).SetTitle("Guilds")
	textChannelsBox := tview.NewBox().SetBorder(true)
	messagesBox := tview.NewBox().SetBorder(true).SetTitle("Messages")
	inputBox := tview.NewBox().SetBorder(true)

	app := tview.NewApplication()

	messagesList := tview.NewList()
	channelsList := tview.NewList()
	guildList := tview.NewList()
	messageInput := tview.NewInputField()

	guildList.Box = guildsBox
	messagesList.Box = messagesBox
	messageInput.Box = inputBox
	channelsList.Box = textChannelsBox

	appGrid := tview.NewGrid().
		SetColumns(20).
		SetRows(20).
		// AddItem(bx, 0, 0, 6, 6, 6, 0, false).          // Top - 1 row
		AddItem(guildList, 0, 0, 6, 1, 1, 1, true). // Left - 6 rows
		// AddItem(channelsList, 0, 1, 6, 1, 0, 0, false). // Left - 6 rows
		AddItem(messagesList, 0, 1, 1, 2, 0, 0, false). // Left - 5 rows
		AddItem(messageInput, 3, 1, 1, 3, 0, 0, false)  // Left - 3 rows
		// AddItem(bx, 0, 3, 3, 3, 0, 0, false) // Right - 3 rows
		// AddItem(bx, 3, 1, 1, 1, 0, 0, false) // Bottom - 1 row
		// AddItem(label, 1, 1, 1, 1, 0, 0, false).
		// AddItem(input, 1, 2, 1, 1versBox
		// messagesList.Box = messagesBox
		// messageInput.Box = inputBox
		// channelsList.Box = , 0, 0, false).
		// AddItem(btn, 2, 1, 1, 2, 0, 0, false)
	for i, guild := range guilds {
		fmt.Printf("Adding guild %s\n", guild.Name)
		guildList.AddItem(guild.Name, string(guild.ID), rune(i), nil)
	}

	guildList.SetSelectedFunc(func(i int, guildName string, guildId string, r rune) {
		fmt.Println("hi")
		if messagesList.GetItemCount() > 0 {
			messagesList.Clear()
		}
		if channelsList.GetItemCount() > 0 {
			channelsList.Clear()
		}
		app.SetFocus(textChannelsBox)
		// appGrid.RemoveItem(guildList)
		channelsList.SetTitle(guildName)
		appGrid.AddItem(channelsList, 0, 0, 6, 1, 1, 1, true) // Left - 6 rows
		url := fmt.Sprintf("/guilds/%X/channels", guildId)
		channels, err := discordApiRequest[[]discordgo.Channel]("GET", url, nil, tokenPtr)
		fmt.Printf("Channels: %+v", channels)
		if err != nil {
			fmt.Printf("Failed to get channels for guild  \"%s\"", guildId)
			panic(err)
		}

		for j, channel := range channels {
			if channel.Type == discordgo.ChannelTypeGuildText {
				channelsList.AddItem("#"+channel.Name, channel.ID, rune(10+j), nil)
			}
		}
	})

	// channelsList.SetSelectedFunc(func(i int, channelName string, channelId string, r rune) {
	// 	// messagesList.Clear()
	// 	// app.SetFocus(messagesBox)
	// 	channel, err := dg.State.Channel(channelId)
	// 	if err != nil {
	// 		fmt.Printf("Failed to get channel by id \"%s\"", channelId)
	// 		panic(err)
	// 	}

	// 	if len(channel.Messages) < 1 {
	// 		messagesList.Clear()
	// 		messagesList.AddItem("No messages in "+channelName, "", 0, nil)
	// 		return
	// 	}

	// 	for j, m := range channel.Messages {
	// 		formattedMessage := formatDiscordMessage(m)
	// 		fmt.Println(m.Author.Username + ": " + formattedMessage)
	// 		messagesList.AddItem(m.Author.Username+": "+formattedMessage, "", rune(j), nil)
	// 	}
	// })

	if err := app.SetRoot(appGrid, true).Run(); err != nil {
		panic(err)
	}
	return messagesList
}

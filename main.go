package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Channel  = flag.String("channel", "", "Channel ID")
	BotToken = flag.String("token", "", "Bot token")
)

type Author struct {
	User string `json:"username"`
	Role string `json:"title"`
}

type Post struct {
	Summary   string `json:"title"`
	Message   string `json:"excerpt"`
	Url       string `json:"url"`
	Author    Author `json:"user"`
	TimeStamp string `json:"created_at"`
}

func init() { flag.Parse() }

func main() {
	const baseUrl = "https://forums.newworld.com"

	const getDevPostsRoute = "/groups/Developer/posts.json"

	c := http.Client{Timeout: time.Duration(1) * time.Second}

	resp, err := c.Get(baseUrl + getDevPostsRoute)

	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}

	var posts []Post

	json.Unmarshal([]byte(body), &posts)

	session, err := discordgo.New("Bot " + *BotToken)

	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}

	getLatestChannelMessage(session)

	for _, post := range posts {
		sendDiscordMessage(session, baseUrl, post)
	}
}

func sendDiscordMessage(session *discordgo.Session, baseUrl string, post Post) {

	var color int

	switch post.Author.Role {
	case "Community Manager":
		color = 3066993
	case "New World Developer":
		color = 15158332
	default:
		color = 8359053
	}

	_, err := session.ChannelMessageSendEmbed(*Channel, &discordgo.MessageEmbed{
		URL:         baseUrl + post.Url,
		Description: post.Message,
		Color:       color,
		Title:       post.Summary,
		Timestamp:   post.TimeStamp,
	})
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func getLatestChannelMessage(session *discordgo.Session) time.Time {
	latestChannelMessage, err := session.ChannelMessages(*Channel, 1, "", "", "")
	if err != nil {
		log.Printf("Could not get latest message: %v", err)
	}
	return latestChannelMessage[0].Timestamp
}

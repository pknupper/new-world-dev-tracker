package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Author struct {
	User string `json:"username"`
	Role string `json:"title"`
}

type Post struct {
	Summary string `json:"title"`
	Message string `json:"excerpt"`
	Url     string `json:"url"`
	Author  Author `json:"user"`
}

type Message struct {
	*discordgo.MessageEmbed
}

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

	var posts []Post

	json.Unmarshal([]byte(body), &posts)

	token := "test123"
	channel := "channel-1"

	discord, err := discordgo.New(token)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	for _, post := range posts {
		var message Message

		message.URL = baseUrl + post.Url
		message.Description = post.Message
		message.Author.Name = post.Author.User + "(" + post.Author.Role + ")"
		message.Color = 15158332
		message.Title = "New Developer Post"

		discord.ChannelMessageSendEmbed(channel, message.MessageEmbed)
	}

}

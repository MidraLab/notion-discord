package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscordWebhook struct {
	UserName  string         `json:"username"`
	AvatarURL string         `json:"avatar_url"`
	Content   string         `json:"content"`
	Embeds    []DiscordEmbed `json:"embeds"`
	TTS       bool           `json:"tts"`
}

type DiscordImage struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type DiscordAuthor struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon_url"`
}

type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type DiscordEmbed struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
	Color       int            `json:"color"`
	Image       DiscordImage   `json:"image"`
	Thumbnail   DiscordImage   `json:"thumbnail"`
	Author      DiscordAuthor  `json:"author"`
	Fields      []DiscordField `json:"fields"`
}

func NewDiscordWebhook(userName, avatarURL, content string, embeds []DiscordEmbed, tts bool) *DiscordWebhook {
	return &DiscordWebhook{
		UserName:  userName,
		AvatarURL: avatarURL,
		Content:   content,
		Embeds:    embeds,
		TTS:       tts,
	}
}

func (webhook *DiscordWebhook) AddEmbeds(embeds ...DiscordEmbed) {
	webhook.Embeds = append(webhook.Embeds, embeds...)
}

func (webhook *DiscordWebhook) SendWebhook(webhookURL string) error {
	jsonData, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf("json error: %s", err.Error())
	}

	request, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("request creation error: %s", err.Error())
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("client error: %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode == 204 {
		fmt.Println("sent", webhook) //成功
	} else {
		return fmt.Errorf("%#v\n", response) //失敗
	}

	return nil
}

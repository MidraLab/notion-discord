package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	notionApiInstance := &NotionAPI{
		DatabaseURL: loadEnv("NOTION_DATABASE_URL"),
		APIKey:      loadEnv("MIDRA_LAB_NOTION_API"),
	}

	pageID, pageURL, err := notionApiInstance.ReadPageID()
	if err != nil {
		log.Fatal(err)
	}

	if err := notionApiInstance.PatchPageTitle(pageID); err != nil {
		log.Fatalf("Failed to patch page title: %v", err)
	}

	discordWebhookInstance := NewDiscordWebhook("NotificationMTG", "https://source.unsplash.com/random", "定例ドキュメントの更新お願いします！！"+pageURL, nil, false)

	discordWebhookUrl := loadEnv("DISCORD_WEBHOOK_URL")

	if err := discordWebhookInstance.SendWebhook(discordWebhookUrl); err != nil {
		log.Fatal(err)
	}
}

func loadEnv(keyName string) string {
	err := godotenv.Load(".env")
	// もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
	// .envの SAMPLE_MESSAGEを取得して、messageに代入します。
	envValue := os.Getenv(keyName)

	return envValue
}

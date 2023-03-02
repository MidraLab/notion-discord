package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	notionAPI := &NotionAPI{
		DatabaseURL: loadEnv("NOTION_DATABASE_URL"),
		APIKey:      loadEnv("MIDRA_LAB_NOTION_API"),
	}

	pageID, pageURL, err := notionAPI.ReadPageID()
	if err != nil {
		log.Fatal(err)
	}

	if err := notionAPI.PatchPageTitle(pageID); err != nil {
		log.Fatalf("failed to patch page title: %v", err)
	}

	dw := NewDiscordWebhook("NotificationMTG", "https://source.unsplash.com/random", " 定例ドキュメントの更新お願いします！！"+pageURL, nil, false)

	whURL := loadEnv("DISCORD_WEBHOOK_URL")

	if err := dw.SendWebhook(whURL); err != nil {
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
	message := os.Getenv(keyName)

	return message
}

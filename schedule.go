package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var BotToken string

func init() {
	err := godotenv.Load("local.env")
	if err != nil {
		fmt.Println("Failed to read environemnt variables.")
	}
	BotToken = os.Getenv("DISCORD_BOT_AUTH_KEY")
}

func main() {
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	discord.AddHandler(trackCreation)

	discord.Identify.Intents = discordgo.IntentsGuildScheduledEvents

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening session,", err)
		return
	}

	defer discord.Close()

	fmt.Println("Made it this far")

	for true {

	}
}

func trackCreation(s *discordgo.Session, m *discordgo.GuildScheduledEventCreate) {
	fmt.Println("Event Created.")
}

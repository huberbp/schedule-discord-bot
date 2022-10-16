package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	fmt.Println("App started:", os.Getpid())

	// Allocate memory for signal handler
	stop := make(chan os.Signal, 1)
	// Set up signal handler to handle all possible outcomes
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// Wait indefinitely for one of these signals (this is blocking)
	<-stop
	// Notify user of shutdown
	fmt.Println("Shut Down")
}

func trackCreation(s *discordgo.Session, m *discordgo.GuildScheduledEventCreate) {
	fmt.Println("Event Created.")
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var BotToken string

func init() {
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("ERROR", err, "- Failed to open log file for log output.")
		return
	}
	log.SetOutput(logFile)
	go log.Println("Initializing new session") // According to stack overflow, log is threadsafe

	err = godotenv.Load("local.env")
	if err != nil {
		fmt.Println("ERROR", err, "- Failed to read environemnt variables.")
		log.Println("ERROR", err, "- Failed to read environemnt variables.")
		return
	}
	BotToken = os.Getenv("DISCORD_BOT_AUTH_KEY")
}

func main() {
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		fmt.Println("ERROR", err, "- Error creating Discord session.")
		log.Println("ERROR", err, "- Error creating Discord session.")
		return
	}

	// Add handlers (trackCreation is a fointer)
	discord.AddHandler(trackCreation)

	// Setup the intents to send to discord on the initial handshake greeting
	discord.Identify.Intents = discordgo.IntentsGuildScheduledEvents

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening session,", err)
		log.Println("Error opening session,", err)
		return
	}

	defer discord.Close()

	fmt.Println("App started:", os.Getpid())
	go log.Println("App started:", os.Getpid())

	// Allocate memory for signal handler
	stop := make(chan os.Signal, 1)
	// Set up signal handler to handle all possible outcomes
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// Wait indefinitely for one of these signals (this is blocking)
	<-stop
	// Notify user of shutdown
	fmt.Println("Shut Down")
	log.Println("Shut Down")
}

func trackCreation(s *discordgo.Session, m *discordgo.GuildScheduledEventCreate) {
	fmt.Println("Event Created.")
}

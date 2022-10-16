package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	BotToken        string
	logFile         *os.File
	err             error
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
)

func init() {
	commandHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

	logFile, err = os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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
	defer logFile.Close()
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		fmt.Println("ERROR", err, "- Error creating Discord session.")
		log.Println("ERROR", err, "- Error creating Discord session.")
		return
	}

	// Check application commands and update them if different
	// appCommands, err := discord.ApplicationCommands(os.Getenv("DISCORD_APP_ID"), "")
	// if err != nil {
	// 	fmt.Println("ERROR", err, "- Error getting app commands.")
	// 	log.Println("ERROR", err, "- Error getting app commands.")
	// 	return
	// }

	directories, err := os.ReadDir("app_commands")
	if err != nil {
		fmt.Println("ERROR", err, "- Error reading local command store.")
		log.Println("ERROR", err, "- Error reading local command store.")
		return
	}

	for _, directory := range directories { // all "directories" will be files anyways here
		var dirName string = directory.Name()

		commandJSON, err := os.ReadFile("app_commands/" + dirName)
		if err != nil {
			fmt.Println("ERROR", err, "- Error reading local", dirName, "command.")
			log.Println("ERROR", err, "- Error reading local", dirName, "command.")
			return
		}

		var localApp discordgo.ApplicationCommand

		json.Unmarshal(commandJSON, &localApp)

		// var discordApp discordgo.ApplicationCommand = *appCommands[index]

		_, err = discord.ApplicationCommandCreate(os.Getenv("DISCORD_APP_ID"), "", &localApp)
		if err != nil {
			fmt.Println("ERROR", err, "- Error updating", localApp.Name, "command.")
			log.Println("ERROR", err, "- Error updating", localApp.Name, "command.")
			return
		}

		switch localApp.Name {
		case "schedule":
			commandHandlers[localApp.Name] = schedule
		}
	}

	// Add handlers (trackCreation is a fointer)
	discord.AddHandler(trackCreation)

	// Add our interaction handlers (for appCommands)
	discord.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[interaction.ApplicationCommandData().Name]; ok {
			handler(session, interaction)
		}
	})

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

func schedule(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Schedule your event!",
		},
	})
}

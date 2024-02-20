package main

import (
	"backupbot/handlers"
	"backupbot/utilities"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/google/go-github/v59/github"
)

type Config struct {
	discordToken    string
	githubToken     string
	repositoryOwner string
	repositoryName  string
}

func ParseConfig() *Config {
	botToken := flag.String("discord-token", "", "Discord Bot access token")
	ghToken := flag.String("github-token", "", "GitHub access token")
	snapRepoOwner := flag.String("repository-owner", "", "The owner of the GitHub repository to store snaps to")
	snapRepoName := flag.String("repository-name", "", "The name of the GitHub repository to store snaps to")
	flag.Parse()

	return &Config{
		discordToken:    *botToken,
		githubToken:     *ghToken,
		repositoryOwner: *snapRepoOwner,
		repositoryName:  *snapRepoName,
	}
}

func main() {
	config := ParseConfig()

	s, err := discordgo.New("Bot " + config.discordToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	neptunebotMessageHandler := CreateNeptuneBotMessageHandler(config)
	var sessionDecorators = []func(*discordgo.Session) *discordgo.Session{
		onReadyNotifier,
		neptunebotMessageHandler.AddMessageHandlers,
	}

	for _, decorator := range sessionDecorators {
		s = decorator(s)
	}

	s.Identify.Intents |= discordgo.IntentsGuildMessages
	s.Identify.Intents |= discordgo.IntentsMessageContent
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	log.Println("Gracefully shutting down.")
}

func onReadyNotifier(s *discordgo.Session) *discordgo.Session {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	return s
}

func CreateNeptuneBotMessageHandler(cfg *Config) *handlers.MessageHandler {
	ghClient := github.NewClient(nil).WithAuthToken(cfg.githubToken)
	dm := utilities.NewDownloadManager()
	snapHandler := handlers.NewSnapHandler(ghClient, dm, cfg.repositoryOwner, cfg.repositoryName)
	messageHandler := handlers.NewMessageHandler(
		[]func(s *discordgo.Session, m *discordgo.MessageCreate){
			snapHandler.Snap,
		},
	)
	return messageHandler
}

package handlers

import (
	"backupbot/utilities"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/go-github/v59/github"
)

const (
	folderDelimiter = "/"
)

type SnapHandler struct {
	githubClient    *github.Client
	downloadManager *utilities.DownloadManager
	repositoryOwner string
	repositoryName  string
}

func NewSnapHandler(
	githubClient *github.Client,
	downloadManager *utilities.DownloadManager,
	backupRepoOwner string,
	backupRepoName string,
) *SnapHandler {
	return &SnapHandler{
		githubClient:    githubClient,
		downloadManager: downloadManager,
		repositoryOwner: backupRepoOwner,
		repositoryName:  backupRepoName,
	}
}

func (sh *SnapHandler) Snap(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Message.Attachments) == 0 {
		return
	}

	// this info could be cached to reduce api calls
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		log.Printf("Could not determine channel name for message %s\n", m.ID)
		return
	}

	// this info could be cached to reduce api calls
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Printf("Could not determine channel name for message %s\n", m.ID)
		return
	}

	for _, a := range m.Attachments {
		data, err := sh.downloadManager.DownloadFile(a.URL)
		if err != nil {
			log.Printf("unable to snap file %s: %s\n", a.ID, err)
			continue
		}

		outputPath := generateOutputPath(guild.Name, channel.Name, a.Filename, m.Timestamp)
		msg := fmt.Sprintf("upload %s from message in %s", outputPath, m.ChannelID)
		fileOptions := &github.RepositoryContentFileOptions{
			Message: &msg,
			Content: data,
		}
		_, _, err = sh.githubClient.Repositories.CreateFile(context.Background(), sh.repositoryOwner, sh.repositoryName, outputPath, fileOptions)
		if err != nil {
			log.Printf("failed to backup file %s: %s\n", a.ID, err)
			continue
		}
		log.Printf("backed up %s\n", outputPath)
	}
}

func generateOutputPath(guildName string, channelName string, inputFileName string, timestamp time.Time) string {
	readableDateString := timestamp.Format("2006-01-02 15:04:05")
	fileName := strings.Join([]string{readableDateString, inputFileName}, " ")
	return strings.Join([]string{guildName, channelName, fileName}, folderDelimiter)
}

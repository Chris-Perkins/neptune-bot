package handlers

import "github.com/bwmarrin/discordgo"

type MessageHandler struct {
	handlers []func(s *discordgo.Session, m *discordgo.MessageCreate)
}

func NewMessageHandler(handlers []func(s *discordgo.Session, m *discordgo.MessageCreate)) *MessageHandler {
	return &MessageHandler{
		handlers: handlers,
	}
}

func (mh *MessageHandler) AddMessageHandlers(s *discordgo.Session) *discordgo.Session {
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		for _, h := range mh.handlers {
			h(s, m)
		}
	})
	return s
}

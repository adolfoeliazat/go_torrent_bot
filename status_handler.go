package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhulik/margelet"
)

type statusHandler struct {
	path               string
	client             torrentClient
	authorizedUsername string
}

func newStatusHandler(authorizedUsername string, path string, client torrentClient) *statusHandler {
	return &statusHandler{path, client, authorizedUsername}
}

func (responder statusHandler) Response(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	if message.From.UserName != responder.authorizedUsername {
		bot.QuickSend(message.Chat.ID, "Sorry, you are not allowed to control me!")
		return nil
	}

	if len(responder.client.Torrents()) == 0 {
		bot.QuickSend(message.Chat.ID, "There is no downloads")
		return nil
	}

	if message.ReplyToMessage != nil {
		message = *message.ReplyToMessage
	}

	torrent, hash, err := findTorrentByMessage(responder.client, message)

	if err != nil {
		for _, t := range responder.client.Torrents() {
			bot.QuickSend(message.Chat.ID, torrentStats(t))
		}
		return nil
	}

	if torrent != nil {
		bot.QuickSend(message.Chat.ID, verboseTorrentStats(responder.path, torrent))
		return nil
	}

	bot.QuickSend(message.Chat.ID, fmt.Sprintf("Cannot find download with hash %s", hash))
	return nil
}

func (responder statusHandler) HelpMessage() string {
	return "Shows status of your downloads"
}

package main

import (
	"fmt"
	"regexp"
	"strings"
	"gitlab.com/NicoNex/echotron"
)


type bot struct {
	chatId int64
	echotron.Engine
}


func NewBot(engine echotron.Engine, chatId int64) echotron.Bot {
	var bot = &bot{
		chatId,
		engine,
	}

	echotron.AddTimer(chatId, "selfDestruct", bot.selfDestruct, 60)
	return bot
}


func (b *bot) Update(update *echotron.Update) {
	if update.Message.Text == "/start" {
		b.SendMessageOptions("Welcome to *Subredditron*!\nSend me any message with a subreddit in the format `r/subreddit` or `/r/subreddit` and I'll send you a link for that subreddit.", b.chatId, echotron.PARSE_MARKDOWN)
	} else if strings.Index(update.Message.Text, "r/") != -1 && strings.Index(update.Message.Text, "reddit.com") == -1 {
		b.SendMessageReply(subreddit(update.Message.Text), b.chatId, update.Message.ID)
		echotron.ResetTimer(b.chatId, "selfDestruct")
	}
}


func (b bot) selfDestruct() {
	echotron.DelTimer(b.chatId, "selfDestruct")
	echotron.DelSession(b.chatId)
}


func subreddit(message string) string {
	re := regexp.MustCompile(`r\/[a-zA-Z]*`)
	sub := re.FindString(message)

	return fmt.Sprintf("https://www.reddit.com/%s", sub)
}


func main() {
	echotron.RunDispatcher("TOKEN", NewBot)
}

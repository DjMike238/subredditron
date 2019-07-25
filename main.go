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
		go echotron.ResetTimer(b.chatId, "selfDestruct")
		sub := subreddit(update.Message.Text)

		if len(sub) > 0 {
			b.SendMessageReply(sub, b.chatId, update.Message.ID)
		}
	}
}


func (b bot) selfDestruct() {
	echotron.DelTimer(b.chatId, "selfDestruct")
	echotron.DelSession(b.chatId)
}


func subreddit(message string) string {
	re := regexp.MustCompile(`(^|[ /])r\/[a-zA-Z_0-9]*`)
	sub := re.FindString(message)
	url := ""

	// Check if the matched string is longer than len(" r/") = 3
	if len(sub) > 3 {
		if sub[:2] == "r/" {
			url = fmt.Sprintf("https://www.reddit.com/%s", sub)
		} else {
			url = fmt.Sprintf("https://www.reddit.com/%s", sub[1:])
		}
	}

	return url
}


func main() {
	echotron.RunDispatcher("TOKEN", NewBot)
}

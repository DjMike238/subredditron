package main

import (
	"fmt"
	"regexp"
	"strings"
	"gitlab.com/NicoNex/echotron"
)


type bot struct {
	chatId int64
	echotron.Api
}


const BOT_NAME = "Subredditron"

var dsp echotron.Dispatcher


func newBot(api echotron.Api, chatId int64) echotron.Bot {
	var bot = &bot{
		chatId,
		api,
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
	dsp.DelSession(b.chatId)
}


func subreddit(message string) string {
	re := regexp.MustCompile(`(^|[ /])r\/[a-zA-Z_0-9]*`)
	sub := re.FindString(message)
	url := ""

	// Check if the matched string is longer than the minimum length for a subreddit
	// name (which is 3) and shorter than the maximum length for a subreddit name
	// (which is 21), both also counting "r/" or "*r/", where * is a character
	// that can be a space (" ") or a slash ("/").
	if len(sub) >= 5 && len(sub) <= 23 && sub[:2] == "r/" {
		url = fmt.Sprintf("https://www.reddit.com/%s", sub)
	} else if len(sub) >= 6 && len(sub) <= 24 && sub[1:3] == "r/" {
		url = fmt.Sprintf("https://www.reddit.com/%s", sub[1:])
	}

	return url
}


func main() {
	dsp = echotron.NewDispatcher("983378957:AAGkoJoydcNsvbHIxU2KGy1ieR1cnDHPnU8", newBot)
	dsp.Run()
}

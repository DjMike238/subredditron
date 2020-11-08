package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"regexp"
	"strings"
	"net/http"
	"github.com/NicoNex/echotron"
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
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error:", err)
			log.Println("Thread recovered. Crysis averted.")
		}
	}()

	if update.Message.Text == "/start" {
		b.SendMessageOptions("Welcome to *Subredditron*!\nSend me any message with a subreddit in the format `r/subreddit` or `/r/subreddit` and I'll send you a link for that subreddit.", b.chatId, echotron.PARSE_MARKDOWN)

	} else if update.Message.Text != "" || update.Message.Caption != "" {
		go echotron.ResetTimer(b.chatId, "selfDestruct")

		var msg string

		if update.Message.Text != "" {
			msg = update.Message.Text
		} else if update.Message.Caption != "" {
			msg = update.Message.Caption
		}

		var sub string

		if strings.Index(msg, "r/") != -1 && strings.Index(msg, "reddit.com") == -1 {
			sub = subreddit(msg)
		}

		var response *http.Response

		if sub != "" {
			b.SendChatAction(echotron.TYPING, b.chatId)
			response, _ = http.Get(sub)
			defer response.Body.Close()

			if response.Status == "404 Not Found" {
				resp := b.SendMessageReply("Subreddit not found.\nThis message will self-destruct in a few seconds.", b.chatId, update.Message.ID)
				time.Sleep(3 * time.Second)
				b.DeleteMessage(b.chatId, resp.Result.ID)
			} else {
				b.SendMessageReply(sub, b.chatId, update.Message.ID)
			}
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
	var url string

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
	_, err := os.Stat(fmt.Sprintf("%s/.log", os.Getenv("HOME")))
	if os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("%s/.log", os.Getenv("HOME")), 0755)
	}
	logfile, err := os.OpenFile(fmt.Sprintf("%s/.log/%s.log", os.Getenv("HOME"), BOT_NAME), os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer logfile.Close()

	log.SetOutput(logfile)
	log.Println(fmt.Sprintf("%s started.", BOT_NAME))
	defer log.Println(fmt.Sprintf("%s stopped.", BOT_NAME))

	dsp = echotron.NewDispatcher("TOKEN", newBot)
	dsp.Run()
}

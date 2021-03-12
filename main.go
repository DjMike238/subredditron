package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/NicoNex/echotron/v2"
)

type bot struct {
	chatID int64
	echotron.API
}

const (
	botName = "Subredditron"
	token   = "token"
)

var dsp *echotron.Dispatcher

func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID,
		echotron.NewAPI(token),
	}
}

func (b *bot) handleMsg(id int, msg string) {
	if strings.HasPrefix(msg, "/start") {
		b.SendMessage(
			"Welcome to *Subredditron*!\nSend me any message with a subreddit in the format `r/subreddit` or `/r/subreddit` and I'll send you a link for that subreddit.",
			b.chatID,
			echotron.ParseMarkdown,
		)

	} else if msg != "" {
		var sub string

		if checkMsg(msg) {
			sub = getSub(msg)
		}

		if sub != "" {
			status := getStatus(sub)

			b.SendChatAction(echotron.Typing, b.chatID)

			if status == 404 {
				resp, err := b.SendMessageReply(
					"Subreddit not found.\nThis message will self-destruct in a few seconds.",
					b.chatID,
					id,
				)
				if err != nil {
					log.Println(err)
				}
				time.Sleep(3 * time.Second)
				b.DeleteMessage(b.chatID, resp.Result.ID)
			} else {
				b.SendMessageReply(sub, b.chatID, id)
			}
		}
	}
}

func (b *bot) handleInline(id, query string) {
	var msg string
	var sub string

	if checkMsg(query) {
		msg = query
	} else if query != "" {
		msg = fmt.Sprintf("r/%s", query)
	}

	if msg != "" {
		sub = getSub(msg)
	}

	if sub != "" {
		status := getStatus(sub)

		if status == 404 {
			_, err := b.AnswerInlineQueryOptions(
				id,
				[]echotron.InlineQueryResult{},
				echotron.InlineQueryOptions{
					CacheTime:         300,
					SwitchPmText:      "Subreddit not found! Try again.",
					SwitchPmParameter: "start",
				},
			)
			if err != nil {
				log.Println(err)
			}
		} else {
			title, desc, thumb := getPreview(msg)

			_, err := b.AnswerInlineQuery(
				id,
				[]echotron.InlineQueryResult{
					&echotron.InlineQueryResultArticle{
						Type:        echotron.ARTICLE,
						ID:          msg,
						Title:       title,
						Description: desc,
						ThumbURL:    thumb,
						InputMessageContent: echotron.InputTextMessageContent{
							MessageText: sub,
						},
					},
				},
			)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (b *bot) Update(update *echotron.Update) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			log.Println("Thread recovered. Crysis averted.")
		}
	}()

	if update.Message != nil {
		b.handleMsg(update.Message.ID, extractMsg(update.Message))
	} else if update.InlineQuery != nil {
		b.handleInline(update.InlineQuery.ID, update.InlineQuery.Query)
	}
}

func main() {
	_, err := os.Stat(fmt.Sprintf("%s/.log", os.Getenv("HOME")))
	if os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("%s/.log", os.Getenv("HOME")), 0755)
	}
	logfile, err := os.OpenFile(fmt.Sprintf("%s/.log/%s.log", os.Getenv("HOME"), botName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer logfile.Close()

	log.SetOutput(logfile)
	log.Println(fmt.Sprintf("%s started.", botName))
	defer log.Println(fmt.Sprintf("%s stopped.", botName))

	dsp = echotron.NewDispatcher(token, newBot)
	dsp.Poll()
}

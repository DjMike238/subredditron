package main

import (
	"fmt"
	"log"
	"os"
	"path"
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
	welcome = `Welcome to *Subredditron*!
Send me any message with a subreddit in the format ` + "`r/subreddit` or `/r/subreddit`" + `and I'll send you a link for that subreddit.

Created by @Dj\_Mike238.
This bot is [open source](https://github.com/DjMike238/subredditron)!`
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
			welcome,
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
			data := getPreview(msg)

			_, err := b.AnswerInlineQuery(
				id,
				[]echotron.InlineQueryResult{
					&echotron.InlineQueryResultArticle{
						Type:        echotron.ARTICLE,
						ID:          getName(data),
						Title:       getTitle(data),
						Description: getDesc(data),
						ThumbURL:    getThumb(data),
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
	logPath := path.Join(os.Getenv("HOME"), ".log")
	_, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		os.Mkdir(logPath, 0755)
	}
	logFile, err := os.OpenFile(path.Join(os.Getenv("HOME"), ".log", botName + ".log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Println(fmt.Sprintf("%s started.", botName))
	defer log.Println(fmt.Sprintf("%s stopped.", botName))

	dsp = echotron.NewDispatcher(token, newBot)
	dsp.Poll()
}

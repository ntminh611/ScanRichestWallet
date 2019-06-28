package service

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"sync"
	"time"
)

type Telegram struct {
	bot     tgbotapi.BotAPI
	token   string
	chatTo  int64
	message chan string
	mutex   sync.Mutex
}

func NewTelegramBot(token string, chatTo int64) *Telegram {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	return &Telegram{bot: *bot, token: token, chatTo: chatTo}
}

func (self *Telegram) sendMessage(msg string) {
	m := tgbotapi.NewMessage(self.chatTo, msg)
	m.ParseMode = "HTML"
	m.DisableWebPagePreview = true
	_, e := self.bot.Send(m)
	if e != nil {
		fmt.Println(e)
	}
}

func (self *Telegram) SendMessage(message string) {
	self.mutex.Lock()
	color.Blue("Receive message")
	self.sendMessage(message)
	color.Green("Sent!")
	time.Sleep(time.Millisecond * 100)
	self.mutex.Unlock()
}

/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package telegram

import (
	"fmt"
	"log"
	"os"
	"time"

	// "gopkg.in/telegram-bot-api.v4"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/bitnami-labs/kubewatch/config"
	kbEvent "github.com/bitnami-labs/kubewatch/pkg/event"
)

var telegramErrMsg = `
%s

You need to set both token and channel for slack notify,
using "--token/-t" and "--channel/-c", or using environment variables:

export KW_TELEGRAM_TOKEN=telegram_token
export KW_TELEGRAM_CHANNEL=telegram_channel

Command line flags will override environment variables

`

// Telegram handler implements handler.Handler interface,
// Notify event to telegram channel
type Telegram struct {
	Token   string
	Channel string
}

type TelegramMessage struct {
	Text string `json:"text"`
}

// Init prepares slack configuration
func (s *Telegram) Init(c *config.Config) error {
	token := c.Handler.Telegram.Token
	channel := c.Handler.Telegram.Channel

	if token == "" {
		token = os.Getenv("KW_TELEGRAM_TOKEN")
	}

	if channel == "" {
		channel = os.Getenv("KW_TELEGRAM_CHANNEL")
	}

	s.Token = token
	s.Channel = channel

	return checkMissingTelegramVars(s)
}

func (s *Telegram) ObjectCreated(obj interface{}) {
	notifyTelegram(s, obj, "created")
}

func (s *Telegram) ObjectDeleted(obj interface{}) {
	notifyTelegram(s, obj, "deleted")
}

func (s *Telegram) ObjectUpdated(oldObj, newObj interface{}) {
	notifyTelegram(s, newObj, "updated")
}

func notifyTelegram(t *Telegram, obj interface{}, action string) {
	e := kbEvent.New(obj, action)
	b, err := tb.NewBot(tb.Settings{
		Token:  t.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	stChannel := &tb.Chat{
		Type:     tb.ChatChannel,
		Username: t.Channel,
	}
	var opts tb.SendOptions
	opts.ParseMode = tb.ModeMarkdown

	b.Send(stChannel, e.Message(), &opts)
	log.Printf("Message successfully sent to channel %s", t.Channel)
}

func checkMissingTelegramVars(t *Telegram) error {
	if t.Token == "" || t.Channel == "" {
		return fmt.Errorf(telegramErrMsg, "Missing telegram token or channel")
	}

	return nil
}

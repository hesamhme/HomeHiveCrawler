package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const token = "8096414698:AAH7ox_ZZkUNSF0fO82kbKsKe3g514Uo2Mc"

func SetupBot() error {
	bot, updateConfig, err2 := initializeBot()

	if err2 != nil {
		return err2
	}

	runBot(bot, updateConfig)

	return nil
}

func initializeBot() (*tgbotapi.BotAPI, tgbotapi.UpdateConfig, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, tgbotapi.UpdateConfig{}, err
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	return bot, updateConfig, nil
}

func runBot(bot *tgbotapi.BotAPI, updateConfig tgbotapi.UpdateConfig) {
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if msg.Text != "" {
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}

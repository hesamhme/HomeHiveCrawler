package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var filterKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("بازه قیمت"),
		tgbotapi.NewKeyboardButton("شهر"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("محله"),
		tgbotapi.NewKeyboardButton("بازه متراژ"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("بازه تعداد اتاق خواب"),
		tgbotapi.NewKeyboardButton("دسته‌بندی اجاره، خرید، رهن"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("رنج سن بنا"),
		tgbotapi.NewKeyboardButton("دسته‌بندی آپارتمانی یا ویلایی"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("بازه طبقه (در صورت آپارتمانی بودن)"),
		tgbotapi.NewKeyboardButton("داشتن انباری"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("داشتن آسانسور"),
		tgbotapi.NewKeyboardButton("بازه تاریخ ایجاد آگهی"),
	), tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("بازگشت"),
	),
)

func SetupBot(token string) error {
	bot, updateConfig, err2 := initializeBot(token)

	if err2 != nil {
		return err2
	}

	runBot(bot, updateConfig)

	return nil
}

func initializeBot(token string) (*tgbotapi.BotAPI, tgbotapi.UpdateConfig, error) {
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
	var lastUserMessageID int
	var lastBotMessageID int

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go func() {
			if lastUserMessageID != 0 {
				deleteUserMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, lastUserMessageID)
				if _, err := bot.Send(deleteUserMsg); err != nil {
					log.Printf("Failed to delete user command: %v", err)
				}
			}
		}()
		go func() {
			if lastBotMessageID != 0 {
				deleteBotMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, lastBotMessageID)
				if _, err := bot.Send(deleteBotMsg); err != nil {
					log.Printf("Failed to delete bot response: %v", err)
				}
			}
		}()

		lastUserMessageID = update.Message.MessageID

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Text {
		case "/start":
			lastBotMessageID = handleStart(bot, &update)
		case "/help":
			lastBotMessageID = handleHelp(bot, &update)
		case "/account":
			handleAccount(bot, &update)
		case "/search":
			msg.Text = "Filter menu opened."
			msg.ReplyMarkup = filterKeyboard
		case "بازگشت":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "بازه قیمت":
			handlePriceRange(bot, update.Message)
		case "شهر":
			handleCity(bot, update.Message)
		case "محله":
			handleNeighborhood(bot, update.Message)
		case "بازه متراژ":
			handleAreaRange(bot, update.Message)
		case "بازه تعداد اتاق خواب":
			handleBedroomCount(bot, update.Message)
		case "دسته‌بندی اجاره، خرید، رهن":
			handleRentBuyMortgage(bot, update.Message)
		case "رنج سن بنا":
			handleBuildingAge(bot, update.Message)
		case "دسته‌بندی آپارتمانی یا ویلایی":
			handleBuildingType(bot, update.Message)
		case "بازه طبقه (در صورت آپارتمانی بودن)":
			handleFloorRange(bot, update.Message)
		case "داشتن انباری":
			handleStorage(bot, update.Message)
		case "داشتن آسانسور":
			handleElevator(bot, update.Message)
		case "بازه تاریخ ایجاد آگهی":
			handleAdCreationDate(bot, update.Message)
		default:
			msg.Text = "Command not recognized."
		}

		if msg.Text != "" {
			sentMsg, err := bot.Send(msg)
			if err != nil {
				log.Panic(err)
			}
			lastBotMessageID = sentMsg.MessageID
		}
	}
}

func handleStart(bot *tgbotapi.BotAPI, update *tgbotapi.Update) int {
	//TODO create user
	username := update.Message.From.UserName
	message := fmt.Sprintf("سلام %s عزیز، خوش آمدی", username)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	sentMsg, err := bot.Send(msg)
	if err != nil {
		return -1
	}
	return sentMsg.MessageID
}
func handleHelp(bot *tgbotapi.BotAPI, update *tgbotapi.Update) int {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "در حال حاضر راهنمایی وجود نداره، سعی کن رو پای خودت وایسی")
	sentMsg, err := bot.Send(msg)
	if err != nil {
		return -1
	}
	return sentMsg.MessageID
}
func handleAccount(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	//TODO show user account data
}

func handlePriceRange(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	//TODO filter based on price
}

func handleCity(bot *tgbotapi.BotAPI, update *tgbotapi.Message) {
	//TODO filter based on city
}

func handleNeighborhood(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	//TODO filter based on neighborhood
}

func handleAreaRange(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	//TODO filter based on area range
}

func handleBedroomCount(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	//TODO filter based on bedroom count
}

func handleRentBuyMortgage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	//TODO filter based on mortgage
}

func handleBuildingAge(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	//TODO filter based on building age
}

func handleBuildingType(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	//TODO filter based on building type
}

func handleFloorRange(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	//TODO filter based on floor range
}

func handleStorage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	//TODO filter based on storage
}

func handleElevator(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	//TODO filter based on elevator
}

func handleAdCreationDate(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	//TODO filter based on creation date
}

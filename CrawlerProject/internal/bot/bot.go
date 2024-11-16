package bot

import (
	"CrawlerProject/internal/model"
	"CrawlerProject/internal/service"
	"fmt"
	"log"
	"strings"
	"strconv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var db *gorm.DB // Global db variable
var userState = make(map[int64]string)

func SetDB(database *gorm.DB) {

	db = database
}

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
		tgbotapi.NewKeyboardButton("اجاره، خرید، رهن"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("رنج سن بنا"),
		tgbotapi.NewKeyboardButton(" آپارتمانی یا ویلایی"),
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
		chatID := update.Message.Chat.ID

		// Check if user is in a specific state
		if state, exists := userState[chatID]; exists {
			if state == "awaiting_city_input" {
				handleCitySearch(bot, update.Message, db) // Assuming `db` is your GORM DB instance
				delete(userState, chatID)                 // Clear user state after handling
				continue
			}
			if state == "awaiting_price_range_input" {
				handlePriceRangeSearch(bot, update.Message, db) // Assuming `db` is your GORM DB instance
				delete(userState, chatID)                 // Clear user state after handling
				continue
			}
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
		case " اجاره، خرید، رهن":
			handleRentBuyMortgage(bot, update.Message)
		case "رنج سن بنا":
			handleBuildingAge(bot, update.Message)
		case " آپارتمانی یا ویلایی":
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
	username := update.Message.From.UserName
	telegramID := update.Message.From.ID

	// Check if the user already exists in the database
	var user model.User
	err := db.Where("telegram_id = ?", telegramID).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		// User does not exist, create a new entry
		newUser := model.User{
			TelegramID: telegramID,
			Username:   username,
		}
		if err := db.Create(&newUser).Error; err != nil {
			log.Printf("Failed to create user: %v", err)
		} else {
			log.Printf("New user created: %s", username)
		}
	}

	// Send welcome message
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
	telegramID := update.Message.From.ID
	var user model.User
	err := db.Where("telegram_id = ?", telegramID).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "حسابی برای شما یافت نشد.")
		bot.Send(msg)
		return
	} else if err != nil {
		log.Printf("Error fetching user account: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "خطا در بازیابی اطلاعات حساب.")
		bot.Send(msg)
		return
	}

	accountInfo := fmt.Sprintf("نام کاربری: %s\nشناسه تلگرام: %d\n", user.Username, user.TelegramID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, accountInfo)
	bot.Send(msg)
}

func handlePriceRange(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    msg := tgbotapi.NewMessage(message.Chat.ID, "لطفا رنج قیمت مورد نظر خود را به صورت (شروع,پایان) وارد نمایید")
    bot.Send(msg)

    // Set user state to expect price range input (assuming you are using userState map)
    userState[message.Chat.ID] = "awaiting_price_range_input"
}

func handlePriceRangeSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
    // Parse input to extract minimum and maximum price
    var minPrice, maxPrice float64
    input := strings.Split(message.Text, ",")
    if len(input) != 2 {
        msg := tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفا رنج قیمت را به صورت (شروع,پایان) وارد نمایید")
        bot.Send(msg)
        return
    }
    
    // Attempt to parse the input values to floats
    var err error
    minPrice, err = strconv.ParseFloat(strings.TrimSpace(input[0]), 64)
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, "مقدار شروع معتبر نیست.")
        bot.Send(msg)
        return
    }
    maxPrice, err = strconv.ParseFloat(strings.TrimSpace(input[1]), 64)
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, "مقدار پایان معتبر نیست.")
        bot.Send(msg)
        return
    }

    // Create filters for price range
    filters := model.Filter{
        PriceMin: minPrice,
        PriceMax: maxPrice,
    }

    // Execute the search query
    results, err := service.GetFilteredListings(db, filters)
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
        bot.Send(msg)
        return
    }

    if len(results) == 0 {
        msg := tgbotapi.NewMessage(message.Chat.ID, "هیچ نتیجه ای برای رنج قیمت مورد نظر پیدا نشد")
        bot.Send(msg)
        return
    }

    // Send the results
	sendFormattedListings(bot, message.Chat.ID, results)
}

func handleCity(bot *tgbotapi.BotAPI, update *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(update.Chat.ID, "لطفاً نام شهر مورد نظر خود را وارد کنید:")
	bot.Send(msg)

	// Set user state to expect city input (assuming you are using userState map)
	userState[update.Chat.ID] = "awaiting_city_input"
}

func handleCitySearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	city := message.Text
	filters := model.Filter{
		City: city,
	}

	results, err := service.GetFilteredListings(db, filters)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
		bot.Send(msg)
		return
	}

	if len(results) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("شهر %s موجود نیست", city))
		bot.Send(msg)
		return
	}

    sendFormattedListings(bot, message.Chat.ID, results)
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

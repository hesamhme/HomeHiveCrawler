package bot

import (
	"CrawlerProject/internal/model"
	"CrawlerProject/internal/service"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"github.com/jung-kurt/gofpdf/v2"
	"github.com/xuri/excelize/v2"
	"os"
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
				handleCitySearch(bot, update.Message, db)
				delete(userState, chatID)
				continue
			}
			if state == "awaiting_price_range_input" {
				handlePriceRangeSearch(bot, update.Message, db)
				delete(userState, chatID)
				continue
			}
			// handleNeighborhoodSearch
			if state == "awaiting_neighborhood_input" {
				handleNeighborhoodSearch(bot, update.Message, db)
				delete(userState, chatID)
				continue
			}
			// handleAreaRangeSearch
			if state == "awaiting_area_range_input" {
				handleAreaRangeSearch(bot, update.Message, db)
				delete(userState, chatID)
				continue
			}
			// handleBedroomCountSearch
			if state == "awaiting_bedroom_count_input" {
				handleBedroomCountSearch(bot, update.Message, db)
				delete(userState, chatID)
				continue
			}

			if state == "awaiting_building_age_input" {
				handleBuildingAgeSearch(bot, update.Message, db)
				delete(userState, chatID)
				continue
			}
			if state == "awaiting_building_type_input" {
				handleBuildingTypeSearch(bot, update.Message, db)
				delete(userState, chatID)
				continue
			}
			if state == "awaiting_floor_range_input" {
				handleFloorRangeSearch(bot, update.Message, db)
				delete(userState, chatID)
				continue
			}
			if state == "awaiting_storage_input" {
				handleStorageSearch(bot, update.Message, db)
				delete(userState, chatID)
				continue
			}
			if state == "awaiting_elevator_input" {
				handleElevatorSearch(bot, update.Message, db)
				delete(userState, chatID)
				continue
			}
			if state == "awaiting_ad_creation_date_input" {
				handleAdCreationDateSearch(bot, update.Message, db)
				delete(userState, chatID)
				continue
			}
			if state == "awaiting_rent_buy_mortgage_input" {
				handleRentBuyMortgageSearch(bot, update.Message, db)
				delete(userState, chatID)
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
		case "اجاره، خرید، رهن":
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
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً نام محله مورد نظر خود را وارد کنید:")
	bot.Send(msg)

	// Set user state to expect neighborhood input (assuming you are using userState map)
	userState[message.Chat.ID] = "awaiting_neighborhood_input"
}

func handleNeighborhoodSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	// Extract the neighborhood name from the user's message
	neighborhood := strings.TrimSpace(message.Text)

	// Create a filter for the search query
	filters := model.Filter{
		Neighborhood: neighborhood,
	}

	// Query the database using the specified filter
	results, err := service.GetFilteredListings(db, filters)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
		bot.Send(msg)
		return
	}

	if len(results) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("هیچ نتیجه ای برای محله %s یافت نشد", neighborhood))
		bot.Send(msg)
		return
	}

	// Use the helper function to format and send results
	sendFormattedListings(bot, message.Chat.ID, results)
}

func handleAreaRange(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً محدوده مساحت مورد نظر خود را به صورت (شروع,پایان) وارد نمایید:")
	bot.Send(msg)

	// Set user state to expect area range input (assuming you are using a userState map)
	userState[message.Chat.ID] = "awaiting_area_range_input"
}

func handleAreaRangeSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	// Parse input to extract minimum and maximum area
	var minArea, maxArea float64
	input := strings.Split(message.Text, ",")
	if len(input) != 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفا محدوده مساحت را به صورت (شروع,پایان) وارد نمایید")
		bot.Send(msg)
		return
	}

	// Attempt to parse the input values to floats
	var err error
	minArea, err = strconv.ParseFloat(strings.TrimSpace(input[0]), 64)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "مقدار شروع معتبر نیست.")
		bot.Send(msg)
		return
	}
	maxArea, err = strconv.ParseFloat(strings.TrimSpace(input[1]), 64)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "مقدار پایان معتبر نیست.")
		bot.Send(msg)
		return
	}

	// Create filters for area range
	filters := model.Filter{
		AreaMin: minArea,
		AreaMax: maxArea,
	}

	// Execute the search query
	results, err := service.GetFilteredListings(db, filters)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
		bot.Send(msg)
		return
	}

	if len(results) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "هیچ نتیجه ای برای محدوده مساحت مورد نظر یافت نشد")
		bot.Send(msg)
		return
	}

	// Use the helper function to format and send results
	sendFormattedListings(bot, message.Chat.ID, results)
}

func handleBedroomCount(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً تعداد اتاق خواب مورد نظر خود را به صورت (حداقل,حداکثر) وارد نمایید:")
	bot.Send(msg)

	// Set user state to expect bedroom count input (assuming you are using a userState map)
	userState[message.Chat.ID] = "awaiting_bedroom_count_input"
}

func handleBedroomCountSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	// Parse input to extract minimum and maximum bedroom count
	var minRooms, maxRooms int
	input := strings.Split(message.Text, ",")
	if len(input) != 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفا تعداد اتاق خواب را به صورت (حداقل,حداکثر) وارد نمایید")
		bot.Send(msg)
		return
	}

	// Attempt to parse the input values to integers
	var err error
	minRooms, err = strconv.Atoi(strings.TrimSpace(input[0]))
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "مقدار حداقل تعداد اتاق معتبر نیست.")
		bot.Send(msg)
		return
	}
	maxRooms, err = strconv.Atoi(strings.TrimSpace(input[1]))
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "مقدار حداکثر تعداد اتاق معتبر نیست.")
		bot.Send(msg)
		return
	}

	// Create filters for bedroom count
	filters := model.Filter{
		RoomsMin: minRooms,
		RoomsMax: maxRooms,
	}

	// Execute the search query
	results, err := service.GetFilteredListings(db, filters)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
		bot.Send(msg)
		return
	}

	if len(results) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "هیچ نتیجه ای برای محدوده تعداد اتاق خواب مورد نظر یافت نشد")
		bot.Send(msg)
		return
	}

	// Use the helper function to format and send results
	sendFormattedListings(bot, message.Chat.ID, results)
}

func handleRentBuyMortgage(bot *tgbotapi.BotAPI, update *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(update.Chat.ID, "لطفاً نوع مورد نظر خود را وارد کنید: خرید/ فروش / اجاره / رهن / اجاره و رهن")
	bot.Send(msg)
	userState[update.Chat.ID] = "awaiting_rent_buy_mortgage_input"
}

func handleRentBuyMortgageSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	// Extract the user's input and trim whitespace
	input := strings.TrimSpace(message.Text)

	// Create filters for status (ad type)
	filters := model.Filter{}

	// Special handling for "اجاره و رهن"
	if input == "اجاره و رهن" {
		filters.Status = "اجاره و رهن" // Special case identifier
	} else {
		filters.Status = input // Direct assignment for other cases
	}

	// Execute the search query
	results, err := service.GetFilteredListings(db, filters)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
		bot.Send(msg)
		return
	}

	if len(results) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("هیچ نتیجه ای برای وضعیت %s یافت نشد", input))
		bot.Send(msg)
		return
	}

	// Use the helper function to format and send results
	sendFormattedListings(bot, message.Chat.ID, results)
}

func handleBuildingAge(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً محدوده سن بنا را به صورت (حداقل,حداکثر) وارد کنید:")
    bot.Send(msg)

    // Set user state to expect building age range input
    userState[message.Chat.ID] = "awaiting_building_age_input"
}

func handleBuildingAgeSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
    input := strings.Split(message.Text, ",")
    if len(input) != 2 {
        msg := tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفاً محدوده سن بنا را به صورت (حداقل,حداکثر) وارد نمایید.")
        bot.Send(msg)
        return
    }

    // Parse the input values for min and max age
    var minAge, maxAge int
    var err error
    minAge, err = strconv.Atoi(strings.TrimSpace(input[0]))
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, "مقدار حداقل سن بنا معتبر نیست.")
        bot.Send(msg)
        return
    }
    maxAge, err = strconv.Atoi(strings.TrimSpace(input[1]))
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, "مقدار حداکثر سن بنا معتبر نیست.")
        bot.Send(msg)
        return
    }

    // Convert ages to match the formula (age = 1403 - age)
    currentYear := 1403
    filters := model.Filter{
        BuildingAgeMin: currentYear - maxAge,
        BuildingAgeMax: currentYear - minAge,
    }

    results, err := service.GetFilteredListings(db, filters)
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
        bot.Send(msg)
        return
    }

    if len(results) == 0 {
        msg := tgbotapi.NewMessage(message.Chat.ID, "هیچ نتیجه ای برای محدوده سن بنا مورد نظر یافت نشد")
        bot.Send(msg)
        return
    }

    sendFormattedListings(bot, message.Chat.ID, results)
}


func handleBuildingType(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    msg := tgbotapi.NewMessage(message.Chat.ID, "  لطفاً نوع ملک مورد نظر خود را وارد کنید آپارتمان/ویلایی/غیره:")
    bot.Send(msg)

    // Set user state to expect building type input
    userState[message.Chat.ID] = "awaiting_building_type_input"
}

func handleBuildingTypeSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
    propertyType := strings.TrimSpace(message.Text)

    filters := model.Filter{
        PropertyType: propertyType,
    }

    results, err := service.GetFilteredListings(db, filters)
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
        bot.Send(msg)
        return
    }

    if len(results) == 0 {
        msg := tgbotapi.NewMessage(message.Chat.ID, "هیچ نتیجه ای برای نوع ملک وارد شده یافت نشد")
        bot.Send(msg)
        return
    }

    sendFormattedListings(bot, message.Chat.ID, results)
}



func handleFloorRange(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً محدوده طبقه را به صورت (حداقل,حداکثر) وارد کنید:")
    bot.Send(msg)

    // Set user state to expect floor range input
    userState[message.Chat.ID] = "awaiting_floor_range_input"
}

func handleFloorRangeSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
    input := strings.Split(message.Text, ",")
    if len(input) != 2 {
        msg := tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفاً محدوده طبقه را به صورت (حداقل,حداکثر) وارد نمایید.")
        bot.Send(msg)
        return
    }

    var minFloor, maxFloor int
    var err error
    minFloor, err = strconv.Atoi(strings.TrimSpace(input[0]))
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, "مقدار حداقل طبقه معتبر نیست.")
        bot.Send(msg)
        return
    }
    maxFloor, err = strconv.Atoi(strings.TrimSpace(input[1]))
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, "مقدار حداکثر طبقه معتبر نیست.")
        bot.Send(msg)
        return
    }

    filters := model.Filter{
        FloorMin: &minFloor,
        FloorMax: &maxFloor,
    }

    results, err := service.GetFilteredListings(db, filters)
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
        bot.Send(msg)
        return
    }

    if len(results) == 0 {
        msg := tgbotapi.NewMessage(message.Chat.ID, "هیچ نتیجه ای برای محدوده طبقه مورد نظر یافت نشد")
        bot.Send(msg)
        return
    }

    sendFormattedListings(bot, message.Chat.ID, results)
}


func handleStorage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    msg := tgbotapi.NewMessage(message.Chat.ID, "آیا ملک دارای انباری باشد؟ (بله/خیر)")
    bot.Send(msg)

    // Set user state to expect storage preference input
    userState[message.Chat.ID] = "awaiting_storage_input"
}

func handleStorageSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
    input := strings.TrimSpace(message.Text)

    // Convert input to a boolean
    var hasStorage bool
    if input == "بله" {
        hasStorage = true
    } else if input == "خیر" {
        hasStorage = false
    } else {
        msg := tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفاً 'بله' یا 'خیر' وارد کنید.")
        bot.Send(msg)
        return
    }

    filters := model.Filter{
        HasStorage: hasStorage,
    }

    results, err := service.GetFilteredListings(db, filters)
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
        bot.Send(msg)
        return
    }

    if len(results) == 0 {
        msg := tgbotapi.NewMessage(message.Chat.ID, "هیچ نتیجه ای برای وضعیت انباری مورد نظر یافت نشد")
        bot.Send(msg)
        return
    }

    sendFormattedListings(bot, message.Chat.ID, results)
}


func handleElevator(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    msg := tgbotapi.NewMessage(message.Chat.ID, "آیا ملک دارای آسانسور باشد؟ (بله/خیر)")
    bot.Send(msg)

    // Set user state to expect elevator preference input
    userState[message.Chat.ID] = "awaiting_elevator_input"
}

func handleElevatorSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
    input := strings.TrimSpace(message.Text)

    // Convert input to a boolean
    var hasElevator bool
    if input == "بله" {
        hasElevator = true
    } else if input == "خیر" {
        hasElevator = false
    } else {
        msg := tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفاً 'بله' یا 'خیر' وارد کنید.")
        bot.Send(msg)
        return
    }

    filters := model.Filter{
        HasElevator: hasElevator,
    }

    results, err := service.GetFilteredListings(db, filters)
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
        bot.Send(msg)
        return
    }

    if len(results) == 0 {
        msg := tgbotapi.NewMessage(message.Chat.ID, "هیچ نتیجه ای برای وضعیت آسانسور مورد نظر یافت نشد")
        bot.Send(msg)
        return
    }

    sendFormattedListings(bot, message.Chat.ID, results)
}


func handleAdCreationDate(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً محدوده تاریخ درج آگهی را به صورت (شروع,پایان) وارد کنید (YYYY-MM-DD):")
    bot.Send(msg)

    // Set user state to expect ad creation date range input
    userState[message.Chat.ID] = "awaiting_ad_creation_date_input"
}

func handleAdCreationDateSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
    input := strings.Split(message.Text, ",")
    if len(input) != 2 {
        msg := tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفاً محدوده تاریخ را به صورت (شروع,پایان) به فرمت YYYY-MM-DD وارد نمایید.")
        bot.Send(msg)
        return
    }

    // Parse the input dates
    startDate, err := time.Parse("2006-01-02", strings.TrimSpace(input[0]))
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, "تاریخ شروع معتبر نیست. لطفاً به فرمت YYYY-MM-DD وارد نمایید.")
        bot.Send(msg)
        return
    }
    endDate, err := time.Parse("2006-01-02", strings.TrimSpace(input[1]))
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, "تاریخ پایان معتبر نیست. لطفاً به فرمت YYYY-MM-DD وارد نمایید.")
        bot.Send(msg)
        return
    }

    filters := model.Filter{
        CreationDateMin: startDate,
        CreationDateMax: endDate,
    }

    results, err := service.GetFilteredListings(db, filters)
    if err != nil {
        msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در جستجوی اطلاعات: %v", err))
        bot.Send(msg)
        return
    }

    if len(results) == 0 {
        msg := tgbotapi.NewMessage(message.Chat.ID, "هیچ نتیجه ای برای محدوده تاریخ درج آگهی مورد نظر یافت نشد")
        bot.Send(msg)
        return
    }

    sendFormattedListings(bot, message.Chat.ID, results)
}


func ExportListingsToExcel(listings []model.Listing, filename string) error {
	f := excelize.NewFile()
	sheetName := "Listings"
	f.NewSheet(sheetName)

	// Persian headers
	headers := []string{
		"عنوان", "قیمت", "شهر", "محله", "متراژ", "تعداد اتاق خواب", "نوع آگهی", "سن بنا", 
		"نوع ملک", "طبقه", "انباری", "آسانسور", "تاریخ ایجاد", "تاریخ بروزرسانی", "لینک",
	}

	// Add headers to Excel
	for col, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+col)))
		f.SetCellValue(sheetName, cell, header)
	}

	// Add data rows
	for row, listing := range listings {
		rowNum := row + 2 // Start from the second row

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), listing.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), listing.Price)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), listing.City)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), listing.Neighborhood)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), fmt.Sprintf("%.2f متر مربع", listing.Area))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), listing.Rooms)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), listing.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), listing.Age)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), listing.HouseType)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), listing.Floor)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), func() string {
			if listing.HasStorage {
				return "دارد"
			}
			return "ندارد"
		}())
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), func() string {
			if listing.HasElevator {
				return "دارد"
			}
			return "ندارد"
		}())
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", rowNum), listing.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", rowNum), listing.UpdatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", rowNum), listing.URL)
	}

	// Save Excel file
	return f.SaveAs(filename)
}


func ExportListingsToPDF(listings []model.Listing, filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)

	// Title
	pdf.CellFormat(0, 10, "Filtered Listings", "", 1, "C", false, 0, "")

	// Listing Content
	for _, listing := range listings {
		pdf.SetFont("Arial", "", 12)

		// Add fields with proper labels in Persian
		data := []struct {
			Label string
			Value string
		}{
			{"عنوان", listing.Title},
			{"قیمت", fmt.Sprintf("%.2f", listing.Price)},
			{"شهر", listing.City},
			{"محله", listing.Neighborhood},
			{"متراژ", fmt.Sprintf("%.2f متر مربع", listing.Area)},
			{"تعداد اتاق خواب", strconv.Itoa(listing.Rooms)},
			{"نوع آگهی", listing.Status},
			{"سن بنا", listing.Age},
			{"نوع ملک", listing.HouseType},
			{"طبقه", strconv.Itoa(listing.Floor)},
			{"انباری", func() string {
				if listing.HasStorage {
					return "دارد"
				}
				return "ندارد"
			}()},
			{"آسانسور", func() string {
				if listing.HasElevator {
					return "دارد"
				}
				return "ندارد"
			}()},
			{"تاریخ ایجاد", listing.CreatedAt.Format("2006-01-02 15:04:05")},
			{"تاریخ بروزرسانی", listing.UpdatedAt.Format("2006-01-02 15:04:05")},
			{"لینک", listing.URL},
		}

		// Print each field
		for _, field := range data {
			pdf.CellFormat(0, 10, fmt.Sprintf("%s: %s", field.Label, field.Value), "", 1, "", false, 0, "")
		}

		// Add spacing between listings
		pdf.Ln(5)
	}

	// Save PDF to file
	return pdf.OutputFileAndClose(filename)
}


// func uploadFileToTelegram(bot *tgbotapi.BotAPI, chatID int64, filename string) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		log.Printf("Failed to open file: %v", err)
// 		return
// 	}
// 	defer file.Close()

// 	doc := tgbotapi.NewDocument(chatID, tgbotapi.FileReader{
// 		Name:   filename,
// 		Reader: file,
// 	})
// 	_, err = bot.Send(doc)
// 	if err != nil {
// 		log.Printf("Failed to send file: %v", err)
// 	}
// }
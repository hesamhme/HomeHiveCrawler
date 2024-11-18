package bot

import (
	"CrawlerProject/internal/model"
	"CrawlerProject/internal/service"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"archive/zip"
	"io"
	"encoding/csv"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var db *gorm.DB // Global db variable
var userState = make(map[int64]string)
var userFilters = make(map[int64]model.Filter)

const (
	awaitingPriceRange      = "awaiting_price_range_input"
	awaitingCity            = "awaiting_city_input"
	awaitingNeighborhood    = "awaiting_neighborhood_input"
	awaitingAreaRange       = "awaiting_area_range_input"
	awaitingBedroomCount    = "awaiting_bedroom_count_input"
	awaitingBuildingAge     = "awaiting_building_age_input"
	awaitingBuildingType    = "awaiting_building_type_input"
	awaitingFloorRange      = "awaiting_floor_range_input"
	awaitingStorage         = "awaiting_storage_input"
	awaitingElevator        = "awaiting_elevator_input"
	awaitingAdCreationDate  = "awaiting_ad_creation_date_input"
	awaitingRentBuyMortgage = "awaiting_rent_buy_mortgage_input"
)

func SetDB(database *gorm.DB) {

	db = database
}

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/start"),
		tgbotapi.NewKeyboardButton("/account"),
		tgbotapi.NewKeyboardButton("/help"),
		tgbotapi.NewKeyboardButton("/search"), // Search/Filter button
	),
)

// Filter options keyboard (shown during search)
var filterKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("جستوجوی پیشرفته"),
		tgbotapi.NewKeyboardButton("تایید فیلترها"),
	),
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
		tgbotapi.NewKeyboardButton("آپارتمانی یا ویلایی"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("بازه طبقه (در صورت آپارتمانی بودن)"),
		tgbotapi.NewKeyboardButton("داشتن انباری"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("داشتن آسانسور"),
		tgbotapi.NewKeyboardButton("بازه تاریخ ایجاد آگهی"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("دریافت نتایج به صورت فایل CSV"), // Download CSV Button
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

		if state, exists := userState[chatID]; exists {
			handleUserState(bot, update.Message, state)
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
			msg.Text = "لطفاً یک فیلتر را انتخاب کنید."
			msg.ReplyMarkup = filterKeyboard
		case "جستوجوی پیشرفته":
			msg.Text = "لطفاً فیلترهای مورد نظر خود را انتخاب کنید. وقتی آماده شدید، 'تایید فیلترها' را بزنید."
			msg.ReplyMarkup = filterKeyboard
		case "تایید فیلترها":
			handleConfirmFilters(bot, update.Message)
			delete(userFilters, update.Message.Chat.ID)
		// Filter cases for specific criteria
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
		case "آپارتمانی یا ویلایی":
			handleBuildingType(bot, update.Message)
		case "بازه طبقه (در صورت آپارتمانی بودن)":
			handleFloorRange(bot, update.Message)
		case "داشتن انباری":
			handleStorage(bot, update.Message)
		case "داشتن آسانسور":
			handleElevator(bot, update.Message)
		case "بازه تاریخ ایجاد آگهی":
			handleAdCreationDate(bot, update.Message)
		case "دریافت نتایج به صورت فایل CSV":
			handleDownloadCSV(bot, update.Message)
		default:
			msg.Text = "دستور شناسایی نشد. لطفاً یکی از گزینه‌های منو را انتخاب کنید."
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

func handleUserState(bot *tgbotapi.BotAPI, message *tgbotapi.Message, state string) {
	switch state {
	case awaitingPriceRange:
		handlePriceRangeSearch(bot, message, db)
	case awaitingCity:
		handleCitySearch(bot, message, db)
	case awaitingNeighborhood:
		handleNeighborhoodSearch(bot, message, db)
	case awaitingAreaRange:
		handleAreaRangeSearch(bot, message, db)
	case awaitingBedroomCount:
		handleBedroomCountSearch(bot, message, db)
	case awaitingBuildingAge:
		handleBuildingAgeSearch(bot, message, db)
	case awaitingBuildingType:
		handleBuildingTypeSearch(bot, message, db)
	case awaitingFloorRange:
		handleFloorRangeSearch(bot, message, db)
	case awaitingStorage:
		handleStorageSearch(bot, message, db)
	case awaitingElevator:
		handleElevatorSearch(bot, message, db)
	case awaitingAdCreationDate:
		handleAdCreationDateSearch(bot, message, db)
	case awaitingRentBuyMortgage:
		handleRentBuyMortgageSearch(bot, message, db)
	default:
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "حالت شناسایی نشد."))
	}
	delete(userState, message.Chat.ID)
}

func handleStart(bot *tgbotapi.BotAPI, update *tgbotapi.Update) int {
	username := update.Message.From.UserName
	telegramID := update.Message.From.ID

	var user model.User
	err := db.Where("telegram_id = ?", telegramID).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		newUser := model.User{
			TelegramID: telegramID,
			Username:   username,
		}
		if err := db.Create(&newUser).Error; err != nil {
			log.Printf("Failed to create user: %v", err)
		}
	}

	message := fmt.Sprintf("سلام %s عزیز، خوش آمدی", username)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	sentMsg, err := bot.Send(msg)
	if err != nil {
		return -1
	}
	return sentMsg.MessageID
}

// A map to store filtered results per user
var userFilteredResults = make(map[int64][]model.Listing)

// Example function where you confirm filters and generate results
func handleConfirmFilters(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    // Assuming you generate filteredResults based on user filters
    filteredResults, err := service.GetFilteredListings(db, userFilters[message.Chat.ID])
    if err != nil {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, "خطا در جستجوی اطلاعات"))
        return
    }

    if len(filteredResults) == 0 {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, "هیچ نتیجه‌ای برای فیلترهای انتخاب‌شده یافت نشد."))
        return
    }

    // Store the filtered results for this user
    userFilteredResults[message.Chat.ID] = filteredResults

    // Display the results (implementation may vary)
    sendFormattedListings(bot, message.Chat.ID, filteredResults)

    // Optionally, prompt user for further actions
    bot.Send(tgbotapi.NewMessage(message.Chat.ID, "برای دانلود نتایج به صورت فایل CSV دکمه مربوطه را فشار دهید."))
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
    //parseRangeInput 
    minPrice, maxPrice, err := parseRangeInput(message.Text)
    if err != nil {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا: %v", err)))
        return
    }

    if _, exists := userFilters[message.Chat.ID]; !exists {
        userFilters[message.Chat.ID] = model.Filter{}
    }
    filter := userFilters[message.Chat.ID] // Retrieve the value (copy)
    filter.PriceMin = minPrice
    filter.PriceMax = maxPrice
    userFilters[message.Chat.ID] = filter // Store the modified value back into the map

    bot.Send(tgbotapi.NewMessage(message.Chat.ID, "رنج قیمت با موفقیت اعمال شد."))
    // Call the function to send the filter menu
    sendFilterMenu(bot, message.Chat.ID)
}



func handleCity(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً نام شهر مورد نظر خود را وارد کنید:")
	bot.Send(msg)
	userState[message.Chat.ID] = awaitingCity

}

func handleCitySearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	city := strings.TrimSpace(message.Text)
	if _, exists := userFilters[message.Chat.ID]; !exists {
		userFilters[message.Chat.ID] = model.Filter{}
	}
	filter := userFilters[message.Chat.ID]
	filter.City = city
	userFilters[message.Chat.ID] = filter
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("شهر '%s' با موفقیت اعمال شد.", city)))
	// Call the function to send the filter menu
	sendFilterMenu(bot, message.Chat.ID)

}

func handleNeighborhood(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً نام محله مورد نظر خود را وارد کنید:")
	bot.Send(msg)
	userState[message.Chat.ID] = "awaiting_neighborhood_input"
 
}

func handleNeighborhoodSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	neighborhood := strings.TrimSpace(message.Text)

	if _, exists := userFilters[message.Chat.ID]; !exists {
		userFilters[message.Chat.ID] = model.Filter{}
	}
	filter := userFilters[message.Chat.ID]
	filter.Neighborhood = neighborhood
	userFilters[message.Chat.ID] = filter

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("محله '%s' با موفقیت اعمال شد.", neighborhood)))
	// Call the function to send the filter menu
	sendFilterMenu(bot, message.Chat.ID)

}

func handleAreaRange(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً محدوده مساحت مورد نظر خود را به صورت (شروع,پایان) وارد نمایید:")
	bot.Send(msg)
	userState[message.Chat.ID] = "awaiting_area_range_input"
 
}

func handleAreaRangeSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
    minArea, maxArea, err := parseRangeInput(message.Text)
    if err != nil {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا: %v", err)))
        return
    }

    if _, exists := userFilters[message.Chat.ID]; !exists {
        userFilters[message.Chat.ID] = model.Filter{}
    }
    filter := userFilters[message.Chat.ID] // Retrieve the value (copy)
    filter.AreaMin = minArea
    filter.AreaMax = maxArea
    userFilters[message.Chat.ID] = filter // Store the modified value back into the map

    bot.Send(tgbotapi.NewMessage(message.Chat.ID, "محدوده مساحت با موفقیت اعمال شد."))
    // Call the function to send the filter menu
    sendFilterMenu(bot, message.Chat.ID)
}



func handleBedroomCount(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً تعداد اتاق خواب مورد نظر خود را به صورت (حداقل,حداکثر) وارد نمایید:")
	bot.Send(msg)
	userState[message.Chat.ID] = "awaiting_bedroom_count_input"
 
}

func handleBedroomCountSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {

    minRoomsFloat, maxRoomsFloat, err := parseRangeInput(message.Text)
    if err != nil {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا: %v", err)))
        return
    }

    minRooms := int(minRoomsFloat)
    maxRooms := int(maxRoomsFloat)

    if _, exists := userFilters[message.Chat.ID]; !exists {
        userFilters[message.Chat.ID] = model.Filter{}
    }
    filter := userFilters[message.Chat.ID] // Retrieve the value (copy)
    filter.RoomsMin = minRooms
    filter.RoomsMax = maxRooms
    userFilters[message.Chat.ID] = filter // Store the modified value back into the map

    bot.Send(tgbotapi.NewMessage(message.Chat.ID, "تعداد اتاق خواب با موفقیت اعمال شد."))
    // Call the function to send the filter menu
    sendFilterMenu(bot, message.Chat.ID)
}


func handleRentBuyMortgage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً نوع مورد نظر خود را وارد کنید: خرید/ فروش / اجاره / رهن / اجاره و رهن")
	bot.Send(msg)
	userState[message.Chat.ID] = "awaiting_rent_buy_mortgage_input"
 
}

func handleRentBuyMortgageSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	input := strings.TrimSpace(message.Text)

	if _, exists := userFilters[message.Chat.ID]; !exists {
		userFilters[message.Chat.ID] = model.Filter{}
	}
	filter := userFilters[message.Chat.ID]
	filter.Status = input
	userFilters[message.Chat.ID] = filter

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("نوع '%s' با موفقیت اعمال شد.", input)))
	// Call the function to send the filter menu
	sendFilterMenu(bot, message.Chat.ID)

}

func handleBuildingAge(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً محدوده سن بنا را به صورت (حداقل,حداکثر) وارد کنید:")
	bot.Send(msg)
	userState[message.Chat.ID] = "awaiting_building_age_input"
 
}

func handleBuildingAgeSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	input := strings.Split(message.Text, ",")
	if len(input) != 2 {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفاً محدوده سن بنا را به صورت (حداقل,حداکثر) وارد نمایید."))
		return
	}
	minAge, err1 := strconv.Atoi(strings.TrimSpace(input[0]))
	maxAge, err2 := strconv.Atoi(strings.TrimSpace(input[1]))
	if err1 != nil || err2 != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفا مقادیر را به درستی وارد نمایید."))
		return
	}

	currentYear := 1403
	if _, exists := userFilters[message.Chat.ID]; !exists {
		userFilters[message.Chat.ID] = model.Filter{}
	}
	filter := userFilters[message.Chat.ID]
	filter.BuildingAgeMin = currentYear - maxAge
	filter.BuildingAgeMax = currentYear - minAge
	userFilters[message.Chat.ID] = filter

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, "محدوده سن بنا با موفقیت اعمال شد."))
	// Call the function to send the filter menu
	sendFilterMenu(bot, message.Chat.ID)

}

func handleBuildingType(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً نوع ملک مورد نظر خود را وارد کنید: آپارتمان/ویلایی/غیره")
	bot.Send(msg)
	userState[message.Chat.ID] = "awaiting_building_type_input"
 
}

func handleBuildingTypeSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	propertyType := strings.TrimSpace(message.Text)

	if _, exists := userFilters[message.Chat.ID]; !exists {
		userFilters[message.Chat.ID] = model.Filter{}
	}
	filter := userFilters[message.Chat.ID]
	filter.PropertyType = propertyType
	userFilters[message.Chat.ID] = filter

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("نوع ملک '%s' با موفقیت اعمال شد.", propertyType)))
	// Call the function to send the filter menu
	sendFilterMenu(bot, message.Chat.ID)

}

func handleFloorRange(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً محدوده طبقه را به صورت (حداقل,حداکثر) وارد کنید:")
	bot.Send(msg)
	userState[message.Chat.ID] = "awaiting_floor_range_input"
 
}

func handleFloorRangeSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
    minFloorFloat, maxFloorFloat, err := parseRangeInput(message.Text)
    if err != nil {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا: %v", err)))
        return
    }

    minFloor := int(minFloorFloat)
    maxFloor := int(maxFloorFloat)

    if _, exists := userFilters[message.Chat.ID]; !exists {
        userFilters[message.Chat.ID] = model.Filter{}
    }
    filter := userFilters[message.Chat.ID] // Retrieve the value (copy)
    filter.FloorMin = &minFloor
    filter.FloorMax = &maxFloor
    userFilters[message.Chat.ID] = filter // Store the modified value back into the map

    bot.Send(tgbotapi.NewMessage(message.Chat.ID, "محدوده طبقه با موفقیت اعمال شد."))
    // Call the function to send the filter menu
    sendFilterMenu(bot, message.Chat.ID)
}


func handleStorage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "آیا ملک دارای انباری باشد؟ (بله/خیر)")
	bot.Send(msg)
	userState[message.Chat.ID] = "awaiting_storage_input"
 
}

func handleStorageSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	input := strings.TrimSpace(message.Text)
	var hasStorage bool
	if input == "بله" {
		hasStorage = true
	} else if input == "خیر" {
		hasStorage = false
	} else {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفاً 'بله' یا 'خیر' وارد کنید."))
		return
	}

	if _, exists := userFilters[message.Chat.ID]; !exists {
		userFilters[message.Chat.ID] = model.Filter{}
	}
	filter := userFilters[message.Chat.ID]
	filter.HasStorage = hasStorage
	userFilters[message.Chat.ID] = filter

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, "وضعیت انباری با موفقیت اعمال شد."))
	// Call the function to send the filter menu
	sendFilterMenu(bot, message.Chat.ID)

}

func handleElevator(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "آیا ملک دارای آسانسور باشد؟ (بله/خیر)")
	bot.Send(msg)
	userState[message.Chat.ID] = "awaiting_elevator_input"
 
}

func handleElevatorSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	input := strings.TrimSpace(message.Text)
	var hasElevator bool
	if input == "بله" {
		hasElevator = true
	} else if input == "خیر" {
		hasElevator = false
	} else {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفاً 'بله' یا 'خیر' وارد کنید."))
		return
	}

	if _, exists := userFilters[message.Chat.ID]; !exists {
		userFilters[message.Chat.ID] = model.Filter{}
	}
	filter := userFilters[message.Chat.ID]
	filter.HasElevator = hasElevator
	userFilters[message.Chat.ID] = filter

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, "وضعیت آسانسور با موفقیت اعمال شد."))
	// Call the function to send the filter menu
	sendFilterMenu(bot, message.Chat.ID)

}

func handleAdCreationDate(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "لطفاً محدوده تاریخ درج آگهی را به صورت (شروع,پایان) وارد کنید (YYYY-MM-DD):")
	bot.Send(msg)
	userState[message.Chat.ID] = "awaiting_ad_creation_date_input"
 
}

func handleAdCreationDateSearch(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *gorm.DB) {
	input := strings.Split(message.Text, ",")
	if len(input) != 2 {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفاً محدوده تاریخ را به صورت (شروع,پایان) به فرمت YYYY-MM-DD وارد نمایید."))
		return
	}
	startDate, err1 := time.Parse("2006-01-02", strings.TrimSpace(input[0]))
	endDate, err2 := time.Parse("2006-01-02", strings.TrimSpace(input[1]))
	if err1 != nil || err2 != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "ورودی نامعتبر است. لطفا تاریخ‌ها را به درستی وارد نمایید."))
		return
	}

	if _, exists := userFilters[message.Chat.ID]; !exists {
		userFilters[message.Chat.ID] = model.Filter{}
	}
	filter := userFilters[message.Chat.ID]
	filter.CreationDateMin = startDate
	filter.CreationDateMax = endDate
	userFilters[message.Chat.ID] = filter

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, "محدوده تاریخ درج آگهی با موفقیت اعمال شد."))
	
	// Call the function to send the filter menu
	sendFilterMenu(bot, message.Chat.ID)

}


func sendFilterMenu(bot *tgbotapi.BotAPI, chatID int64) {
    menuMsg := tgbotapi.NewMessage(chatID, "منوی فیلترها باز است. می‌توانید فیلترهای دیگری انتخاب کنید یا \"تایید فیلترها\" را بزنید.")
    menuMsg.ReplyMarkup = filterKeyboard
    bot.Send(menuMsg)
}

func parseRangeInput(input string) (float64, float64, error) {
    parts := strings.Split(input, ",")
    if len(parts) == 1 {

        max, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
        if err != nil {
            return 0, 0, fmt.Errorf("مقدار وارد شده معتبر نیست: %v", err)
        }
        return 0, max, nil
    } else if len(parts) == 2 {

        min, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
        max, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
        if err1 != nil || err2 != nil {
            return 0, 0, fmt.Errorf("مقادیر وارد شده معتبر نیستند")
        }
        return min, max, nil
    } else {

        return 0, 0, fmt.Errorf("ورودی نامعتبر است. لطفاً فقط یک یا دو عدد وارد کنید.")
    }
}

// Generates a CSV file with the given data
func generateCSVFile(fileName string, data []model.Listing) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	headers := []string{"ID", "Title", "Price", "City", "Bedrooms", "Area", "Ad Type", "Creation Date"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error writing header to CSV: %v", err)
	}

	// Write data rows
	for _, listing := range data {
		record := []string{
			fmt.Sprintf("%d", listing.ListingID),
			listing.Title,
			fmt.Sprintf("%f", listing.Price),
			listing.City,
			fmt.Sprintf("%d", listing.Bedrooms),
			fmt.Sprintf("%f", listing.Meterage),
			listing.AdType,
			listing.CreatedAt.String(),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record to CSV: %v", err)
		}
	}

	return nil
}


func createZipFile(zipFileName, fileName string) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileToZip, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = fileName
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}

func handleDownloadCSV(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    // Retrieve filtered results for the user
    filteredResults, exists := userFilteredResults[message.Chat.ID]
    if !exists || len(filteredResults) == 0 {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, "هیچ نتیجه‌ای برای دانلود یافت نشد. لطفاً ابتدا جستجو کنید."))
        return
    }

    csvFileName := "results.csv"
    zipFileName := "results.zip"

    // Generate CSV with the retrieved filtered results
    err := generateCSVFile(csvFileName, filteredResults)
    if err != nil {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در ایجاد فایل CSV: %v", err)))
        return
    }

    // Create ZIP file
    err = createZipFile(zipFileName, csvFileName)
    if err != nil {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در ایجاد فایل ZIP: %v", err)))
        return
    }

    // Send the ZIP file
    file := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FilePath(zipFileName))
    _, err = bot.Send(file)
    if err != nil {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("خطا در ارسال فایل: %v", err)))
        return
    }
}
package bot

import (
	"CrawlerProject/internal/model"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendFormattedListings(bot *tgbotapi.BotAPI, chatID int64, results []model.Listing) {
	for _, result := range results {
		msgText := fmt.Sprintf(
			"عنوان: %s\n"+
				"قیمت: %f\n"+
				"شهر: %s\n"+
				"محله: %s\n"+
				"متراژ: %d متر مربع\n"+
				"تعداد اتاق خواب: %d\n"+
				"نوع آگهی: %s\n"+
				"سن بنا: %s سال\n"+
				"نوع ملک: %s\n"+
				"طبقه: %d\n"+
				"انباری: %s\n"+
				"آسانسور: %s\n"+
				"تاریخ درج آگهی: %s\n"+
				"تاریخ ایجاد: %s\n"+
				"تاریخ بروزرسانی: %s\n"+
				"تصاویر: %s\n"+
				"[لینک](%s)",
			result.Title,
			result.Price,
			result.City,
			result.Neighborhood,
			result.Meterage,
			result.Bedrooms,
			result.AdType,
			result.Age,
			result.HouseType,
			result.Floor,
			func() string {
				if result.Warehouse {
					return "دارد"
				}
				return "ندارد"
			}(),
			func() string {
				if result.Elevator {
					return "دارد"
				}
				return "ندارد"
			}(),
			result.AdCreateDate,
			result.CreatedAt,
			result.UpdatedAt,
			result.Images,
			result.Link)

		// Create an inline keyboard for bookmarking and downloading as ZIP
		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("بوکمارک کردن", fmt.Sprintf("bookmark_%d", result.ListingID)),
			),
		)

		msg := tgbotapi.NewMessage(chatID, msgText)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = markup
		if _, err := bot.Send(msg); err == nil {
			fmt.Println("Result sent to Telegram successfully")
		} else {
			fmt.Printf("Failed to send result to Telegram: %v\n", err)
		}
	}
}

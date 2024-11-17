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
				"[لینک](%s)\n"+
				"دانلود zip"+
				"ارسال به ایمیل",
			result.Title,
			result.Price,
			result.City,
			result.Neighborhood,
			result.Area,
			result.Rooms,
			result.Status,
			result.Age,
			result.HouseType,
			result.Floor,
			func() string {
				if result.HasStorage {
					return "دارد"
				}
				return "ندارد"
			}(),
			func() string {
				if result.HasElevator {
					return "دارد"
				}
				return "ندارد"
			}(),
			result.CreatedAt,
			result.CreatedAt,
			result.UpdatedAt,
			result.Images,
			result.URL)

		msg := tgbotapi.NewMessage(chatID, msgText)
		if _, err := bot.Send(msg); err == nil {
			fmt.Println("Result sent to Telegram successfully")
		} else {
			fmt.Printf("Failed to send result to Telegram: %v\n", err)
		}
	}
}

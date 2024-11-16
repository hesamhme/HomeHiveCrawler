package bot


import (
    "fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"CrawlerProject/internal/model"

)

// sendFormattedListings formats and sends a list of listings to a given chat ID
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
                "[لینک](%s)\n",
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

        msg := tgbotapi.NewMessage(chatID, msgText)
        bot.Send(msg)
    }
}

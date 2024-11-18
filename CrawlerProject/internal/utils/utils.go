package utils

import (
	"CrawlerProject/internal/model"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var persianDigitMap = map[rune]rune{
	'۰': '0', '۱': '1', '۲': '2', '۳': '3', '۴': '4',
	'۵': '5', '۶': '6', '۷': '7', '۸': '8', '۹': '9',
}

const persianToEnglishJS = `
function persianToEnglish(str) {
	var numbers = {'۰':'0','۱':'1','۲':'2','۳':'3','۴':'4','۵':'5','۶':'6','۷':'7','۸':'8','۹':'9'};
	return str.replace(/[۰-۹]/g, function(d) { return numbers[d]; });
}
`

var persianToGregorianMonth = map[string]int{
	"فروردین":  1,  // March
	"اردیبهشت": 2,  // April
	"خرداد":    3,  // May
	"تیر":      4,  // June
	"مرداد":    5,  // July
	"شهریور":   6,  // August
	"مهر":      7,  // September
	"آبان":     8,  // October
	"آذر":      9,  // November
	"دی":       10, // December
	"بهمن":     11, // January
	"اسفند":    12, // February
}

// persianToLatinDigits converts Persian digits to Latin digits
var persianToLatinDigits = map[rune]rune{
	'۰': '0',
	'۱': '1',
	'۲': '2',
	'۳': '3',
	'۴': '4',
	'۵': '5',
	'۶': '6',
	'۷': '7',
	'۸': '8',
	'۹': '9',
}

func evaluateScript(selector string, conversion string) string {
	return fmt.Sprintf(`
		(function() {
			%s
			var el = %s;
			if (!el) return '';
			return %s;
		})()
	`, persianToEnglishJS, selector, conversion)
}

func EvaluateNumericScript(selector string) string {
	return fmt.Sprintf(`
		(function() {
			%s
			var el = %s;
			if (!el) return 0;
			var text = el.innerText.trim();
			return parseInt(persianToEnglish(text)) || 0;
		})()
	`, persianToEnglishJS, selector)
}

func UniqueAds(ads []model.Listing) []model.Listing {
	// return ads
	seen := make(map[string]bool)
	unique := []model.Listing{}

	for _, ad := range ads {
		key := ad.Title + "|" + ad.URL
		if !seen[key] && strings.TrimSpace(ad.Title) != "" {
			seen[key] = true
			unique = append(unique, ad)
		}
	}
	return unique
}

func convertPersianToLatinDigits(str string) string {
	result := make([]rune, 0, len(str))
	for _, r := range str {
		if latin, ok := persianToLatinDigits[r]; ok {
			result = append(result, latin)
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func ExtractPersianDate(text string) (time.Time, error) {
	// Regular expression to match Persian date pattern
	re := regexp.MustCompile(`(\d{1,2})\s+([\p{L}]+)\s+(\d{4})`)

	// Convert Persian digits to Latin digits
	latinText := convertPersianToLatinDigits(text)

	// Find the date pattern
	matches := re.FindStringSubmatch(latinText)
	if len(matches) != 4 {
		return time.Time{}, fmt.Errorf("date pattern not found")
	}

	// Extract components
	day := matches[1]
	month := matches[2]
	year := matches[3]

	// Convert Persian year to Gregorian year
	persianYear, err := strconv.Atoi(year)
	if err != nil {
		return time.Time{}, err
	}
	gregorianYear := persianYear // Approximate conversion

	// Get Gregorian month number
	gregorianMonth, ok := persianToGregorianMonth[month]
	if !ok {
		return time.Time{}, fmt.Errorf("invalid Persian month: %s", month)
	}

	// Convert day to integer
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		return time.Time{}, err
	}

	// Create time.Time object
	// Note: This is an approximation as exact Persian to Gregorian conversion
	// requires more complex calculations
	return time.Date(gregorianYear, time.Month(gregorianMonth), dayInt,
		0, 0, 0, 0, time.UTC), nil
}

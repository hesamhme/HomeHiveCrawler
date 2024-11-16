package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/exp/rand"
	model "CrawlerProject/internal/model"
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

func NewCrawler(config model.CrawlerConfig) *model.Crawler {
	return &Crawler{
		config:           config,
		urlSemaphore:     make(chan struct{}, config.MaxURLConcurrency),
		adsSemaphore:     make(chan struct{}, config.MaxAdConcurrency),
		errorChan:        make(chan error, len(config.Cities)*len(config.Types)),
		resultsChan:      make(chan HouseAd, 10000),
		goroutineMonitor: NewGoroutineMonitor(),
	}
}

// DefaultConfig returns the default configuration
func DefaultConfig() CrawlerConfig {
	return CrawlerConfig{
		RunInterval:        5 * time.Hour,
		MinTimeBetweenRuns: time.Duration(float64(5*time.Hour) * 0.9),
		PageTimeout:        30 * time.Minute,
		AdTimeout:          20 * time.Minute,
		MaxURLConcurrency:  2,
		MaxAdConcurrency:   5,
		// Cities:             []string{"tabriz", "yazd", "tehran", "shiraz"},
		Cities: []string{"yazd"},
		// Cities:             []string{"tabriz", "azarshahr", "ahar", "bonab", "sarab", "sahand", "maragheh", "marand", "mianeh", "urmia", "oshnavieh", "bukan", "piranshahr", "khoy", "sardasht", "salmas", "shahin-dej", "maku", "mahabad", "miandoab", "naqadeh", "ardabil", "parsabad", "khalkhal", "sarein", "germi", "meshgin-shahr", "namin", "isfahan", "aran-va-bidgol", "abrisham-isfahan", "khomeyni-shahr", "khansar", "khour", "daran", "semirom", "shahin-shahr", "falavarjan", "foolad-shahr", "ghamsar", "kashan", "golpayegan", "lenjan", "mobarakeh", "najafabad", "karaj", "asara", "eshtehard", "tankaman", "charbagh-alborz", "taleqan", "fardis", "koohsar", "garmdareh", "mahdasht", "mohammad-shahr", "nazarabad", "hashtgerd", "abdanan", "ilam", "eyvan", "dehloran", "mehran", "borazjan", "dayyer", "bandar-kangan", "bandar-ganaveh", "bushehr", "jam", "khormoj", "tehran", "absard", "abali", "arjmand", "eslamshahr", "andisheh-new-town", "baghershahr", "bumehen", "pakdasht", "pardis", "parand", "pishva", "javadabad", "chahar-dangeh", "damavand", "robat-karim", "rudehen", "shahr-e-rey", "shahedshahr", "shemshak", "shahriar", "sabashahr", "safadasht-industrial-city", "ferdosiye", "fasham", "firuzkooh", "qods", "qarchak", "kahrizak", "kilan", "golestan-baharestan", "lavasan", "nasimshahr", "vahidieh", "varamin", "boroujen", "saman", "shahrekord", "farrokhshahr", "lordegan", "birjand", "tabas", "ferdows", "ghayen", "mashhad", "bardaskan", "taybad", "torbat-jam", "torbat-heydariyeh", "chenaran", "khaf", "sabzevar", "shandiz", "torghabeh", "qasemabad-khaf", "quchan", "golbahar", "gonabad", "molkabad", "neyshabur", "ashkhaneh", "esfarāyen", "bojnurd", "shirvan", "ahvaz", "abadan", "omidiyeh", "andimeshk", "izeh", "bandar-imam-khomeini", "bandar-mahshahr", "behbahan", "chamran-town", "hamidiyeh", "khorramshahr", "dezful", "ramshir", "ramhormoz", "susangerd", "shadeghan", "shush", "shooshtar", "masjed-soleyman", "hendijan", "abhar", "khorramdarreh", "zanjan", "qeydar", "damghan", "semnan", "shahroud", "garmsar", "iranshahr", "chabahar", "khash", "zabol", "zahedan", "zahak", "saravan", "konarak", "shiraz", "abadeh", "eqlid", "jahrom", "khoour", "darab", "zarghan", "sadra", "fasa", "firuzabad", "kazeroon", "lar", "lamerd", "marvdasht", "mohr", "norabad", "neyriz", "abyek", "eqbaliyeh", "alvand", "takestan", "shal", "qazvin", "mohammadiyeh", "qom", "baneh", "bijar", "dehgolan", "saqqez", "sanandaj", "qorveh", "kamyaran", "marivan", "baft", "bardsir", "boluk", "bam", "jiroft", "rafsanjan", "zarand", "sirjan", "kerman", "kahnooj", "mahan", "kermanshah", "eslamabad-gharb", "bisotun", "javanrud", "sarpol-zahab", "sonqor", "sahneh", "kangavar", "gahvareh", "harsin", "dogonbadan", "dehdasht", "sisakht", "yasuj", "azadshahr-golestan", "aq-qala", "bandar-torkaman", "aliabad-katul", "kordkuy", "kalale", "galikesh", "gorgan", "gomishan", "gonbad-kavus", "minoodasht", "sangdovin", "sorkhan-kalateh", "faragi", "sadegh-abad", "bandar-gaz", "maraveh-tapeh", "daland", "negin-shahr", "ramiyan", "khan-bin", "jelin", "dozin", "nokandeh", "goli-dagh", "nodeh-khandoz", "anbaralum", "fazel-abad", "mazrae-katool", "yanghagh", "sijval", "simin-shahr", "tatar-olya", "alghajar", "ghorogh", "inche-borun", "rasht", "astara", "astaneh-ashrafiyeh", "ahmadsar-gourab", "asalem", "amlash", "barah-sar", "bandar-anzali", "pareh-sar", "talesh", "toutkabon", "jirandeh", "chaboksar", "chaf-chamkhale", "chobar", "haviq", "khoshkbijar", "khomam", "deylaman", "rankouh", "rahim-abad", "rostam-abad", "rezvanshahr", "rudbar", "roudbaneh", "rudsar", "zibakenar", "sangar", "siahkal", "shaft", "shelman", "someh-sara", "fuman", "kelachay", "kouchesfahan", "koumeleh", "kiashahr", "gourab-zarmikh", "lahijan", "lashtenesha", "langarud", "loshan", "loulman", "lavandevil", "lisar", "masal", "masuleh", "makloan", "manjil", "vajargah", "tahergurab", "shanderman", "ziyabar", "otaghvar", "tulam-shahr", "pirbazar", "azna", "aleshtar", "aligudarz", "borujerd", "pol-dokhtar", "khorramabad", "dorud", "kuhdasht", "nurabad", "aalasht", "amol", "amirkala", "izadshahr", "babol", "babolsar", "baladeh", "behshahr", "bahnamir", "polsefid", "tonekabon", "juybar", "chalus", "chamestan", "khalil-shahr", "khoshroud-pey", "ramsar", "rostamkola", "royan", "reyneh", "ziraab", "sari", "sorkhrood", "salman-shahr", "sourek", "shirgah", "abbasabad-mazandaran", "farahabad", "fereydunkenar", "farim", "qaemshahr", "katalem-sadatshahr", "kelarabad", "kelarestan", "kouhi-kheyl", "kiasar", "kiakola", "gatab", "gazanak", "galougah-babol", "mahmudabad", "marzan-abad", "marzikola", "nashtarud", "neka", "nur", "nowshahr", "paeen-holar", "dalkhani", "galugah-babol", "hadi-shahr", "babakan", "zargarshahr", "arateh", "emamzadeh-abdollah", "shirud", "dabudasht", "akand", "astaneh-sara", "pool", "tabaghdeh", "kojur", "khoram-abad", "hachirud", "arak", "khomein", "delijan", "saveh", "shazand", "mahalat", "mohajeran", "bandar-abbas", "takht", "dargahan", "qeshm", "kish", "minab", "hormuz", "asadabad", "bahar", "tuyserkan", "kabudrahang", "malayer", "nahavand", "hamedan", "ardakan", "bafq", "taft", "hamidia", "mehriz", "meybod", "yazd"},
		// Types: []string{"buy-apartment"},
		// Types: []string{"buy-apartment", "buy-villa", "rent-apartment", "rent-villa"},
		Types:     []string{"buy-villa"},
		OutputDir: "crawler_output",
		ChromeFlags: append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
		),
	}
}

// Start begins the crawler's operation
func (c *Crawler) Start(ctx context.Context) error {
	log.Printf("Starting crawler with interval: %v", c.config.RunInterval)

	// Run immediately on startup
	if err := c.RunOnce(ctx); err != nil {
		log.Printf("Initial run failed: %v", err)
	}

	// Setup ticker for periodic runs
	ticker := time.NewTicker(c.config.RunInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := c.checkAndRun(ctx); err != nil {
				log.Printf("Scheduled run failed: %v", err)
			}
		}
	}
}

// checkAndRun verifies if enough time has passed and runs the crawler
func (c *Crawler) checkAndRun(ctx context.Context) error {
	c.runMutex.Lock()
	if time.Since(c.lastRunTime) < c.config.MinTimeBetweenRuns {
		c.runMutex.Unlock()
		return nil
	}
	c.runMutex.Unlock()
	return c.RunOnce(ctx)
}

// RunOnce performs a single crawl operation
func (c *Crawler) RunOnce(ctx context.Context) error {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	log.Printf("Starting crawl at %v", time.Now())
	defer func() {
		c.lastRunTime = time.Now()
		log.Printf("Completed crawl at %v", time.Now())
	}()

	return c.crawl(ctx)
}

// crawl performs the actual crawling operation
func (c *Crawler) crawl(ctx context.Context) error {
	// Setup browser context
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, c.config.ChromeFlags...)
	defer allocCancel()

	// Create timeout context for entire operation
	crawlCtx, cancel := context.WithTimeout(allocCtx, c.config.PageTimeout)
	defer cancel()

	var wg sync.WaitGroup
	var allAds []HouseAd

	// Process each URL concurrently
	for _, city := range c.config.Cities {
		for _, _type := range c.config.Types {
			wg.Add(1)
			go func(city, _type string) {
				defer wg.Done()

				// Start monitoring this goroutine
				stats := c.goroutineMonitor.StartTracking(city, _type)
				defer c.goroutineMonitor.StopTracking(stats.GoroutineID)

				if err := c.processURL(crawlCtx, city, _type, stats, &allAds); err != nil {
					select {
					case c.errorChan <- err:
					default:
					}
				}
			}(city, _type)
		}
	}

	wg.Wait()

	// Save goroutine statistics
	if err := c.goroutineMonitor.SaveStats(c.config.OutputDir); err != nil {
		log.Printf("Error saving goroutine stats: %v", err)
	}

	// Process gathered ads
	return c.processAds(crawlCtx, &allAds)
}

// processURL handles crawling a single URL
func (c *Crawler) processURL(ctx context.Context, city, _type string, stats *GoroutineStats, allAds *[]HouseAd) error {
	// Acquire URL semaphore
	c.urlSemaphore <- struct{}{}
	defer func() { <-c.urlSemaphore }()

	url := "https://divar.ir/s/" + city + "/" + _type
	stats.URL = url

	// Create new browser context
	browserCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var urlAds []HouseAd
	var adsWg sync.WaitGroup
	adsWg.Add(1)

	log.Printf("Processing URL: %s", url)

	// Run chromedp for this URL
	err := chromedp.Run(browserCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second),
		c.scrollAndScrape(&urlAds, &adsWg),
	)

	if err != nil {
		return fmt.Errorf("error processing URL %s: %w", url, err)
	}

	adsWg.Wait()

	// Update statistics
	stats.NumAdsFound = len(urlAds)

	// Safely append the ads
	c.adsMutex.Lock()
	*allAds = append(*allAds, urlAds...)
	c.adsMutex.Unlock()

	log.Printf("Completed URL %s: Found %d ads", url, len(urlAds))
	return nil
}

// scrollAndScrape implements the scrolling and scraping logic
// Original function with issues identified and fixed
func (c *Crawler) scrollAndScrape(ads *[]HouseAd, wg *sync.WaitGroup) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		// Only one Done() call is needed at the end
		defer wg.Done()

		var previousHeight int64
		adChannel := make(chan []HouseAd)
		done := make(chan struct{})
		var mu sync.Mutex
		var once sync.Once

		// Goroutine for scraping and scrolling
		go func() {
			defer once.Do(func() { close(adChannel) })
			defer close(done)

			i := 0
			for i < 3 && ctx.Err() == nil {

				var currentHeight int64
				if err := chromedp.Evaluate(`document.documentElement.scrollHeight`, &currentHeight).Do(ctx); err != nil {
					log.Println("Error getting current height:", err)
					return
				}

				mu.Lock()
				prevHeight := previousHeight
				mu.Unlock()

				if currentHeight == prevHeight {
					var hasMoreButton bool
					log.Println("Checking for 'show more' button...")
					err := chromedp.Evaluate(`
                        (() => {
                            const button = document.querySelector('.post-list__load-more-btn-be092');
                            if (button && button.innerText.includes('آگهی‌های بیشتر')) {
                                button.click();
                                return true;
                            }
                            return false;
                        })();
                    `, &hasMoreButton).Do(ctx)

					if err != nil {
						log.Println("Error evaluating 'show more' button:", err)
						return
					}

					if !hasMoreButton {
						log.Println("No 'show more' button found, scraping finished")
						return
					}
					i++
				}

				var newAds []HouseAd
				err := chromedp.Evaluate(`
                    Array.from(document.querySelectorAll('.kt-post-card')).map(card => ({
                        title: card.querySelector('.kt-post-card__title')?.innerText || '',
                        link: card.querySelector('a')?.href || ''
                    }))
                `, &newAds).Do(ctx)

				if err != nil {
					log.Println("Error extracting ads:", err)
					return
				}

				log.Printf("Found %d new ads", len(newAds))

				select {
				case adChannel <- newAds:
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					log.Println("Timeout sending ads to channel")
					return
				}

				if err := chromedp.Evaluate(`window.scrollTo(0, document.documentElement.scrollHeight)`, nil).Do(ctx); err != nil {
					log.Println("Error scrolling page:", err)
					return
				}

				select {
				case <-time.After(1 * time.Second):
				case <-ctx.Done():
					return
				}

				mu.Lock()
				previousHeight = currentHeight
				mu.Unlock()
			}
		}()

		// Collect ads from the channel
		go func() {
			for {
				select {
				case newAds, ok := <-adChannel:
					if !ok {
						return
					}
					mu.Lock()
					*ads = uniqueAds(append(*ads, newAds...))
					log.Printf("Current number of unique ads found: %d", len(*ads))
					mu.Unlock()
				case <-ctx.Done():
					return
				}
			}
		}()

		// Wait until scraping and collecting is done
		select {
		case <-done:
			log.Println("Scraping finished.")
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	}
}

// processAds handles the processing of gathered ads
func (c *Crawler) processAds(ctx context.Context, ads *[]HouseAd) error {
	totalAds := len(*ads)
	if totalAds == 0 {
		return fmt.Errorf("no ads found during scraping")
	}

	log.Printf("Processing details for %d ads", totalAds)

	var wg sync.WaitGroup
	// for i := range *ads {
	for i := range 6 {
		wg.Add(1)
		go func(ad *HouseAd, index int) {
			defer wg.Done()

			// Reactivate semaphore control
			c.adsSemaphore <- struct{}{}
			defer func() { <-c.adsSemaphore }()

			// Add retry logic
			maxRetries := 3
			var err error
			for retry := 0; retry < maxRetries; retry++ {
				if err = c.processAdDetails(ctx, ad, index); err == nil {
					break
				}
				if retry < maxRetries-1 {
					time.Sleep(time.Duration(1<<uint(retry)) * time.Second)
				}
			}

			if err != nil {
				select {
				case c.errorChan <- fmt.Errorf("failed after %d retries: %w", maxRetries, err):
				default:
				}
				return
			}

			select {
			case c.resultsChan <- *ad:
			case <-ctx.Done():
			}
		}(&(*ads)[i], i)
	}

	go func() {
		wg.Wait()
		close(c.errorChan)
		close(c.resultsChan)
	}()

	var processedAds []HouseAd
	var errors []error

	for c.errorChan != nil || c.resultsChan != nil {
		select {
		case err, ok := <-c.errorChan:
			if !ok {
				c.errorChan = nil
				continue
			}
			errors = append(errors, err)
		case ad, ok := <-c.resultsChan:
			if !ok {
				c.resultsChan = nil
				continue
			}
			processedAds = append(processedAds, ad)
		}
	}

	if len(errors) > 0 {
		log.Printf("Encountered %d errors during processing:", len(errors))
		for _, err := range errors {
			log.Printf("- %v", err)
		}
	}

	c.SaveResults(&processedAds)
	return nil
}

// processAdDetails handles fetching details for a single ad
func (c *Crawler) processAdDetails(ctx context.Context, ad *HouseAd, index int) error {
	// Create new browser context for each ad
	fmt.Println("crawling ", ad.Link)
	browserCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Add timeout
	timeoutCtx, timeoutCancel := context.WithTimeout(browserCtx, c.config.AdTimeout)
	defer timeoutCancel()

	// Navigate to ad page first
	if err := chromedp.Run(timeoutCtx, chromedp.Navigate(ad.Link)); err != nil {
		return fmt.Errorf("failed to navigate to ad page: %w", err)
	}

	// Add random delay
	delay := time.Duration(1000+rand.Intn(1000)) * time.Millisecond
	select {
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	case <-time.After(delay):
	}
	// Create a slice of tasks to execute
	tasks := []struct {
		description string
		action      func(context.Context) error
	}{
		{
			description: "Get meterage",
			action: func(adCtx context.Context) error {
				// Try primary selector first
				err := chromedp.Evaluate(evaluateNumericScript(
					"document.querySelectorAll('.kt-group-row__data-row .kt-group-row-item__value')[0]",
				), &ad.Meterage).Do(adCtx)
				// Fall back to secondary selector if primary fails
				if err != nil || ad.Meterage == 0 {
					return chromedp.Evaluate(evaluateNumericScript(
						"document.querySelector('.kt-unexpandable-row__value')",
					), &ad.Meterage).Do(adCtx)
				}
				return err
			},
		},
		{
			description: "Get bedrooms",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(evaluateNumericScript(
					"document.querySelector('.kt-group-row__data-row td:nth-child(3)')",
				), &ad.Bedrooms).Do(adCtx)
			},
		},
		{
			description: "Get city",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(`
						(function() {
							var el = document.querySelector('.kt-page-title__subtitle');
							if (!el) return '';
							var text = el.innerText || '';
							var parts = text.split('در');
							return parts.length > 1 ? parts[1].trim() : '';
						})()
					`, &ad.City).Do(adCtx)
			},
		},
		{
			description: "Get description",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(`(() => {
						const el = document.querySelector('.kt-description-row__text.kt-description-row__text--primary');
						return el ? el.innerText.trim() : '';
					})()`,
					&ad.Description).Do(adCtx)
			},
		},
		{
			description: "Get seller contact",
			action: func(adCtx context.Context) error {
				a := chromedp.Run(adCtx,
					chromedp.WaitVisible(`.post-actions__get-contact`),

					// Click the button
					chromedp.Click(`.post-actions__get-contact`),

					// Wait a bit for the phone number to appear
					chromedp.Sleep(1*time.Second),

					// Run JavaScript to extract and convert the phone number
					chromedp.Evaluate(`
								(() => {
									function persianToEnglish(str) {
										const numbers = {'۰':'0','۱':'1','۲':'2','۳':'3','۴':'4','۵':'5','۶':'6','۷':'7','۸':'8','۹':'9'};
										return str.replace(/[۰-۹]/g, d => numbers[d]);
									}

									const phoneElement = document.querySelector('.copy-row a.kt-unexpandable-row__action');
									if (!phoneElement) {
										return '';
									}

									return persianToEnglish(phoneElement.textContent.trim());
								})()
								`, &ad.Seller),
				)
				fmt.Println(">>>>>>", ad.Seller)
				return a

			},
		},
		{
			description: "Get house type",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(`
						(function() {
							// Find the navigation button in the post page section
							var buttonSpan = document.querySelector('.post-page__section--padded .kt-chip span');
							if (!buttonSpan) {
								return '';
							}

							// Get the full text and split by spaces
							var fullText = buttonSpan.innerText || buttonSpan.textContent || '';
							var words = fullText.split(' ');

							// Return everything after the first word
							if (words.length <= 1) {
								return '';
							}

							// Join all words after the first one
							return words.slice(1).join(' ');
						})()
					`, &ad.HouseType).Do(adCtx)
			},
		},
		{
			description: "Get ad type",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(`
												(function() {
							// Find the navigation button in the post page section
							var buttonSpan = document.querySelector('.post-page__section--padded .kt-chip span');
							if (!buttonSpan) {
								return '';
							}
							// Get the full text and split by spaces
							var fullText = buttonSpan.innerText || buttonSpan.textContent || '';
							var words = fullText.split(' ');
							// Return the first word
							if (words.length < 1) {
								return '';
							}
							// Return first word
							return words[0];
						})()
					`, &ad.AdType).Do(adCtx)
			},
		},
		{
			description: "Check amenities",
			action: func(adCtx context.Context) error {
				amenitiesScript := `
						(function() {
							var amenities = Array.from(document.querySelectorAll('.kt-group-row__data-row .kt-body.kt-body--stable'))
								.map(function(el) { return el.textContent; });
							return {
								hasElevator: amenities.indexOf('آسانسور') !== -1,
								hasWarehouse: amenities.indexOf('انباری') !== -1
							};
						})()
					`
				var amenities struct {
					HasElevator  bool `json:"hasElevator"`
					HasWarehouse bool `json:"hasWarehouse"`
				}
				if err := chromedp.Evaluate(amenitiesScript, &amenities).Do(adCtx); err != nil {
					return err
				}
				ad.Elevator = amenities.HasElevator
				ad.WareHouse = amenities.HasWarehouse
				return nil
			},
		},
		{
			description: "Get floor number",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(`
						(function() {
							var numbers = {'۰':'0','۱':'1','۲':'2','۳':'3','۴':'4','۵':'5','۶':'6','۷':'7','۸':'8','۹':'9'};
							var floors = Array.from(document.querySelectorAll('.kt-unexpandable-row__title-box p'));
							var floorEl = null;
							for (var i = 0; i < floors.length; i++) {
								if (floors[i].innerText.includes("طبقه")) {
									floorEl = floors[i];
									break;
								}
							}
							if (!floorEl) return 0;
							var parent = floorEl.closest('.kt-unexpandable-row__title-box');
							if (!parent) return 0;
							var next = parent.nextElementSibling;
							if (!next) return 0;
							var value = next.querySelector('.kt-unexpandable-row__value');
							if (!value) return 0;
							var text = value.innerText || '0';
							var english = text.replace(/[۰-۹]/g, function(d) { return numbers[d]; });
							return parseInt(english, 10) || 0;
						})()
					`, &ad.Floor).Do(adCtx)
			},
		},
		{
			description: "Get age",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(`
						(() => {
							const persianToLatin = {
								'۰': '0', '۱': '1', '۲': '2', '۳': '3', '۴': '4',
								'۵': '5', '۶': '6', '۷': '7', '۸': '8', '۹': '9'
							};

							const headerCells = document.querySelectorAll('table.kt-group-row th');
							const yearColumnIndex = Array.from(headerCells).findIndex(
								th => th.textContent.includes('ساخت')
							);
							if (yearColumnIndex === -1) return null;

							const dataRow = document.querySelector('table.kt-group-row tbody tr');
							if (!dataRow) return null;

							const yearCell = dataRow.querySelectorAll('td')[yearColumnIndex];
							if (!yearCell) return null;

							// Get the Persian numeral and convert to Latin
							const persianYear = yearCell.textContent.trim();
							const latinYear = persianYear.split('').map(char => persianToLatin[char] || char).join('');

							return latinYear;
						})()
					`, &ad.Age).Do(adCtx)
			},
		},
		{
			description: "Get price",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(`
						(() => {
						const persianToLatin = {
							'۰': '0', '۱': '1', '۲': '2', '۳': '3', '۴': '4',
							'۵': '5', '۶': '6', '۷': '7', '۸': '8', '۹': '9'
						};
						const rows = document.querySelectorAll('.kt-base-row');
						for (const row of rows) {
							if (row.textContent.includes('قیمت کل')) {
								const priceEl = row.querySelector('.kt-unexpandable-row__value');
								if (!priceEl) return 0;
								const priceText = priceEl.textContent
									.replace('تومان', '')
									.replace(/,/g, '')
									.replace(/٬/g, '')
									.trim();
								const latinPrice = priceText
									.split('')
									.map(char => persianToLatin[char] || '')
									.join('');
								return parseInt(latinPrice) || 0;
							}
						}
						return 0;
					})()
					`, &ad.Price).Do(adCtx)
			},
		},
		{
			description: "Get images",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(`
						 Array.from(document.querySelectorAll('picture img')).map(img => {
							return img.src || img.getAttribute('data-src');
						}).filter(url => url && !url.includes('placeholder'))
					`, &ad.Images).Do(adCtx)
			},
		},
		{
			description: "Get creation date",
			action: func(adCtx context.Context) error {
				var title string
				err := chromedp.Evaluate(`document.title`, &title).Do(adCtx)
				if err != nil {
					return err
				}
				date, err := ExtractPersianDate(title)

				if err != nil {
					return err
				}
				ad.AdCreateDate = date
				return nil
			},
		},
	}
	return chromedp.Run(timeoutCtx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, task := range tasks {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					if err := task.action(ctx); err != nil {
						log.Printf("Error in %s for ad %s: %v", task.description, ad.Link, err)
					}
				}
			}
			return nil
		}),
	)
}

// SaveResults saves the crawled results to storage
func (c *Crawler) SaveResults(ads *[]HouseAd) error {
	if err := os.MkdirAll(c.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := filepath.Join(c.config.OutputDir,
		fmt.Sprintf("crawl_results_%s.json", time.Now().Format("2006-01-02_15-04-05")))

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(ads); err != nil {
		return fmt.Errorf("failed to encode results: %w", err)
	}

	log.Printf("Results saved to %s", filename)
	return nil
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

func evaluateNumericScript(selector string) string {
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

func uniqueAds(ads []HouseAd) []HouseAd {
	// return ads
	seen := make(map[string]bool)
	unique := []HouseAd{}

	for _, ad := range ads {
		key := ad.Title + "|" + ad.Link
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

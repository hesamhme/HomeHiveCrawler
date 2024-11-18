package crawler

import (
	"CrawlerProject/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	model "CrawlerProject/internal/model"
	utils "CrawlerProject/internal/utils"
	"CrawlerProject/pkg/config"
	"CrawlerProject/pkg/logger"

	"github.com/chromedp/chromedp"
	"golang.org/x/exp/rand"
)

type MyCrawler struct {
	model.Crawler
}

func NewCrawler(config model.CrawlerConfig) *MyCrawler {
	return &MyCrawler{
		Crawler: model.Crawler{
			Config:           config,
			UrlSemaphore:     make(chan struct{}, config.MaxURLConcurrency),
			AdsSemaphore:     make(chan struct{}, config.MaxAdConcurrency),
			ErrorChan:        make(chan error, len(config.Cities)*len(config.Types)),
			ResultsChan:      make(chan model.Listing, 10000),
			GoroutineMonitor: model.NewGoroutineMonitor(),
		},
	}
}

// DefaultConfig returns the default configuration
func DefaultConfig() model.CrawlerConfig {
	config, err := config.InitConfig()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error while initializing config")
		os.Exit(3)
	}
	return model.CrawlerConfig{
		RunInterval:        time.Duration(config.Interval) * time.Hour,
		MinTimeBetweenRuns: time.Duration(float64(5*time.Hour) * 0.9),
		PageTimeout:        30 * time.Minute,
		AdTimeout:          20 * time.Minute,
		MaxURLConcurrency:  config.MaxURLConcurrency,
		MaxAdConcurrency:   config.MaxAdConcurrency,
		Cities:             []string{"tehran"},
		Types:              []string{"buy-apartment", "buy-villa", "rent-apartment", "rent-villa"},
		// Cities:          []string{"tabriz", "azarshahr", "ahar", "bonab", "sarab", "sahand", "maragheh", "marand", "mianeh", "urmia", "oshnavieh", "bukan", "piranshahr", "khoy", "sardasht", "salmas", "shahin-dej", "maku", "mahabad", "miandoab", "naqadeh", "ardabil", "parsabad", "khalkhal", "sarein", "germi", "meshgin-shahr", "namin", "isfahan", "aran-va-bidgol", "abrisham-isfahan", "khomeyni-shahr", "khansar", "khour", "daran", "semirom", "shahin-shahr", "falavarjan", "foolad-shahr", "ghamsar", "kashan", "golpayegan", "lenjan", "mobarakeh", "najafabad", "karaj", "asara", "eshtehard", "tankaman", "charbagh-alborz", "taleqan", "fardis", "koohsar", "garmdareh", "mahdasht", "mohammad-shahr", "nazarabad", "hashtgerd", "abdanan", "ilam", "eyvan", "dehloran", "mehran", "borazjan", "dayyer", "bandar-kangan", "bandar-ganaveh", "bushehr", "jam", "khormoj", "tehran", "absard", "abali", "arjmand", "eslamshahr", "andisheh-new-town", "baghershahr", "bumehen", "pakdasht", "pardis", "parand", "pishva", "javadabad", "chahar-dangeh", "damavand", "robat-karim", "rudehen", "shahr-e-rey", "shahedshahr", "shemshak", "shahriar", "sabashahr", "safadasht-industrial-city", "ferdosiye", "fasham", "firuzkooh", "qods", "qarchak", "kahrizak", "kilan", "golestan-baharestan", "lavasan", "nasimshahr", "vahidieh", "varamin", "boroujen", "saman", "shahrekord", "farrokhshahr", "lordegan", "birjand", "tabas", "ferdows", "ghayen", "mashhad", "bardaskan", "taybad", "torbat-jam", "torbat-heydariyeh", "chenaran", "khaf", "sabzevar", "shandiz", "torghabeh", "qasemabad-khaf", "quchan", "golbahar", "gonabad", "molkabad", "neyshabur", "ashkhaneh", "esfarāyen", "bojnurd", "shirvan", "ahvaz", "abadan", "omidiyeh", "andimeshk", "izeh", "bandar-imam-khomeini", "bandar-mahshahr", "behbahan", "chamran-town", "hamidiyeh", "khorramshahr", "dezful", "ramshir", "ramhormoz", "susangerd", "shadeghan", "shush", "shooshtar", "masjed-soleyman", "hendijan", "abhar", "khorramdarreh", "zanjan", "qeydar", "damghan", "semnan", "shahroud", "garmsar", "iranshahr", "chabahar", "khash", "zabol", "zahedan", "zahak", "saravan", "konarak", "shiraz", "abadeh", "eqlid", "jahrom", "khoour", "darab", "zarghan", "sadra", "fasa", "firuzabad", "kazeroon", "lar", "lamerd", "marvdasht", "mohr", "norabad", "neyriz", "abyek", "eqbaliyeh", "alvand", "takestan", "shal", "qazvin", "mohammadiyeh", "qom", "baneh", "bijar", "dehgolan", "saqqez", "sanandaj", "qorveh", "kamyaran", "marivan", "baft", "bardsir", "boluk", "bam", "jiroft", "rafsanjan", "zarand", "sirjan", "kerman", "kahnooj", "mahan", "kermanshah", "eslamabad-gharb", "bisotun", "javanrud", "sarpol-zahab", "sonqor", "sahneh", "kangavar", "gahvareh", "harsin", "dogonbadan", "dehdasht", "sisakht", "yasuj", "azadshahr-golestan", "aq-qala", "bandar-torkaman", "aliabad-katul", "kordkuy", "kalale", "galikesh", "gorgan", "gomishan", "gonbad-kavus", "minoodasht", "sangdovin", "sorkhan-kalateh", "faragi", "sadegh-abad", "bandar-gaz", "maraveh-tapeh", "daland", "negin-shahr", "ramiyan", "khan-bin", "jelin", "dozin", "nokandeh", "goli-dagh", "nodeh-khandoz", "anbaralum", "fazel-abad", "mazrae-katool", "yanghagh", "sijval", "simin-shahr", "tatar-olya", "alghajar", "ghorogh", "inche-borun", "rasht", "astara", "astaneh-ashrafiyeh", "ahmadsar-gourab", "asalem", "amlash", "barah-sar", "bandar-anzali", "pareh-sar", "talesh", "toutkabon", "jirandeh", "chaboksar", "chaf-chamkhale", "chobar", "haviq", "khoshkbijar", "khomam", "deylaman", "rankouh", "rahim-abad", "rostam-abad", "rezvanshahr", "rudbar", "roudbaneh", "rudsar", "zibakenar", "sangar", "siahkal", "shaft", "shelman", "someh-sara", "fuman", "kelachay", "kouchesfahan", "koumeleh", "kiashahr", "gourab-zarmikh", "lahijan", "lashtenesha", "langarud", "loshan", "loulman", "lavandevil", "lisar", "masal", "masuleh", "makloan", "manjil", "vajargah", "tahergurab", "shanderman", "ziyabar", "otaghvar", "tulam-shahr", "pirbazar", "azna", "aleshtar", "aligudarz", "borujerd", "pol-dokhtar", "khorramabad", "dorud", "kuhdasht", "nurabad", "aalasht", "amol", "amirkala", "izadshahr", "babol", "babolsar", "baladeh", "behshahr", "bahnamir", "polsefid", "tonekabon", "juybar", "chalus", "chamestan", "khalil-shahr", "khoshroud-pey", "ramsar", "rostamkola", "royan", "reyneh", "ziraab", "sari", "sorkhrood", "salman-shahr", "sourek", "shirgah", "abbasabad-mazandaran", "farahabad", "fereydunkenar", "farim", "qaemshahr", "katalem-sadatshahr", "kelarabad", "kelarestan", "kouhi-kheyl", "kiasar", "kiakola", "gatab", "gazanak", "galougah-babol", "mahmudabad", "marzan-abad", "marzikola", "nashtarud", "neka", "nur", "nowshahr", "paeen-holar", "dalkhani", "galugah-babol", "hadi-shahr", "babakan", "zargarshahr", "arateh", "emamzadeh-abdollah", "shirud", "dabudasht", "akand", "astaneh-sara", "pool", "tabaghdeh", "kojur", "khoram-abad", "hachirud", "arak", "khomein", "delijan", "saveh", "shazand", "mahalat", "mohajeran", "bandar-abbas", "takht", "dargahan", "qeshm", "kish", "minab", "hormuz", "asadabad", "bahar", "tuyserkan", "kabudrahang", "malayer", "nahavand", "hamedan", "ardakan", "bafq", "taft", "hamidia", "mehriz", "meybod", "yazd"},
		// Types: 			[]string{"buy-apartment"},
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
func (c *MyCrawler) Start(ctx context.Context) error {
	log.Printf("Starting crawler with interval: %v", c.Config.RunInterval)

	// Run immediately on startup
	if err := c.RunOnce(ctx); err != nil {
		log.Printf("Initial run failed: %v", err)
	}

	// Setup ticker for periodic runs
	ticker := time.NewTicker(c.Config.RunInterval)
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
func (c *MyCrawler) checkAndRun(ctx context.Context) error {
	c.RunMutex.Lock()
	if time.Since(c.LastRunTime) < c.Config.MinTimeBetweenRuns {
		c.RunMutex.Unlock()
		return nil
	}
	c.RunMutex.Unlock()
	return c.RunOnce(ctx)
}

// RunOnce performs a single crawl operation
func (c *MyCrawler) RunOnce(ctx context.Context) error {
	c.RunMutex.Lock()
	defer c.RunMutex.Unlock()

	log.Printf("Starting crawl at %v", time.Now())
	defer func() {
		c.LastRunTime = time.Now()
		log.Printf("Completed crawl at %v", time.Now())
	}()

	return c.crawl(ctx)
}

// crawl performs the actual crawling operation
func (c *MyCrawler) crawl(ctx context.Context) error {
	// Setup browser context
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, c.Config.ChromeFlags...)
	defer allocCancel()

	// Create timeout context for entire operation
	crawlCtx, cancel := context.WithTimeout(allocCtx, c.Config.PageTimeout)
	defer cancel()

	var wg sync.WaitGroup
	var allAds []model.Listing

	// Process each URL concurrently
	for _, city := range c.Config.Cities {
		for _, _type := range c.Config.Types {
			wg.Add(1)
			go func(city, _type string) {
				defer wg.Done()

				// Start monitoring this goroutine

				stats := c.GoroutineMonitor.StartTracking(city, _type)
				defer c.GoroutineMonitor.StopTracking(stats.GoroutineID)

				if err := c.processURL(crawlCtx, city, _type, stats, &allAds); err != nil {
					select {
					case c.ErrorChan <- err:
					default:
					}
				}
			}(city, _type)
		}
	}

	wg.Wait()

	// Save goroutine statistics
	if err := c.GoroutineMonitor.SaveStats(c.Config.OutputDir); err != nil {
		log.Printf("Error saving goroutine stats: %v", err)
	}

	// Process gathered ads
	return c.processAds(crawlCtx, &allAds)
}

// processURL handles crawling a single URL
func (c *MyCrawler) processURL(ctx context.Context, city, _type string, stats *model.GoroutineStats, allAds *[]model.Listing) error {
	// Acquire URL semaphore
	c.UrlSemaphore <- struct{}{}
	defer func() { <-c.UrlSemaphore }()

	url := "https://divar.ir/s/" + city + "/" + _type
	stats.URL = url

	// Create new browser context
	browserCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var urlAds []model.Listing
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
	c.AdsMutex.Lock()
	*allAds = append(*allAds, urlAds...)
	c.AdsMutex.Unlock()

	log.Printf("Completed URL %s: Found %d ads", url, len(urlAds))
	return nil
}

// scrollAndScrape implements the scrolling and scraping logic
// Original function with issues identified and fixed
func (c *MyCrawler) scrollAndScrape(ads *[]model.Listing, wg *sync.WaitGroup) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		// Only one Done() call is needed at the end
		defer wg.Done()

		var previousHeight int64
		adChannel := make(chan []model.Listing)
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

				var newAds []model.Listing
				err := chromedp.Evaluate(`
                    Array.from(document.querySelectorAll('.kt-post-card')).map(card => ({
                        title: card.querySelector('.kt-post-card__title')?.innerText || '',
                        url: card.querySelector('a')?.href || ''
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
					*ads = utils.UniqueAds(append(*ads, newAds...))
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
func (c *MyCrawler) processAds(ctx context.Context, ads *[]model.Listing) error {
	totalAds := len(*ads)
	if totalAds == 0 {
		return fmt.Errorf("no ads found during scraping")
	}

	log.Printf("Processing details for %d ads", totalAds)

	var wg sync.WaitGroup
	for i := range *ads {
		// for i := range 6 {
		wg.Add(1)
		go func(ad *model.Listing, index int) {
			defer wg.Done()

			// Reactivate semaphore control
			c.AdsSemaphore <- struct{}{}
			defer func() { <-c.AdsSemaphore }()

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
				case c.ErrorChan <- fmt.Errorf("failed after %d retries: %w", maxRetries, err):
				default:
				}
				return
			}

			select {
			case c.ResultsChan <- *ad:
			case <-ctx.Done():
			}
		}(&(*ads)[i], i)
	}

	go func() {
		wg.Wait()
		close(c.ErrorChan)
		close(c.ResultsChan)
	}()

	var processedAds []model.Listing
	var errors []error

	for c.ErrorChan != nil || c.ResultsChan != nil {
		select {
		case err, ok := <-c.ErrorChan:
			if !ok {
				c.ErrorChan = nil
				continue
			}
			errors = append(errors, err)
		case ad, ok := <-c.ResultsChan:
			if !ok {
				c.ResultsChan = nil
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
func (c *MyCrawler) processAdDetails(ctx context.Context, ad *model.Listing, index int) error {
	// Create new browser context for each ad
	fmt.Println("crawling ", ad.URL)
	browserCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Add timeout
	timeoutCtx, timeoutCancel := context.WithTimeout(browserCtx, c.Config.AdTimeout)
	defer timeoutCancel()

	// Navigate to ad page first
	if err := chromedp.Run(timeoutCtx, chromedp.Navigate(ad.URL)); err != nil {
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
				err := chromedp.Evaluate(utils.EvaluateNumericScript(
					"document.querySelectorAll('.kt-group-row__data-row .kt-group-row-item__value')[0]",
				), &ad.Area).Do(adCtx)
				// Fall back to secondary selector if primary fails
				if err != nil || ad.Area == 0 {
					return chromedp.Evaluate(utils.EvaluateNumericScript(
						"document.querySelector('.kt-unexpandable-row__value')",
					), &ad.Area).Do(adCtx)
				}
				return err
			},
		},
		{
			description: "Get bedrooms",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(utils.EvaluateNumericScript(
					"document.querySelector('.kt-group-row__data-row td:nth-child(3)')",
				), &ad.Rooms).Do(adCtx)
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
							var city = parts[1].split('،'); 
							return city.length > 1 ? city[0].trim() : '';
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
					`, &ad.Status).Do(adCtx)
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
				ad.HasElevator = amenities.HasElevator
				ad.HasStorage = amenities.HasWarehouse
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
				date, err := utils.ExtractPersianDate(title)

				if err != nil {
					return err
				}
				ad.CreatedAt = date
				return nil
			},
		},
		{
			description: "Get neighbourhood",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(`
						const text = document.querySelector('.kt-page-title__subtitle').textContent;
						const parts = text.split(/[,،]/);
						parts.length > 1 ? parts[1].trim() : '';
					`, &ad.Neighborhood).Do(adCtx)
			},
		},
		{
			description: "Get parking",
			action: func(adCtx context.Context) error {
				return chromedp.Evaluate(`
						 (() => {
						try {
							// First find the section title
							const sectionTitle = Array.from(document.querySelectorAll('.kt-section-title__title'))
								.find(el => el.textContent === 'ویژگی‌ها و امکانات');
							
							if (!sectionTitle) {
								return false;
							}
							
							// Find the closest table that contains features
							const featureTable = sectionTitle.closest('.kt-section-title--alt-padded')
								.nextElementSibling;
							
							if (!featureTable) {
								return false;
							}
							
							// Look for پارکینگ in the table
							const hasParking = Array.from(featureTable.querySelectorAll('.kt-body--stable'))
								.some(el => el.textContent === 'پارکینگ');
								
							return hasParking;
						} catch (e) {
							console.error('Error checking for parking:', e);
							return false;
						}
					})();
					`, &ad.HasParking).Do(adCtx)
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
						log.Printf("Error in %s for ad %s: %v", task.description, ad.URL, err)
					}
				}
			}
			return nil
		}),
	)
}

// SaveResults saves the crawled results to storage
func (c *MyCrawler) SaveResults(ads *[]model.Listing) error {
	// save to database
	for _, ad := range *ads {
		service.StoreListing(nil, ad)
	}

	if err := os.MkdirAll(c.Config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := filepath.Join(c.Config.OutputDir,
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

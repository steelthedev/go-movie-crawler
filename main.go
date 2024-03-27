package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

type CrawlerHandler struct {
	Collector    *colly.Collector
	PageLimit    int
	PageToScrape string
	Update       tgbotapi.UpdatesChannel
	Bot          *tgbotapi.BotAPI
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func NewCrawlerHandler(collector *colly.Collector, pageLimit int, pageToScrape string, update tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) *CrawlerHandler {
	return &CrawlerHandler{
		Collector:    collector,
		PageLimit:    pageLimit,
		PageToScrape: pageToScrape,
		Update:       update,
		Bot:          bot,
	}

}

func (h CrawlerHandler) ScrapeMovieData() ([]string, error) {
	var data []string
	h.Collector.OnHTML("div.titles-many", func(h *colly.HTMLElement) {
		links := h.ChildAttrs("a", "href")
		data = append(data, links...)
	})

	h.Collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := h.Collector.Visit(h.PageToScrape)

	if err != nil {
		slog.Error("could not visit website")
		return nil, err
	}

	return data, nil
}

func (h *CrawlerHandler) InitCrawler() {
	var pagesToScrape []string
	i := 1

	pagesDiscovered := []string{h.PageToScrape}

	h.Collector.OnHTML("ul.pagination", func(e *colly.HTMLElement) {
		var nextPage string
		nextPageLinks := e.ChildAttrs("a", "href")

		if len(nextPageLinks) < 2 && len(nextPageLinks) != 0 {
			nextPage = nextPageLinks[0]
		}
		if len(nextPageLinks) != 0 && len(nextPageLinks) > 1 {
			nextPage = nextPageLinks[1]
		}

		// Check if the discovered page is new and should be scraped
		if !contains(pagesToScrape, nextPage) && !contains(pagesDiscovered, nextPage) {
			pagesToScrape = append(pagesToScrape, nextPage)
			pagesDiscovered = append(pagesDiscovered, nextPage)
		}
	})

	h.Collector.OnScraped(func(r *colly.Response) {
		// keep scraping data till there are no new pages again
		for len(pagesToScrape) != 0 && i < h.PageLimit {
			// Getting pages to scrap, adding them to page that should be scrapped and removing them from the pagesToScrape list
			h.PageToScrape = pagesToScrape[0]

			pagesToScrape = pagesToScrape[1:]

			i++
			h.Collector.Visit(h.PageToScrape)
		}
	})
	// Scrap logic (Returns a slice of data containing several links)
	data, err := h.ScrapeMovieData()
	if err != nil {
		slog.Error("Could not finish %s", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(data))

	for _, link := range data {
		go func(link string) {
			defer wg.Done()
			fmt.Println("Sending messages")
			msg := tgbotapi.NewMessageToChannel("@awafim_crawler_bot", link)
			_, err := h.Bot.Send(msg)
			if err != nil {
				slog.Error(err.Error())
			}
		}(link)

	}
	wg.Wait()

	h.Collector.Visit(h.PageToScrape)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(string(os.Getenv("BOT_TOKEN")))
	if err != nil {
		slog.Info(err.Error())
	}
	bot.Debug = true
	fmt.Println("Authorized ", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	c := colly.NewCollector()
	handler := NewCrawlerHandler(c, 5, "https://www.awafim.tv/browse/page/1", updates, bot)
	handler.InitCrawler()

}

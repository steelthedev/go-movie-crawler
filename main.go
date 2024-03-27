package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gocolly/colly"
)

type CrawlerHandler struct {
	Collector    *colly.Collector
	PageLimit    int
	PageToScrape string
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func NewCrawlerHandler(collector *colly.Collector, pageLimit int, pageToScrape string) *CrawlerHandler {
	return &CrawlerHandler{
		Collector:    collector,
		PageLimit:    pageLimit,
		PageToScrape: pageToScrape,
	}

}

func (h CrawlerHandler) ScrapeMovieData() ([]string, error) {
	time.Sleep(30 * time.Second)
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

func (h CrawlerHandler) InitCrawler() {
	var pagesToScrape []string
	i := 1

	pagesDiscovered := []string{h.PageToScrape}

	h.Collector.OnHTML("ul.pagination", func(h *colly.HTMLElement) {
		nextPage := h.ChildAttr("a", "href")
		fmt.Println(nextPage)

		// check if the dicovered page is new
		if !contains(pagesToScrape, nextPage) {

			// check if the discovered page should be scrapped
			if !contains(pagesDiscovered, nextPage) {
				pagesToScrape = append(pagesToScrape, nextPage)
			}
			pagesDiscovered = append(pagesDiscovered, nextPage)
		}
	})

	// Scrap logic (Returns a slice of data containing several links)
	_, err := h.ScrapeMovieData()
	if err != nil {
		slog.Error("Could not finish %s", err)
		return
	}

	h.Collector.OnScraped(func(r *colly.Response) {
		// keep scrapping data till there are no new pages again

		if len(pagesToScrape) != 0 && i < h.PageLimit {
			// Getting pages to scrap, adding them to page that should be scrapped and removing them from the pagesToScrap list
			h.PageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			i++
			h.Collector.Visit(h.PageToScrape)
		}
	})

	h.Collector.Visit(h.PageToScrape)
}

func main() {

	c := colly.NewCollector()
	handler := NewCrawlerHandler(c, 5, "https://www.awafim.tv/browse/page/1")
	handler.InitCrawler()

}

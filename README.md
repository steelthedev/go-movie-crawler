### Crawler Readme

This Go package is designed to scrape movie data from a https://awafim.tv and send the collected information to a Telegram channel using a bot. Below is an overview of the package and its functionalities.

#### Overview

The package consists of several components:

- **Main**: The main function initializes the bot, sets up the crawler handler, and starts the scraping process.
- **CrawlerHandler**: This struct encapsulates the functionality related to crawling and scraping data. It utilizes the Colly library for web scraping and the Telegram Bot API for sending messages.
- **Functions**:
  - `NewCrawlerHandler`: Initializes a new crawler handler with the necessary parameters.
  - `ScrapeMovieData`: Scrapes movie data from the specified webpage.
  - `GetTargetElements`: Retrieves pagination links for further scraping.
  - `onScrapeData`: Initiates the scraping process for each discovered page.
  - `sendTelegramMessageToChannel`: Sends messages containing scraped data to a specified Telegram channel.
  - `InitCrawler`: Initializes the crawling process, including scraping and message sending.

#### Usage

To use this package, follow these steps:

1. Ensure you have a valid Telegram bot token and set it in the `.env` file.
2. Import the necessary packages (`github.com/go-telegram-bot-api/telegram-bot-api/v5`, `github.com/gocolly/colly`, `github.com/joho/godotenv`).
3. Initialize a new bot using the provided token.
4. Create a new Colly collector.
5. Initialize a new CrawlerHandler instance with the collector, page limit, starting URL, updates channel, and bot.
6. Call the `InitCrawler` function to start the scraping process.

#### Dependencies

- `github.com/go-telegram-bot-api/telegram-bot-api/v5`: Telegram Bot API library for Go.
- `github.com/gocolly/colly`: Elegant scraping framework for Golang.
- `github.com/joho/godotenv`: GoDotEnv loads environment variables from a `.env` file.



#### Note

- Ensure that the `.env` file contains the correct `BOT_TOKEN` variable.
- Modify the `PageLimit` and `PageToScrape` parameters in `NewCrawlerHandler` according to your requirements.
- Adjust the target Telegram channel in `sendTelegramMessageToChannel` function.
- Customize the scraping logic in `ScrapeMovieData` and other relevant functions as needed.

Feel free to extend and modify this package according to your specific use case.
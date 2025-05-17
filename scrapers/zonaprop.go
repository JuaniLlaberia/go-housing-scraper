package scrapers

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
	"zonaprop-scraper/structs"
	"zonaprop-scraper/utils"

	"github.com/gocolly/colly"
)

func ZonapropScraper(c *colly.Collector, url structs.Url) {
	var properties []structs.Property
	var visitedUrls sync.Map

	c.OnHTML("h1.postingsTitle-module__title", func(e *colly.HTMLElement) {
		title := (strings.Split(e.Text, " ")[0])
		amount, err := strconv.Atoi(title)
		if err != nil {
			fmt.Println("Failed to parse amount of properties.")
		}

		fmt.Printf("Preparing to scrape %v pages...\n", math.Ceil(float64(amount)/30.0))
	})

	c.OnHTML("div.postingCard-module__posting-top", func(e *colly.HTMLElement) {
		property := structs.Property{}

		priceStr := e.ChildText("div.postingPrices-module__price")
		parts := strings.Fields(priceStr)
		if len(parts) < 2 {
			fmt.Println("Invalid price format")
			return
		}

		price := strings.ReplaceAll(parts[1], ".", "")
		property.Price = price

		address := e.ChildText("div.postingLocations-module__location-address")
		formattedAddress := utils.CleanAddress(address, url.Neighbour)
		property.Address = formattedAddress

		squareStr := e.ChildText("postingMainFeatures-module__posting-main-features-listing")
		parts = strings.Fields(squareStr)
		if len(parts) < 2 {
			fmt.Println("Invalid square number format")
			return
		}

		square := strings.ReplaceAll(parts[1], ".", "")
		property.Square = square

		numPrice, errPrice := strconv.ParseInt(price, 10, 64)
		numSquare, errSquare := strconv.ParseInt(square, 10, 64)

		if errPrice == nil && errSquare == nil {
			property.PricePerSquare = fmt.Sprintf("%v", numPrice/numSquare)
		} else {
			fmt.Println("Failed to convert price and square to number. So price per square was not calculated.")
		}

		url := e.ChildAttr("h3.postingCard-module__posting-description a", "href")
		property.Url = "https://www.zonaprop.com.ar" + url

		properties = append(properties, property)
	})

	c.OnHTML("[data-qa='PAGING_NEXT']", func(e *colly.HTMLElement) {
		url.Page += 1
		pageUrl := utils.UrlBuilder(url)

		if _, found := visitedUrls.Load(pageUrl); !found {
			delay := rand.Intn(5) + 1
			fmt.Printf("Delaying requests %v seconds to avoid anti-bot system\n", delay)
			time.Sleep(time.Duration(delay) * time.Second)

			visitedUrls.Store(pageUrl, struct{}{})
			e.Request.Visit(pageUrl)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		// Sort properties by price
		utils.SortPropertiesByPrice(&properties)

		// Create CSV file from properties
		err := utils.WriteCsv("propiedades.csv", properties)
		if err != nil {
			fmt.Println("Failed to create CSV file. Error message:", err.Error())
			fmt.Println("Printing results:", properties)
		}
		fmt.Println("Scrape finished correctly.")
	})

	urlToScrape := utils.UrlBuilder(url)
	c.Visit(urlToScrape)
}

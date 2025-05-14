package main

import (
	"fmt"
	"math"
	"math/rand"
	"scraper/structs"
	"scraper/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("=========================")
	fmt.Println("Starting scraping process")
	fmt.Println("=========================")

	var properties []structs.Property
	var visitedUrls sync.Map

	url := structs.Url{
		Base:         "https://www.zonaprop.com.ar",
		PropertyType: structs.Departamentos,
		Operation:    structs.Venta,
		Neighbour:    "Belgrano",
		PriceRange: structs.PriceRange{
			Min: 600_000,
			Max: 800_000,
		},
		Areas:          3,
		Rooms:          2,
		Bathrooms:      1,
		ProfesionalUse: false,
		Page:           1,
	}

	c := colly.NewCollector(
		colly.AllowedDomains("www.zonaprop.com.ar"),
	)

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.RandomUserAgent())
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept-Language", "es-AR")
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visiting:", r.Request.URL)
	})

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

		property.Price = e.ChildText("div.postingPrices-module__price")
		property.Address = e.ChildText("div.postingLocations-module__location-address")
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
		err := utils.WriteCsv("propiedades.csv", properties)
		if err != nil {
			fmt.Println("Failed to create CSV file. Error message:", err.Error())
			fmt.Println("Printing results:", properties)
		}
	})

	urlToScrape := utils.UrlBuilder(url)

	// Start by scraping the first page
	// => Giving us the properties of page #1 and the total properties so we can do parallelism
	c.Visit(urlToScrape)
}

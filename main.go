package main

import (
	"fmt"
	"io"
	"log"
	"net/http/cookiejar"
	"os"
	"scraper/analyzers"
	"scraper/scrapers"
	"scraper/structs"
	"scraper/utils"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("=========================")
	fmt.Println("Starting scraping process")
	fmt.Println("=========================")

	// Defining url and filters
	url := structs.Url{
		Base:         "https://www.zonaprop.com.ar",
		PropertyType: structs.Departamentos,
		Operation:    structs.Venta,
		Neighbour:    "Belgrano",
		PriceRange: structs.PriceRange{
			Min: 600_000,
			Max: 700_000,
		},
		Areas:          3,
		Rooms:          2,
		Bathrooms:      1,
		ProfesionalUse: false,
		Page:           1,
	}

	c := colly.NewCollector(
		colly.AllowedDomains("www.zonaprop.com.ar", "www.zonaprop.com.ar:443"),
	)

	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	if jar != nil {
		c.SetCookieJar(jar)
	}

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err.Error())
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.RandomUserAgent())
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Referer", "https://www.google.com/")
	})

	// Gettings the data
	scrapers.ZonapropScraper(c, url)

	// Analyze data with Gemini
	csvFile, err := os.Open("propiedades.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	bytes, err := io.ReadAll(csvFile)
	if err != nil {
		fmt.Println("Failed to convert file to []byte")
	}

	analyzers.BasicPropertiesAnalyzer(bytes)
}

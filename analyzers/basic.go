package analyzers

import (
	"encoding/json"
	"fmt"
	"os"
	"scraper/configs"
	"scraper/structs"

	"google.golang.org/genai"
)

type MarketOverview struct {
	AveragePrice          int64 `json:"average_price"`
	MedianPrice           int64 `json:"median_price"`
	AveragePricePerSquare int64 `json:"average_price_per_square"`
}

type AnalyzerResponse struct {
	MarketOverview      MarketOverview     `json:"market_overview"`
	Summary             string             `json:"summary"`
	BestValueProperties []structs.Property `json:"best_value_properties"`
}

func PropertiesAnalyzer(file []byte, url structs.Url) {
	prompt := fmt.Sprintf(`
		Analyze the provided CSV data from properties. Here is some extra data about the properties that may help you:
			- Areas: %v.
			- Bedrooms: %v.
			- Bathrooms: %v or more.
		(All properties have the same characteristics in terms of rooms, areas, etc... What changes is price and square meter).

		After analysing the data provide:
			1) A proffesional summary report (between 300-500 words) for a real state agent to use with clients.
			   Focus on price insigths and best-value observations. Example:
			   		"This week’s listing summary includes 85 properties, with an average price of $320,000 and an average size of 68 sqm. The average price per square meter is $4,700, but there are notable opportunities below this.
					The best value listings include a 75 sqm unit at $2,800/sqm, which is nearly 40%% below market average. In total, 14 listings are priced more than 25%% below average, indicating strong potential for buyers looking for value.
					Price generally increases with square footage (r² = 0.87), suggesting consistent pricing logic across the dataset. However, larger apartments tend to offer slightly better value per sqm.
					The top 3 value deals this week are all under $3,000/sqm, making them attractive options for both investors and home buyers seeking affordability."

			2) Return a ranked list of the 10 best-value properties, based on lowest price per square meter.

			3) Provide some properties overview values (average price, median price and average price per square meter)

			Follow this JSON format:
				"market_overview": {
					"average_price": int64,
					"median_price": int64,
					"average_price_per_square": int64,
				},
				"summary": string,
				"best_value_properties": [
					{
						"address": string,
						"price": string,
						"square": string,
						"pricePerSquare": string,
						"url": string
					}
				],
	`, url.Areas, url.Rooms, url.Bathrooms)

	parts := []*genai.Part{
		{
			InlineData: &genai.Blob{
				MIMEType: "text/csv",
				Data:     file,
			},
		},
		genai.NewPartFromText(prompt),
	}
	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := configs.Gemini(contents)
	if err != nil {
		fmt.Printf("Something went wrong with gemini: %v", err.Error())
		return
	}

	var response AnalyzerResponse
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		fmt.Printf("Something went wrong parsing gemini response: %v", err.Error())
		return
	}

	data, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err.Error())
		return
	}

	fmt.Println("Creating json file with report...")
	err = os.WriteFile("report.json", data, 0644)
	if err != nil {
		fmt.Printf("Failed to write JSON to file: %v\n", err.Error())
		return
	}

	fmt.Println("Gemini has analyzed the properties and created a report.")
}

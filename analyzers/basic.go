package analyzers

import (
	"fmt"
	"scraper/configs"

	"google.golang.org/genai"
)

func BasicPropertiesAnalyzer(file []byte) {
	prompt := `
		Analyze the provided CSV data from properties. 
		I want you to make a report including:
			- Top 10 most affordable properties.
			- Top 10 properties with the best price per square foot.

		In addition generate some statistics about all properties:
			- Average price
			- Median price
			- Average price per sqft

		Follow this JSON format:
			"market_overview": {
				"average_price": int64,
				"median_price": int64,
				"average_price_per_sqft": int64,
    		},
			"top_10_best_value_properties": [
				{
					"address": string,
					"price": string,
					"square": string,
					"url": string
				}
			],
			"top_10_most_affortable_properties": [
				{
					"address": string,
					"price": string,
					"square": string,
					"url": string
				}
			]

	`

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

	fmt.Println(result)
}

package utils

import (
	"encoding/csv"
	"os"
	"zonaprop-scraper/structs"
)

func WriteCsv(path string, properties []structs.Property) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"Address",
		"Price",
		"Square",
		"PricePerSquare",
		"Url",
	}
	writer.Write(headers)

	for _, property := range properties {
		record := []string{
			property.Address,
			property.Price,
			property.Square,
			property.PricePerSquare,
			property.Url,
		}

		writer.Write(record)
	}
	defer writer.Flush()

	return nil
}

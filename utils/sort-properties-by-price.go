package utils

import (
	"sort"
	"zonaprop-scraper/structs"
)

func SortPropertiesByPrice(properties *[]structs.Property) {
	sort.Slice(*properties, func(i, j int) bool {
		return (*properties)[i].Price < (*properties)[j].Price
	})
}

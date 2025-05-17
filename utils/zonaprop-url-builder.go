package utils

import (
	"fmt"
	"strings"
	"zonaprop-scraper/structs"
)

func UrlBuilder(url structs.Url) string {
	baseUrl := url.Base + "/"

	// Property Type
	if url.PropertyType == structs.Casas {
		baseUrl += "casas-"
	} else if url.PropertyType == structs.Departamentos {
		baseUrl += "departamentos-"
	} else {
		baseUrl += "inmuebles-"
	}

	// Operation
	if url.Operation == structs.Alquiler {
		baseUrl += "alquiler-"
	} else if url.Operation == structs.Venta {
		baseUrl += "venta-"
	} else if url.Operation == structs.Temporal {
		baseUrl += "alquiler-temporal-"
	}

	// Format neighbour
	if url.Neighbour != "" {
		neighbour := strings.Join(strings.Split(url.Neighbour, " "), "-")
		lowerNeighbour := strings.ToLower(neighbour)

		baseUrl += fmt.Sprintf("%v-", lowerNeighbour)
	}

	if url.ProfesionalUse {
		baseUrl += "con-apto-profesional-"
	}

	// Bathrooms
	if url.Bathrooms != 0 {
		baseUrl += fmt.Sprintf("mas-de-%v-bano-", url.Bathrooms)
	}

	// Rooms
	if url.Rooms != 0 {
		baseUrl += fmt.Sprintf("%v-habitaciones-", url.Rooms)
	}

	// Areas
	if url.Areas != 0 {
		baseUrl += fmt.Sprintf("%v-ambientes-", url.Areas)
	}

	// Price range
	if url.PriceRange.Min != 0 {
		baseUrl += fmt.Sprintf("%v-", url.PriceRange.Min)

		if url.PriceRange.Max != 0 {
			baseUrl += fmt.Sprintf("%v-", url.PriceRange.Max)
		}

		if url.Currency == structs.Dollar {
			baseUrl += "dolar-"
		} else if url.Currency == structs.Pesos {
			baseUrl += "pesos-"
		}
	}

	if url.Page != 1 {
		baseUrl += fmt.Sprintf("pagina-%v", url.Page)
	}

	return baseUrl + ".html"
}

package structs

type PropertyType int

const (
	Departamentos PropertyType = iota
	Casas
)

type OperationType int

const (
	Venta OperationType = iota
	Alquiler
	Temporal
)

type Currency int

const (
	Dollar Currency = iota
	Pesos
)

type PriceRange struct {
	Min int64
	Max int64
}

type Url struct {
	Base           string
	PropertyType   PropertyType
	Operation      OperationType
	Neighbour      string
	PriceRange     PriceRange
	Currency       Currency
	Areas          int64
	Rooms          int64
	Bathrooms      int64
	ProfesionalUse bool
	Page           int64
}

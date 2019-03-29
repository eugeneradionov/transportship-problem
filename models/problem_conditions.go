package models

type ProblemConditions struct {
	Suppliers     []Supplier  `json:"suppliers"`
	Consumers     []Consumer  `json:"consumers"`
	TransportCost [][]float64 `json:"transport_cost"`
}

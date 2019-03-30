package models

type Supplier struct {
	ID      int     `json:"id"`
	Stock   float64 `json:"stock"`
	fictive bool
}

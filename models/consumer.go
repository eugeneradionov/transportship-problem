package models

type Consumer struct {
	ID      int     `json:"id"`
	Demand  float64 `json:"demand"`
	fictive bool
}

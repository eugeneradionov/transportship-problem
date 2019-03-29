package models

type Consumer struct {
	ID     int     `json:"id,omitempty"`
	Demand float64 `json:"demand"`
}

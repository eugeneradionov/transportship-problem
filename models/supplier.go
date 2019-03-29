package models

type Supplier struct {
	ID    int     `json:"id,omitempty"`
	Stock float64 `json:"stock"`
}

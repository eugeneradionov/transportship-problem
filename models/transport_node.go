package models

import "fmt"

type transportNode struct {
	Sup    *Supplier `json:"from,omitempty"`
	Cons   *Consumer `json:"to,omitempty"`
	Amount float64   `json:"amount"`
}

func (t transportNode) String() string {
	return fmt.Sprintf("Sup{ID: %v}, Cons{ID: %v}, Amount: %f", t.Sup.ID, t.Cons.ID, t.Amount)
}

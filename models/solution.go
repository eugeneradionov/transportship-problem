package models

import (
	"errors"
	"fmt"
	"log"
	"math"
)

const elipsis = 0.00001

type Solution struct {
	ProblemConditions `json:"-"`

	Route []transportNode `json:"transport_route"`
	Cost  float64         `json:"total_cost"`

	pivotN  int
	pivotM  int
	route   [][]float64
	pots    [][]float64
	numSup  int
	numCons int
}

func NewSolution(cond ProblemConditions) *Solution {
	s := Solution{ProblemConditions: cond}
	s.pivotN = -1
	s.pivotM = -1

	return &s
}

func (s *Solution) Solve() error {
	if err := s.checkInputValues(); err != nil {
		return err
	}

	s.fixImbalance()
	s.initData()

	s.northWest()
	log.Printf("[INFO] at Solution.Solve(): northWest completed: %v", s.route)

	nOpt, err := s.notOptimal()
	if err != nil {
		return err
	}

	for nOpt {
		err := s.findOptimal()
		if err != nil {
			return err
		}

		nOpt, err = s.notOptimal()
		if err != nil {
			return err
		}
	}
	log.Printf("[INFO] at Solution.Solve(): Found optimal solution: dual: %v, route: %v", s.pots, s.route)

	s.restoreInitialConditions()

	s.createRoute()
	log.Printf("[INFO] at Solution.Solve(): Got solution: Route:{%v}", s.Route)

	return nil
}

func (s *Solution) fixImbalance() {
	balance := 0.0

	for _, sup := range s.Suppliers {
		balance += sup.Stock
	}

	for _, cons := range s.Consumers {
		balance -= cons.Demand
	}

	if balance > 0 {
		s.Consumers = append(s.Consumers, Consumer{ID: -1, Demand: balance, fictive: true})
	} else if balance < 0 {
		s.Suppliers = append(s.Suppliers, Supplier{ID: -1, Stock: -balance, fictive: true})
	}
}

func (s *Solution) checkInputValues() error {
	if len(s.Suppliers) <= 0 {
		return errors.New("there is must be at least one supplier")
	}

	if len(s.Consumers) <= 0 {
		return errors.New("there is must be at least one consumer")
	}

	if len(s.TransportCost) != len(s.Suppliers) {
		return fmt.Errorf("invalid transport cost matrix: must have %v rows, got %v", len(s.Suppliers), len(s.TransportCost))
	}

	for _, row := range s.TransportCost {
		if len(row) != len(s.Consumers) {
			return fmt.Errorf("invalid transport cost matrix row: must have %v columns, got %v", len(s.Consumers), len(row))
		}
	}

	return nil
}

func (s *Solution) initData() {
	s.numSup = len(s.Suppliers)
	s.numCons = len(s.Consumers)

	for i := range s.Consumers {
		s.Consumers[i].Demand += elipsis
	}
	s.Suppliers[s.numSup-1].Stock += elipsis * float64(s.numCons)

	newTransportCost := make([][]float64, s.numSup)
	for i := range newTransportCost {
		newTransportCost[i] = make([]float64, s.numCons)
	}

	for i, row := range s.TransportCost {
		for j, cost := range row {
			newTransportCost[i][j] = cost
		}
	}
	s.TransportCost = newTransportCost

	s.route = make([][]float64, s.numSup)
	for i := range s.route {
		s.route[i] = make([]float64, s.numCons)
	}

	s.pots = make([][]float64, s.numSup)
	for i := range s.pots {
		sl := make([]float64, s.numCons)
		for i := range sl {
			sl[i] = -1.0
		}
		s.pots[i] = sl
	}

}

func (s *Solution) northWest() {
	var i, j int

	sup := make([]float64, s.numCons)
	cons := make([]float64, s.numSup)

	for i <= s.numSup-1 && j <= s.numCons-1 {
		if s.Consumers[j].Demand-sup[j] < s.Suppliers[i].Stock-cons[i] {
			delta := s.Consumers[j].Demand - sup[j]
			s.route[i][j] = delta
			sup[j] += delta
			cons[i] += delta
			j += 1
		} else {
			delta := s.Suppliers[i].Stock - cons[i]
			s.route[i][j] = delta
			sup[j] += delta
			cons[i] += delta
			i += 1
		}
	}
}

func (s *Solution) findPath(i, j int) ([][]int, error) {
	path := [][]int{
		{i, j},
	}

	if !s.findHorizontally(&path, i, j, i, j) {
		log.Printf("solution.findPath(): path error, cannot find path horizontally, path = %v, i = %v, j = %v", path, i, j)
		return nil, errors.New("path error: cannot find path horizontally")
	}

	return path, nil
}

func (s *Solution) findHorizontally(path *[][]int, u, v, u1, v1 int) bool {
	for i := 0; i < s.numCons; i++ {
		if i != v && s.route[u][i] != 0 {
			if i == v1 {
				*path = append(*path, []int{u, i})
				return true
			}

			if s.findVertically(path, u, i, u1, v1) {
				*path = append(*path, []int{u, i})
				return true
			}
		}
	}

	return false
}

func (s *Solution) findVertically(path *[][]int, u, v, u1, v1 int) bool {
	for i := 0; i < s.numSup; i++ {
		if i != u && s.route[i][v] != 0 {
			if s.findHorizontally(path, i, v, u1, v1) {
				*path = append(*path, []int{i, v})
				return true
			}
		}
	}

	return false
}

func (s *Solution) findOptimal() error {
	path, err := s.findPath(s.pivotN, s.pivotM)
	if err != nil {
		return err
	}

	min := math.MaxFloat64

	for i := 1; i < len(path); i += 2 {
		t := s.route[path[i][0]][path[i][1]]
		if t < min {
			min = t
		}
	}

	for i := 1; i < len(path); i += 2 {
		s.route[path[i][0]][path[i][1]] -= min
		s.route[path[i-1][0]][path[i-1][1]] += min
	}

	return nil
}

func (s *Solution) notOptimal() (bool, error) {
	nMax := -math.MaxFloat64
	err := s.calcPotentials()
	if err != nil {
		return false, err
	}

	for i := 0; i < s.numSup; i++ {
		for j := 0; j < s.numCons; j++ {
			x := s.pots[i][j]
			if x > nMax {
				nMax = x
				s.pivotN = i
				s.pivotM = j
			}
		}
	}

	return nMax > 0, nil
}

func (s *Solution) calcPotentials() error {
	for i := 0; i < s.numSup; i++ {
		for j := 0; j < s.numCons; j++ {
			s.pots[i][j] = -0.5

			if s.route[i][j] == 0 {
				path, err := s.findPath(i, j)
				if err != nil {
					return err
				}

				v := -1.0
				x := 0.0

				for _, node := range path {
					x += v * s.TransportCost[node[0]][node[1]]
					v *= -1
				}

				s.pots[i][j] = x
			}
		}
	}

	return nil
}

func (s *Solution) restoreInitialConditions() {
	for i := range s.Consumers {
		s.Consumers[i].Demand -= elipsis
	}
	s.Suppliers[s.numSup-1].Stock -= elipsis * float64(s.numCons)
}

func (s *Solution) createRoute() {
	for i := 0; i < s.numSup; i++ {
		for j := 0; j < s.numCons; j++ {
			amount := math.Round(s.route[i][j])
			if amount == 0 {
				continue
			}

			sup := s.Suppliers[i]
			if sup.fictive {
				continue
			}

			cons := s.Consumers[j]
			if cons.fictive {
				continue
			}

			s.Cost += math.Round(s.TransportCost[i][j] * amount)

			s.Route = append(s.Route, transportNode{
				Sup: &sup, Cons: &cons, Amount: amount,
			})
		}
	}
}

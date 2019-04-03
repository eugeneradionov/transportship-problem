package models

import (
	"math"
	"reflect"
	"testing"
)

func TestNewSolution(t *testing.T) {
	type args struct {
		cond ProblemConditions
	}
	tests := []struct {
		name string
		args args
		want *Solution
	}{
		{name: "New solution", want: &Solution{pivotM: -1, pivotN: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSolution(tt.args.cond); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSolution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSolution_checkInputValues(t *testing.T) {
	type fields struct {
		ProblemConditions ProblemConditions
	}

	suppliers := []Supplier{
		Supplier{
			ID:    1,
			Stock: 30,
		},
		Supplier{
			ID:    2,
			Stock: 40,
		},
		Supplier{
			ID:    3,
			Stock: 20,
		},
	}

	consumers := []Consumer{
		Consumer{
			ID:     1,
			Demand: 20,
		},
		Consumer{
			ID:     2,
			Demand: 30,
		},
		Consumer{
			ID:     3,
			Demand: 30,
		},
		Consumer{
			ID:     4,
			Demand: 10,
		},
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: suppliers,
				Consumers: consumers,
				TransportCost: [][]float64{
					[]float64{},
				},
			},
		}},
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: suppliers,
				Consumers: consumers,
				TransportCost: [][]float64{
					[]float64{},
					[]float64{},
					[]float64{},
				},
			},
		}},
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: suppliers,
				Consumers: consumers,
				TransportCost: [][]float64{
					[]float64{1, 2, 3, 4},
					[]float64{5, 6, 7, 8},
					[]float64{9, 10, 11, 12},
				},
			},
		}},
	}

	tt := tests[0]
	t.Run(tt.name, func(t *testing.T) {
		s := &Solution{
			ProblemConditions: tt.fields.ProblemConditions,
		}
		got := s.checkInputValues()
		if got == nil {
			t.Errorf("Solution.checkInputValues() = \"%v\" want \"%v\"", got, "error")
		}
	})

	tt = tests[1]
	t.Run(tt.name, func(t *testing.T) {
		s := &Solution{
			ProblemConditions: tt.fields.ProblemConditions,
		}
		got := s.checkInputValues()
		if got == nil {
			t.Errorf("Solution.checkInputValues() = \"%v\" want \"%v\"", got, "error")
		}
	})

	tt = tests[2]
	t.Run(tt.name, func(t *testing.T) {
		s := &Solution{
			ProblemConditions: tt.fields.ProblemConditions,
		}
		got := s.checkInputValues()
		if got != nil {
			t.Errorf("Solution.checkInputValues() = \"%v\" want \"%v\"", got, nil)
		}
	})
}

func TestSolution_fixImbalance(t *testing.T) {
	type fields struct {
		ProblemConditions ProblemConditions
	}

	suppliers := []Supplier{
		Supplier{
			ID:    1,
			Stock: 30,
		},
		Supplier{
			ID:    2,
			Stock: 40,
		},
		Supplier{
			ID:    3,
			Stock: 20,
		},
	}

	consumers := []Consumer{
		Consumer{
			ID:     1,
			Demand: 20,
		},
		Consumer{
			ID:     2,
			Demand: 30,
		},
		Consumer{
			ID:     3,
			Demand: 30,
		},
		Consumer{
			ID:     4,
			Demand: 10,
		},
	}

	tests := []struct {
		name        string
		fields      fields
		wantSupLen  int
		wantConsLen int
	}{
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: suppliers,
				Consumers: consumers,
			},
		},
			wantSupLen:  3,
			wantConsLen: 4,
		},
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: []Supplier{
					Supplier{
						ID:    1,
						Stock: 30,
					},
					Supplier{
						ID:    2,
						Stock: 50,
					},
					Supplier{
						ID:    3,
						Stock: 20,
					},
				},
				Consumers: consumers,
			},
		},
			wantSupLen:  3,
			wantConsLen: 5,
		},
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: suppliers,
				Consumers: []Consumer{
					Consumer{
						ID:     1,
						Demand: 30,
					},
					Consumer{
						ID:     2,
						Demand: 30,
					},
					Consumer{
						ID:     3,
						Demand: 30,
					},
					Consumer{
						ID:     4,
						Demand: 10,
					},
				},
			},
		},
			wantSupLen:  4,
			wantConsLen: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Solution{
				ProblemConditions: tt.fields.ProblemConditions,
			}

			s.fixImbalance()

			supLen := len(s.Suppliers)
			if supLen != tt.wantSupLen {
				t.Errorf("Solution.fixImbalance() invalid suppliers length got: %v; want %v", supLen, tt.wantSupLen)
			}

			consLen := len(s.Consumers)
			if consLen != tt.wantConsLen {
				t.Errorf("Solution.fixImbalance() invalid consumers length got: %v; want %v", consLen, tt.wantConsLen)
			}
		})
	}
}

func TestSolution_Solve(t *testing.T) {
	type fields struct {
		ProblemConditions ProblemConditions
		Route             []transportNode
		Cost              float64
		pivotN            int
		pivotM            int
		route             [][]float64
		pots              [][]float64
		numSup            int
		numCons           int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: []Supplier{
					Supplier{
						ID:    1,
						Stock: 30,
					},
					Supplier{
						ID:    2,
						Stock: 40,
					},
					Supplier{
						ID:    3,
						Stock: 20,
					},
				},
				Consumers: []Consumer{
					Consumer{
						ID:     1,
						Demand: 20,
					},
					Consumer{
						ID:     2,
						Demand: 30,
					},
					Consumer{
						ID:     3,
						Demand: 30,
					},
					Consumer{
						ID:     4,
						Demand: 10,
					},
				},
				TransportCost: [][]float64{
					[]float64{2, 3, 2, 4},
					[]float64{3, 2, 5, 1},
					[]float64{4, 3, 2, 6},
				},
			},
		}},
	}

	route := []transportNode{
		transportNode{
			Sup:    &Supplier{ID: 1},
			Cons:   &Consumer{ID: 1},
			Amount: 20,
		},
		transportNode{
			Sup:    &Supplier{ID: 1},
			Cons:   &Consumer{ID: 3},
			Amount: 10,
		},
		transportNode{
			Sup:    &Supplier{ID: 2},
			Cons:   &Consumer{ID: 2},
			Amount: 30,
		},
		transportNode{
			Sup:    &Supplier{ID: 2},
			Cons:   &Consumer{ID: 4},
			Amount: 10,
		},
		transportNode{
			Sup:    &Supplier{ID: 3},
			Cons:   &Consumer{ID: 3},
			Amount: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Solution{
				ProblemConditions: tt.fields.ProblemConditions,
				Route:             tt.fields.Route,
				Cost:              tt.fields.Cost,
				pivotN:            tt.fields.pivotN,
				pivotM:            tt.fields.pivotM,
				route:             tt.fields.route,
				pots:              tt.fields.pots,
				numSup:            tt.fields.numSup,
				numCons:           tt.fields.numCons,
			}
			err := s.Solve()
			if err != nil {
				t.Error(err)
			}

			if s.Cost != 170 {
				t.Errorf("Solution.Solve() s.Cost = %v want %v", s.Cost, 170)
			}

			for i, node := range s.Route {
				wNode := route[i]

				if node.Amount != wNode.Amount {
					t.Errorf("Solution.Solve() s.Route[%v].Amount = %v want %v", i, node.Amount, wNode.Amount)
				}

				if node.Sup.ID != wNode.Sup.ID {
					t.Errorf("Solution.Solve() s.Route[%v].Sup.ID = %v want %v", i, node.Sup.ID, wNode.Sup.ID)
				}

				if node.Cons.ID != wNode.Cons.ID {
					t.Errorf("Solution.Solve() s.Route[%v].Cons.ID = %v want %v", i, node.Cons.ID, wNode.Cons.ID)
				}
			}
		})
	}
}

func TestSolution_createRoute(t *testing.T) {
	type fields struct {
		ProblemConditions ProblemConditions
		Route             []transportNode
		Cost              float64
		pivotN            int
		pivotM            int
		route             [][]float64
		pots              [][]float64
		numSup            int
		numCons           int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: []Supplier{
					Supplier{
						ID:    1,
						Stock: 30,
					},
					Supplier{
						ID:    2,
						Stock: 40,
					},
					Supplier{
						ID:    3,
						Stock: 20,
					},
				},
				Consumers: []Consumer{
					Consumer{
						ID:     1,
						Demand: 20,
					},
					Consumer{
						ID:     2,
						Demand: 30,
					},
					Consumer{
						ID:     3,
						Demand: 30,
					},
					Consumer{
						ID:     4,
						Demand: 10,
					},
				},
				TransportCost: [][]float64{
					[]float64{2, 3, 2, 4},
					[]float64{3, 2, 5, 1},
					[]float64{4, 3, 2, 6},
				},
			},
		}},
	}

	route := []transportNode{
		transportNode{
			Sup:    &Supplier{ID: 1},
			Cons:   &Consumer{ID: 1},
			Amount: 20,
		},
		transportNode{
			Sup:    &Supplier{ID: 1},
			Cons:   &Consumer{ID: 3},
			Amount: 10,
		},
		transportNode{
			Sup:    &Supplier{ID: 2},
			Cons:   &Consumer{ID: 2},
			Amount: 30,
		},
		transportNode{
			Sup:    &Supplier{ID: 2},
			Cons:   &Consumer{ID: 4},
			Amount: 10,
		},
		transportNode{
			Sup:    &Supplier{ID: 3},
			Cons:   &Consumer{ID: 3},
			Amount: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Solution{
				ProblemConditions: tt.fields.ProblemConditions,
				Route:             tt.fields.Route,
				Cost:              tt.fields.Cost,
				pivotN:            tt.fields.pivotN,
				pivotM:            tt.fields.pivotM,
				route:             tt.fields.route,
				pots:              tt.fields.pots,
				numSup:            tt.fields.numSup,
				numCons:           tt.fields.numCons,
			}
			s.initData()
			s.northWest()

			nOpt, err := s.notOptimal()
			if err != nil {
				t.Error(err)
			}

			for nOpt {
				err := s.findOptimal()
				if err != nil {
					t.Error(err)
				}

				nOpt, err = s.notOptimal()
				if err != nil {
					t.Error(err)
				}
			}
			s.createRoute()

			if s.Cost != 170 {
				t.Errorf("Solution.Solve() s.Cost = %v want %v", s.Cost, 170)
			}

			for i, node := range s.Route {
				wNode := route[i]

				if node.Amount != wNode.Amount {
					t.Errorf("Solution.Solve() s.Route[%v].Amount = %v want %v", i, node.Amount, wNode.Amount)
				}

				if node.Sup.ID != wNode.Sup.ID {
					t.Errorf("Solution.Solve() s.Route[%v].Sup.ID = %v want %v", i, node.Sup.ID, wNode.Sup.ID)
				}

				if node.Cons.ID != wNode.Cons.ID {
					t.Errorf("Solution.Solve() s.Route[%v].Cons.ID = %v want %v", i, node.Cons.ID, wNode.Cons.ID)
				}
			}
		})
	}
}

func TestSolution_restoreInitialConditions(t *testing.T) {
	type fields struct {
		ProblemConditions ProblemConditions
		Route             []transportNode
		Cost              float64
		pivotN            int
		pivotM            int
		route             [][]float64
		pots              [][]float64
		numSup            int
		numCons           int
	}
	tests := []struct {
		name   string
		fields fields
		want   ProblemConditions
	}{
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: []Supplier{
					Supplier{
						ID:    1,
						Stock: 10,
					},
					Supplier{
						ID:    2,
						Stock: 20 + elipsis*2,
					},
				},
				Consumers: []Consumer{
					Consumer{
						ID:     1,
						Demand: 20 + elipsis,
					},
					Consumer{
						ID:     2,
						Demand: 10 + elipsis,
					},
				},
				TransportCost: [][]float64{
					[]float64{1, 2},
					[]float64{2, 1},
				},
			},
			numSup:  2,
			numCons: 2,
		},
			want: ProblemConditions{
				Suppliers: []Supplier{
					Supplier{
						ID:    1,
						Stock: 10,
					},
					Supplier{
						ID:    2,
						Stock: 20,
					},
				},
				Consumers: []Consumer{
					Consumer{
						ID:     1,
						Demand: 20,
					},
					Consumer{
						ID:     2,
						Demand: 10,
					},
				},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Solution{
				ProblemConditions: tt.fields.ProblemConditions,
				Route:             tt.fields.Route,
				Cost:              tt.fields.Cost,
				pivotN:            tt.fields.pivotN,
				pivotM:            tt.fields.pivotM,
				route:             tt.fields.route,
				pots:              tt.fields.pots,
				numSup:            tt.fields.numSup,
				numCons:           tt.fields.numCons,
			}
			s.restoreInitialConditions()

			for i, sup := range tt.fields.ProblemConditions.Suppliers {
				got := sup.Stock
				want := tt.want.Suppliers[i].Stock
				if got != want {
					t.Errorf("solution.restoreInitialConditions() invalid supplier stock: got: %v, want: %v", got, want)
				}
			}

			for i, cons := range tt.fields.ProblemConditions.Consumers {
				got := cons.Demand
				want := tt.want.Consumers[i].Demand
				if got != want {
					t.Errorf("solution.restoreInitialConditions() invalid consumer demand: got: %v, want: %v", got, want)
				}
			}

		})
	}
}

func TestSolution_initData(t *testing.T) {
	type fields struct {
		ProblemConditions ProblemConditions
		Route             []transportNode
		Cost              float64
		pivotN            int
		pivotM            int
		route             [][]float64
		pots              [][]float64
		numSup            int
		numCons           int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: []Supplier{
					Supplier{
						ID:    1,
						Stock: 10,
					},
					Supplier{
						ID:    2,
						Stock: 20,
					},
				},
				Consumers: []Consumer{
					Consumer{
						ID:     1,
						Demand: 20,
					},
					Consumer{
						ID:     2,
						Demand: 10,
					},
				},
				TransportCost: [][]float64{
					[]float64{1, 2},
					[]float64{2, 1},
				},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Solution{
				ProblemConditions: tt.fields.ProblemConditions,
				Route:             tt.fields.Route,
				Cost:              tt.fields.Cost,
				pivotN:            tt.fields.pivotN,
				pivotM:            tt.fields.pivotM,
				route:             tt.fields.route,
				pots:              tt.fields.pots,
				numSup:            tt.fields.numSup,
				numCons:           tt.fields.numCons,
			}
			s.initData()

			if s.numSup != len(tt.fields.ProblemConditions.Suppliers) {
				t.Errorf("Solution.initData() = %v, want %v", s.numSup, len(tt.fields.ProblemConditions.Suppliers))
			}

			if s.numCons != len(tt.fields.ProblemConditions.Consumers) {
				t.Errorf("Solution.initData() = %v, want %v", s.numCons, len(tt.fields.ProblemConditions.Consumers))
			}

			want := 20.00001
			if s.Consumers[0].Demand != want {
				t.Errorf("Solution.initData() s.Consumers[0].Demand = %v, want %v", s.Consumers[0].Demand, want)
			}

			want = 10.00001
			if s.Consumers[1].Demand != want {
				t.Errorf("Solution.initData() s.Consumers[1].Demand = %v, want %v", s.Consumers[1].Demand, want)
			}

		})
	}
}

func TestSolution_northWest(t *testing.T) {
	type fields struct {
		ProblemConditions ProblemConditions
		Route             []transportNode
		Cost              float64
		pivotN            int
		pivotM            int
		route             [][]float64
		pots              [][]float64
		numSup            int
		numCons           int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: []Supplier{
					Supplier{
						ID:    1,
						Stock: 10,
					},
					Supplier{
						ID:    2,
						Stock: 20,
					},
				},
				Consumers: []Consumer{
					Consumer{
						ID:     1,
						Demand: 20,
					},
					Consumer{
						ID:     2,
						Demand: 10,
					},
				},
				TransportCost: [][]float64{
					[]float64{1, 2},
					[]float64{2, 1},
				},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Solution{
				ProblemConditions: tt.fields.ProblemConditions,
				Route:             tt.fields.Route,
				Cost:              tt.fields.Cost,
				pivotN:            tt.fields.pivotN,
				pivotM:            tt.fields.pivotM,
				route:             tt.fields.route,
				pots:              tt.fields.pots,
				numSup:            tt.fields.numSup,
				numCons:           tt.fields.numCons,
			}
			s.initData()
			s.northWest()

			nwStartSolution := [][]float64{
				[]float64{10, 0},
				[]float64{10, 10},
			}

			for i, row := range s.route {
				for j, val := range row {
					val = math.Round(val)
					if val != nwStartSolution[i][j] {
						t.Errorf("Solution.northWest() s.route[%d][%d] = %v, want %v", i, j, val, nwStartSolution[i][j])
					}
				}
			}

		})
	}
}

func TestSolution_notOptimal(t *testing.T) {
	type fields struct {
		ProblemConditions ProblemConditions
		Route             []transportNode
		Cost              float64
		pivotN            int
		pivotM            int
		route             [][]float64
		pots              [][]float64
		numSup            int
		numCons           int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: []Supplier{
					Supplier{
						ID:    1,
						Stock: 10,
					},
					Supplier{
						ID:    2,
						Stock: 20,
					},
				},
				Consumers: []Consumer{
					Consumer{
						ID:     1,
						Demand: 20,
					},
					Consumer{
						ID:     2,
						Demand: 10,
					},
				},
				TransportCost: [][]float64{
					[]float64{1, 2},
					[]float64{2, 1},
				},
			},
		}},
		{fields: fields{
			ProblemConditions: ProblemConditions{
				Suppliers: []Supplier{
					Supplier{
						ID:    1,
						Stock: 30,
					},
					Supplier{
						ID:    2,
						Stock: 40,
					},
					Supplier{
						ID:    3,
						Stock: 20,
					},
				},
				Consumers: []Consumer{
					Consumer{
						ID:     1,
						Demand: 20,
					},
					Consumer{
						ID:     2,
						Demand: 30,
					},
					Consumer{
						ID:     3,
						Demand: 30,
					},
					Consumer{
						ID:     4,
						Demand: 10,
					},
				},
				TransportCost: [][]float64{
					[]float64{2, 3, 2, 4},
					[]float64{3, 2, 5, 1},
					[]float64{4, 3, 2, 6},
				},
			},
		},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Solution{
				ProblemConditions: tt.fields.ProblemConditions,
				Route:             tt.fields.Route,
				Cost:              tt.fields.Cost,
				pivotN:            tt.fields.pivotN,
				pivotM:            tt.fields.pivotM,
				route:             tt.fields.route,
				pots:              tt.fields.pots,
				numSup:            tt.fields.numSup,
				numCons:           tt.fields.numCons,
			}
			s.initData()
			s.northWest()
			got, err := s.notOptimal()
			if err != nil {
				t.Error(err)
			}

			if got != tt.want {
				t.Errorf("Solution.notOptimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

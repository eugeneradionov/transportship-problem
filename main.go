package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/eugeneradionov/transportship-problem/models"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", solve).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func solve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cond models.ProblemConditions

	if err := json.NewDecoder(r.Body).Decode(&cond); err != nil {
		http.Error(w, "invalid input data", http.StatusUnprocessableEntity)
		return
	}
	log.Printf("[INFO] at solve(): Received Body: %+v", cond)

	s := models.NewSolution(cond)
	if err := s.Solve(); err != nil {
		http.Error(w, fmt.Sprintf("something went wrong: %v", err), http.StatusUnprocessableEntity)
		return
	}
	log.Printf("[INFO] at solve(): Got solution: Route:{%v}", s.Route)

	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, "something went wrong", http.StatusUnprocessableEntity)
		return
	}
}

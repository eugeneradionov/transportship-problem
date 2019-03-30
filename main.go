package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eugeneradionov/transportship-problem/httperrors"

	"github.com/eugeneradionov/transportship-problem/models"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", solve).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", (os.Getenv("PORT"))), r))
}

func solve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cond models.ProblemConditions

	if err := json.NewDecoder(r.Body).Decode(&cond); err != nil {
		httperrors.SendErrorJSON(w, err, http.StatusUnprocessableEntity)
		return
	}
	log.Printf("[INFO] at main.solve(): Received Body: %+v", cond)

	s := models.NewSolution(cond)
	if err := s.Solve(); err != nil {
		httperrors.SendErrorJSON(w, err, http.StatusUnprocessableEntity)
		return
	}
	log.Printf("[INFO] at main.solve(): Got solution: Route:{%v}", s.Route)

	if err := json.NewEncoder(w).Encode(s); err != nil {
		httperrors.SendErrorJSON(w, err, http.StatusUnprocessableEntity)
		return
	}
}

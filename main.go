package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/eugeneradionov/transportship-problem/httperrors"

	"github.com/eugeneradionov/transportship-problem/models"
	"github.com/gorilla/mux"
)

const (
	STATIC_DIR = "./public/"
)

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(STATIC_DIR))).Methods(http.MethodGet)
	r.HandleFunc("/solve", solve).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}

func solve(w http.ResponseWriter, r *http.Request) {
	var cond models.ProblemConditions

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewDecoder(r.Body).Decode(&cond); err != nil {
		httperrors.SendErrorJSON(w, err, http.StatusBadRequest)
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

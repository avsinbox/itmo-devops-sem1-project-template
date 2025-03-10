package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"itmo-devops-sem1-project-template/internal/handlers"
)

func CreateRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// POST /api/v0/prices
	router.HandleFunc("/api/v0/prices", handlers.PostPrices(db)).Methods(http.MethodPost)

	// GET /api/v0/prices
	router.HandleFunc("/api/v0/prices", handlers.GetPrices(db)).Methods(http.MethodGet)

	return router
}

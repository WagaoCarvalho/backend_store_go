package routes

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", handlers.GetHome).Methods("GET")

	//protectedRoutes := r.PathPrefix("/").Subrouter()

	return r
}

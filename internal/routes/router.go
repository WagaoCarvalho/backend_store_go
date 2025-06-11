package routes

import (
	"log"
	"net/http"

	homeHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/home"
	"github.com/WagaoCarvalho/backend_store_go/internal/middlewares"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/db_postgres"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(middlewares.Logging)
	r.Use(middlewares.RecoverPanic)
	r.Use(middlewares.RateLimiter)
	r.Use(middlewares.CORS)

	db, err := repo.Connect(&repo.RealPgxPool{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Rotas públicas
	r.HandleFunc("/", homeHandlers.GetHome).Methods(http.MethodGet)

	// Módulos de rota
	RegisterLoginRoutes(r, db)
	RegisterUserRoutes(r, db)
	RegisterUserCategoryRoutes(r, db)
	RegisterProductRoutes(r, db)
	RegisterAddressRoutes(r, db)
	RegisterContactRoutes(r, db)
	RegisterSupplierRoutes(r, db)
	RegisterSupplierCategoryRoutes(r, db)

	return r
}

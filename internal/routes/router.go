package routes

import (
	"net/http"

	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/home"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/middlewares"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/db_postgres"
	"github.com/gorilla/mux"
)

func NewRouter(log *logger.LoggerAdapter) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(middlewares.RequestIDMiddleware())
	r.Use(middlewares.RecoverMiddleware(log))
	r.Use(middlewares.LoggingMiddleware(log))
	r.Use(middlewares.RateLimiter)
	r.Use(middlewares.CORS)

	db, err := repo.Connect(&repo.RealPgxPool{})
	if err != nil {
		log.Error(nil, err, "Erro ao conectar ao banco de dados", nil)
	}

	r.HandleFunc("/", handlers.GetHome).Methods(http.MethodGet)

	RegisterLoginRoutes(r, db)
	RegisterUserRoutes(r, db)
	RegisterUserCategoryRoutes(r, db)
	RegisterUserCategoryRelationRoutes(r, db)
	RegisterProductRoutes(r, db)
	RegisterAddressRoutes(r, db)
	RegisterContactRoutes(r, db)
	RegisterSupplierRoutes(r, db)
	RegisterSupplierCategoryRoutes(r, db)

	return r
}

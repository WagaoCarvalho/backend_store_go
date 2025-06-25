package routes

import (
	"net/http"

	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/home"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	cors "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/cors"
	logging "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/logging"
	rate_limiter "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/rate_limiter"
	recover "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/recover"
	request "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/request"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/db_postgres"
	"github.com/gorilla/mux"
)

func NewRouter(log *logger.LoggerAdapter) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(request.RequestIDMiddleware())
	r.Use(recover.RecoverMiddleware(log))
	r.Use(logging.LoggingMiddleware(log))
	r.Use(rate_limiter.RateLimiter)
	r.Use(cors.CORS)

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

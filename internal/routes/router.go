package routes

import (
	"context"
	"net/http"

	redis "github.com/WagaoCarvalho/backend_store_go/internal/auth/blacklist_redis"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handler/home"
	cors "github.com/WagaoCarvalho/backend_store_go/internal/middleware/cors"
	logging "github.com/WagaoCarvalho/backend_store_go/internal/middleware/logging"
	rate_limiter "github.com/WagaoCarvalho/backend_store_go/internal/middleware/rate_limiter"
	recover "github.com/WagaoCarvalho/backend_store_go/internal/middleware/recover"
	request "github.com/WagaoCarvalho/backend_store_go/internal/middleware/request"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/db_postgres"
	routes_product "github.com/WagaoCarvalho/backend_store_go/internal/routes/products"
	routes_supplier "github.com/WagaoCarvalho/backend_store_go/internal/routes/suppliers"
	routes_user "github.com/WagaoCarvalho/backend_store_go/internal/routes/users"
	"github.com/WagaoCarvalho/backend_store_go/logger"
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
		log.Error(context.TODO(), err, "Erro ao conectar ao banco de dados", nil)
	}

	blacklist := redis.NewRedisTokenBlacklist("localhost:6379", "", 0)

	r.HandleFunc("/", handlers.GetHome).Methods(http.MethodGet)

	RegisterLoginRoutes(r, db, log, blacklist)

	//Suupliers
	routes_supplier.RegisterSupplierRoutes(r, db, log, blacklist)
	routes_supplier.RegisterSupplierFullRoutes(r, db, log, blacklist)
	routes_supplier.RegisterSupplierCategoryRoutes(r, db, log, blacklist)
	routes_supplier.RegisterSupplierCategoryRelationRoutes(r, db, log, blacklist)

	//Users
	routes_user.RegisterUserRoutes(r, db, log, blacklist)
	routes_user.RegisterUserFullRoutes(r, db, log, blacklist)
	routes_user.RegisterUserCategoryRoutes(r, db, log, blacklist)
	routes_user.RegisterUserCategoryRelationRoutes(r, db, log, blacklist)

	//Products
	routes_product.RegisterProductRoutes(r, db, log, blacklist)

	//Adressess
	RegisterAddressRoutes(r, db, log, blacklist)

	//Contacts
	RegisterContactRoutes(r, db, log, blacklist)

	return r
}

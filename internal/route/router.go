package routes

import (
	"context"
	"net/http"

	redis "github.com/WagaoCarvalho/backend_store_go/infra/db/redis"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handler/home"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	cors "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/cors"
	logging "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/logging"
	rateLimiter "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/rate_limiter"
	recover "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/recover"
	request "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/request"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db_postgres"
	routesClient "github.com/WagaoCarvalho/backend_store_go/internal/route/client"
	routesProduct "github.com/WagaoCarvalho/backend_store_go/internal/route/product"
	routesSale "github.com/WagaoCarvalho/backend_store_go/internal/route/sale"
	routesSupplier "github.com/WagaoCarvalho/backend_store_go/internal/route/supplier"
	routesUser "github.com/WagaoCarvalho/backend_store_go/internal/route/user"
	"github.com/gorilla/mux"
)

func NewRouter(log *logger.LogAdapter) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(request.RequestIDMiddleware())
	r.Use(recover.RecoverMiddleware(log))
	r.Use(logging.LoggingMiddleware(log))
	r.Use(rateLimiter.RateLimiter)
	r.Use(cors.CORS)

	db, err := repo.Connect(&repo.RealPgxPool{})
	if err != nil {
		log.Error(context.TODO(), err, "Erro ao conectar ao banco de dados", nil)
	}

	blacklist := redis.NewRedisTokenBlacklist("localhost:6379", "", 0)

	r.HandleFunc("/", handlers.GetHome).Methods(http.MethodGet)

	RegisterLoginRoutes(r, db, log, blacklist)

	//Suupliers
	routesSupplier.RegisterSupplierRoutes(r, db, log, blacklist)
	routesSupplier.RegisterSupplierFullRoutes(r, db, log, blacklist)
	routesSupplier.RegisterSupplierCategoryRoutes(r, db, log, blacklist)
	routesSupplier.RegisterSupplierCategoryRelationRoutes(r, db, log, blacklist)

	//Users
	routesUser.RegisterUserRoutes(r, db, log, blacklist)
	routesUser.RegisterUserFullRoutes(r, db, log, blacklist)
	routesUser.RegisterUserCategoryRoutes(r, db, log, blacklist)
	routesUser.RegisterUserCategoryRelationRoutes(r, db, log, blacklist)

	//Clients
	routesClient.RegisterClientRoutes(r, db, log, blacklist)

	//Products
	routesProduct.RegisterProductRoutes(r, db, log, blacklist)

	//Sale
	routesSale.RegisterSaleRoutes(r, db, log, blacklist)

	//Adressess
	RegisterAddressRoutes(r, db, log, blacklist)

	//Contacts
	RegisterContactRoutes(r, db, log, blacklist)

	return r
}

package routes

import (
	"net/http"

	jwt_auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handler/product"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middleware/jwt"
	repo_product "github.com/WagaoCarvalho/backend_store_go/internal/repositories/products"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/products"
	"github.com/WagaoCarvalho/backend_store_go/logger"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterProductRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repo_product := repo_product.NewProductRepository(db, log)

	productService := services.NewProductService(repo_product, log)
	handler := handlers.NewProductHandler(productService, log)

	// Config JWT
	jwtCfg := config.LoadJwtConfig()
	jwtManager := jwt_auth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	// Rotas protegidas
	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, jwtManager))

	s.HandleFunc("/product", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/products", handler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/product/{id:[0-9]+}", handler.GetById).Methods(http.MethodGet)
	s.HandleFunc("/product/name/{name}", handler.GetByName).Methods(http.MethodGet)
	s.HandleFunc("/product/manufacturer/{manufacturer}", handler.GetByManufacturer).Methods(http.MethodGet)
	s.HandleFunc("/product/version/{id:[0-9]+}", handler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/product/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/product/enable/{id:[0-9]+}", handler.Enable).Methods(http.MethodPatch)
	s.HandleFunc("/product/disable/{id:[0-9]+}", handler.Disable).Methods(http.MethodPatch)
	s.HandleFunc("/product/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}

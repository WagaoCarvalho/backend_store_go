package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/product"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product"
	jwt_auth "github.com/WagaoCarvalho/backend_store_go/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/pkg/middleware/jwt"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterProductRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repo_product := repo.NewProductRepository(db, log)

	productService := service.NewProductService(repo_product, log)
	handler := handler.NewProductHandler(productService, log)

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

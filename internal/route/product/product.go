package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/handler/product/filter"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/product/product"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repoFilter "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/filter"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/product"
	serviceFilter "github.com/WagaoCarvalho/backend_store_go/internal/service/product/filter"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product/product"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterProductRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repo := repo.NewProduct(db)
	productService := service.NewProductService(repo)
	handler := handler.NewProductHandler(productService, log)

	repoFilter := repoFilter.NewFilterProduct(db)
	serviceFilter := serviceFilter.NewProductFilterService(repoFilter)
	filter := filter.NewProductFilterHandler(serviceFilter, log)

	// Config JWT
	jwtCfg := config.LoadJwtConfig()
	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	// Rotas protegidas
	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, jwtManager))

	s.HandleFunc("/product", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/product/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/product/update/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/product/delete/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)

	s.HandleFunc("/product/enable/{id:[0-9]+}", handler.EnableProduct).Methods(http.MethodPatch)
	s.HandleFunc("/product/disable/{id:[0-9]+}", handler.DisableProduct).Methods(http.MethodPatch)

	s.HandleFunc("/product/version/{id:[0-9]+}", handler.GetVersionByID).Methods(http.MethodGet)

	s.HandleFunc("/product/stock/{id:[0-9]+}", handler.UpdateStock).Methods(http.MethodPatch)
	s.HandleFunc("/product/increase-stock/{id:[0-9]+}", handler.IncreaseStock).Methods(http.MethodPatch)
	s.HandleFunc("/product/decrease-stock/{id:[0-9]+}", handler.DecreaseStock).Methods(http.MethodPatch)
	s.HandleFunc("/product/get-stock/{id:[0-9]+}", handler.GetStock).Methods(http.MethodGet)

	s.HandleFunc("/product/enable-discount/{id:[0-9]+}", handler.EnableDiscount).Methods(http.MethodPatch)
	s.HandleFunc("/product/disable-discount/{id:[0-9]+}", handler.DisableDiscount).Methods(http.MethodPatch)
	s.HandleFunc("/product/apply-discount/{id:[0-9]+}", handler.ApplyDiscount).Methods(http.MethodPatch)

	s.HandleFunc("/products/filter", filter.Filter).Methods(http.MethodGet)

}

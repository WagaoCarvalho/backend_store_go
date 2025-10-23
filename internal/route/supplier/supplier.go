package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/supplier/supplier"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSupplierRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repoSupplier := repo.NewSupplier(db)
	supplierService := service.NewSupplier(repoSupplier)
	handler := handler.NewSupplierHandler(supplierService, log)

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

	s.HandleFunc("/supplier", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/suppliers", handler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/supplier/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/supplier/name/{name}", handler.GetByName).Methods(http.MethodGet)
	s.HandleFunc("/supplier/version/{id:[0-9]+}", handler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/supplier/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/supplier/enable/{id:[0-9]+}", handler.Enable).Methods(http.MethodPatch)
	s.HandleFunc("/supplier/disable/{id:[0-9]+}", handler.Disable).Methods(http.MethodPatch)
	s.HandleFunc("/supplier/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}

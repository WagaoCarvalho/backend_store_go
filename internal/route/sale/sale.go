package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/handler/sale/filter"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/sale/sale"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repoFilter "github.com/WagaoCarvalho/backend_store_go/internal/repo/sale/filter"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/sale/sale"
	serviceFilter "github.com/WagaoCarvalho/backend_store_go/internal/service/sale/filter"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/sale/sale"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSaleRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repoSale := repo.NewSale(db)
	saleService := service.NewSaleService(repoSale)
	handler := handler.NewSaleHandler(saleService, log)

	repoFilter := repoFilter.NewFilterSale(db)
	serviceFilter := serviceFilter.NewSaleFilterService(repoFilter)
	filter := filter.NewSaleFilterHandler(serviceFilter, log)

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

	s.HandleFunc("/sale", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/sale/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/sale/client/{client_id:[0-9]+}", handler.GetByClientID).Methods(http.MethodGet)
	s.HandleFunc("/sale/user/{user_id:[0-9]+}", handler.GetByUserID).Methods(http.MethodGet)
	s.HandleFunc("/sale/status/{status}", handler.GetByStatus).Methods(http.MethodGet)
	s.HandleFunc("/sale/version/{id:[0-9]+}", handler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/sale/date-range/{start}/{end}", handler.GetByDateRange).Methods(http.MethodGet)
	s.HandleFunc("/sale/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/sale/delete/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
	s.HandleFunc("/sale/{id:[0-9]+}/activate", handler.Activate).Methods(http.MethodPatch)
	s.HandleFunc("/sale/{id:[0-9]+}/cancel", handler.Cancel).Methods(http.MethodPatch)
	s.HandleFunc("/sale/{id:[0-9]+}/complete", handler.Complete).Methods(http.MethodPatch)
	s.HandleFunc("/sale/{id:[0-9]+}/returned", handler.Returned).Methods(http.MethodPatch)

	s.HandleFunc("/sales/filter", filter.Filter).Methods(http.MethodGet)
}

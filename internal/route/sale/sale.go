package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/sale/sale"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/sale"
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
	saleService := service.NewSale(repoSale)
	handler := handler.NewSaleHandler(saleService, log)

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
	s.HandleFunc("/sale/date-range/{start}/{end}", handler.GetByDateRange).Methods(http.MethodGet)
	s.HandleFunc("/sale/update/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/sale/delete/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}

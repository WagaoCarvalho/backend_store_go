package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handler/supplier/full"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoContact "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
	repoCatRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/category_relation"
	repoContactRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/contact_relation"
	repoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/full"
	services "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/full"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSupplierFullRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repoSupplier := repoSupplier.NewSupplierFull(db)
	repoAddress := repoAddress.NewAddressTx()
	repoContact := repoContact.NewContactTx(db)
	repoCatRel := repoCatRel.NewSupplierCategoryRelationTx(db)
	repoContactRel := repoContactRel.NewSupplierContactRelation(db)

	supplierService := services.NewSupplierFullService(repoSupplier, repoAddress, repoContact, repoCatRel, repoContactRel)
	handler := handlers.NewSupplierFull(supplierService, log)

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

	s.HandleFunc("/supplier-full", handler.CreateFull).Methods(http.MethodPost)
}

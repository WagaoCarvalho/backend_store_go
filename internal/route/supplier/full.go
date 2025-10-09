package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handler/supplier/supplier_full"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoContact "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
	repoCatRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_category_relations"
	repoContactRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_contact_relations"
	repoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_full_repositories"
	services "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier_full_services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSupplierFullRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repoSupplier := repoSupplier.NewSupplierFullRepository(db)
	repoAddress := repoAddress.NewAddressRepository(db)
	repoContact := repoContact.NewContactRepository(db)
	repoCatRel := repoCatRel.NewSupplierCategoryRelationRepo(db)
	repoContactRel := repoContactRel.NewSupplierContactRelationRepositories(db)

	supplierService := services.NewSupplierFullService(repoSupplier, repoAddress, repoContact, repoCatRel, repoContactRel)
	handler := handlers.NewSupplierFullHandler(supplierService, log)

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

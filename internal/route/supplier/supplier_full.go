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
	repoSupplierCatRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_category_relations"
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
	repoUser := repoSupplier.NewSupplierFullRepository(db)
	repoAddress := repoAddress.NewAddressRepository(db)
	repoContact := repoContact.NewContactRepository(db)
	repoUserCatRel := repoSupplierCatRel.NewSupplierCategoryRelationRepo(db)

	userService := services.NewSupplierFullService(repoUser, repoAddress, repoContact, repoUserCatRel)
	handler := handlers.NewSupplierFullHandler(userService, log)

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

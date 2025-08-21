package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handler/supplier/supplier_full"
	jwt_auth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo_address "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repo_contact "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
	repo_supplier_cat_rel "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_category_relations"
	repo_supplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_full_repositories"
	services "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier_full_services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSupplierFullRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repo_user := repo_supplier.NewSupplierFullRepository(db, log)
	repo_address := repo_address.NewAddressRepository(db, log)
	repo_contact := repo_contact.NewContactRepository(db, log)
	repo_user_cat_rel := repo_supplier_cat_rel.NewSupplierCategoryRelationRepo(db, log)

	userService := services.NewSupplierFullService(repo_user, repo_address, repo_contact, repo_user_cat_rel, log)
	handler := handlers.NewSupplierFullHandler(userService, log)

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

	s.HandleFunc("/supplier-full", handler.CreateFull).Methods(http.MethodPost)
}

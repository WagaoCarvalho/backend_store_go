package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handler/user/user_full"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	auth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoContact "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
	repoUserCatRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_category_relation"
	repoUserContactRel "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_contact_relation"
	repoUser "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_full"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user_full"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserFullRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repoUser := repoUser.NewUser(db)
	repoAddress := repoAddress.NewAddress(db)
	repoContact := repoContact.NewContact(db)
	repoUserCatRel := repoUserCatRel.NewUserCategoryRelation(db)
	repoUserContactRel := repoUserContactRel.NewUserContactRelation(db)
	hasher := auth.BcryptHasher{}

	userService := service.NewUserFull(repoUser, repoAddress, repoContact, repoUserCatRel, repoUserContactRel, hasher)
	handler := handlers.NewUserFullHandler(userService, log)

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
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, jwtManager)) // <- passa jwtManager, não string SecretKey

	s.HandleFunc("/user-full", handler.CreateFull).Methods(http.MethodPost)
}

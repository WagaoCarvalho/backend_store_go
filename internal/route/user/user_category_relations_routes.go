package routes

import (
	"net/http"

	jwt_auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/user/user_category_relations"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_category_relations"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user_category_relations"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserCategoryRelationRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist,
) {
	relationRepo := repo.NewUserCategoryRelationRepositories(db, log)
	relationService := service.NewUserCategoryRelationServices(relationRepo, log)
	relationHandler := handler.NewUserCategoryRelationHandler(relationService, log)

	// Carregar config JWT
	jwtCfg := config.LoadJwtConfig()

	// Criar jwtManager que implementa JWTService
	jwtManager := jwt_auth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, jwtManager))

	s.HandleFunc("/user-category-relation", relationHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/user-category-relations/{user_id:[0-9]+}", relationHandler.GetAllRelationsByUserID).Methods(http.MethodGet)
	s.HandleFunc("/user-category-relation/{user_id:[0-9]+}/category/{category_id:[0-9]+}/exists", relationHandler.HasUserCategoryRelation).Methods(http.MethodGet)
	s.HandleFunc("/user-category-relation/{user_id:[0-9]+}/category/{category_id:[0-9]+}", relationHandler.Delete).Methods(http.MethodDelete)
	s.HandleFunc("/user-category-relation/{user_id:[0-9]+}", relationHandler.DeleteAll).Methods(http.MethodDelete)
}

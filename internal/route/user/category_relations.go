package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/user/category_relation"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/category_relation"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/category_relation"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserCategoryRelationRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	relationRepo := repo.NewUserCategoryRelation(db)
	relationService := service.NewUserCategoryRelationService(relationRepo)
	relationHandler := handler.NewUserCategoryRelation(relationService, log)

	// Carregar config JWT
	jwtCfg := config.LoadJwtConfig()

	// Criar jwtManager que implementa JWTService
	jwtManager := jwtAuth.NewJWTManager(
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

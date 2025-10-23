package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/user/user_contact_relation"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_contact_relation"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user_contact_relation"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserContactRelationRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	relationRepo := repo.NewUserContactRelation(db)
	relationService := service.NewUserContactRelation(relationRepo)
	relationHandler := handler.NewUserContactRelation(relationService, log)

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

	s.HandleFunc("/user-contact-relation", relationHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/user-contact-relations/{user_id:[0-9]+}", relationHandler.GetAllByUserID).Methods(http.MethodGet)
	s.HandleFunc("/user-contact-relation/{user_id:[0-9]+}/contact/{contact_id:[0-9]+}/exists", relationHandler.HasRelation).Methods(http.MethodGet)
	s.HandleFunc("/user-contact-relation/{user_id:[0-9]+}/contact/{contact_id:[0-9]+}", relationHandler.Delete).Methods(http.MethodDelete)
	s.HandleFunc("/user-contact-relation/{user_id:[0-9]+}", relationHandler.DeleteAll).Methods(http.MethodDelete)
}

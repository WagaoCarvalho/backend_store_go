package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user/user_category_relations_handler"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/jwt"
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	service "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_category_relations"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserCategoryRelationRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist, // <- injeta blacklist aqui
) {
	relationRepo := repository.NewUserCategoryRelationRepositories(db, log)
	relationService := service.NewUserCategoryRelationServices(relationRepo, log)
	relationHandler := handler.NewUserCategoryRelationHandler(relationService, log)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, config.LoadConfig().Jwt.SecretKey)) // <- middleware atualizado com blacklist

	s.HandleFunc("/user-category-relation", relationHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/user-category-relations/{user_id:[0-9]+}", relationHandler.GetAllRelationsByUserID).Methods(http.MethodGet)
	s.HandleFunc("/user-category-relation/{user_id:[0-9]+}/category/{category_id:[0-9]+}/exists", relationHandler.HasUserCategoryRelation).Methods(http.MethodGet)
	s.HandleFunc("/user-category-relation/{user_id:[0-9]+}/category/{category_id:[0-9]+}", relationHandler.Delete).Methods(http.MethodDelete)
	s.HandleFunc("/user-category-relation/{user_id:[0-9]+}", relationHandler.DeleteAll).Methods(http.MethodDelete)
}

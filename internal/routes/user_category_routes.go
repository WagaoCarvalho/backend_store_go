package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user/user_categories_handler"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/jwt"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_categories"
	service "github.com/WagaoCarvalho/backend_store_go/internal/services/users/user_categories"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserCategoryRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist, // <- injeta blacklist
) {
	userCategoryRepo := repositories.NewUserCategoryRepository(db, log)
	userCategoryService := service.NewUserCategoryService(userCategoryRepo, log)
	userCategoryHandler := handler.NewUserCategoryHandler(userCategoryService, log)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, config.LoadConfig().Jwt.SecretKey)) // <- aplica middleware com logger e blacklist

	s.HandleFunc("/user-category", userCategoryHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/user-category/{id:[0-9]+}", userCategoryHandler.GetById).Methods(http.MethodGet)
	s.HandleFunc("/user-categories", userCategoryHandler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/user-category/{id:[0-9]+}", userCategoryHandler.Update).Methods(http.MethodPut)
	s.HandleFunc("/user-category/{id:[0-9]+}", userCategoryHandler.Delete).Methods(http.MethodDelete)
}

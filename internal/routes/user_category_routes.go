package routes

import (
	"net/http"

	userCategoryHandler "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user/user_categories_handler"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/jwt"
	userCategoryRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_categories"
	userCategoryService "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_categories"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserCategoryRoutes(r *mux.Router, db *pgxpool.Pool, log *logger.LoggerAdapter) {
	userCategoryRepo := userCategoryRepositories.NewUserCategoryRepository(db, log)
	userCategoryService := userCategoryService.NewUserCategoryService(userCategoryRepo)
	userCategoryHandler := userCategoryHandler.NewUserCategoryHandler(userCategoryService)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken)

	s.HandleFunc("/user-category", userCategoryHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/user-category/{id:[0-9]+}", userCategoryHandler.GetById).Methods(http.MethodGet)
	s.HandleFunc("/user-categories", userCategoryHandler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/user-category/{id:[0-9]+}", userCategoryHandler.Update).Methods(http.MethodPut)
	s.HandleFunc("/user-category/{id:[0-9]+}", userCategoryHandler.Delete).Methods(http.MethodDelete)
}

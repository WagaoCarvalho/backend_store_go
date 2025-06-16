package routes

import (
	"net/http"

	userCategoryHandler "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user/user_categories_handler"
	"github.com/WagaoCarvalho/backend_store_go/internal/middlewares"
	userCategoryRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_categories"
	userCategoryService "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_categories"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserCategoryRoutes(r *mux.Router, db *pgxpool.Pool) {
	userCategoryRepo := userCategoryRepositories.NewUserCategoryRepository(db)
	userCategoryService := userCategoryService.NewUserCategoryService(userCategoryRepo)
	userCategoryHandler := userCategoryHandler.NewUserCategoryHandler(userCategoryService)

	s := r.PathPrefix("/").Subrouter()
	s.Use(middlewares.IsAuthByBearerToken)

	s.HandleFunc("/user-category", userCategoryHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/user-category/{id:[0-9]+}", userCategoryHandler.GetById).Methods(http.MethodGet)
	s.HandleFunc("/user-category/{id:[0-9]+}", userCategoryHandler.GetVersionByID).Methods(http.MethodPatch)
	s.HandleFunc("/user-category-relation/version/{user_id:[0-9]+}", userCategoryHandler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/user-categories", userCategoryHandler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/user-category/{id:[0-9]+}", userCategoryHandler.Update).Methods(http.MethodPut)
	s.HandleFunc("/user-category/{id:[0-9]+}", userCategoryHandler.Delete).Methods(http.MethodDelete)
}

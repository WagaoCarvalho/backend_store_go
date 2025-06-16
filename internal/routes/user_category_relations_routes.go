package routes

import (
	"net/http"

	userCategoryRelationHandler "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user/user_category_relations_handler"
	"github.com/WagaoCarvalho/backend_store_go/internal/middlewares"
	userCategoryRelationRepository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	userCategoryRelationService "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_category_relations"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserCategoryRelationRoutes(r *mux.Router, db *pgxpool.Pool) {
	relationRepo := userCategoryRelationRepository.NewUserCategoryRelationRepositories(db)
	relationService := userCategoryRelationService.NewUserCategoryRelationServices(relationRepo)
	relationHandler := userCategoryRelationHandler.NewUserCategoryRelationHandler(relationService)

	s := r.PathPrefix("/").Subrouter()
	s.Use(middlewares.IsAuthByBearerToken)

	s.HandleFunc("/user-category-relation", relationHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/user-category-relations/{user_id:[0-9]+}", relationHandler.GetAllRelationsByUserID).Methods(http.MethodGet)
	s.HandleFunc("/user-category-relation/{user_id:[0-9]+}/category/{category_id:[0-9]+}", relationHandler.Delete).Methods(http.MethodDelete)
	s.HandleFunc("/user-category-relation/{user_id:[0-9]+}", relationHandler.DeleteAll).Methods(http.MethodDelete)
}

package routes

import (
	"net/http"

	loginHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/login"
	userRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	loginServices "github.com/WagaoCarvalho/backend_store_go/internal/services/login"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterLoginRoutes(r *mux.Router, db *pgxpool.Pool) {
	userRepo := userRepositories.NewUserRepository(db)
	service := loginServices.NewLoginService(userRepo)
	handler := loginHandlers.NewLoginHandler(service)

	r.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
}

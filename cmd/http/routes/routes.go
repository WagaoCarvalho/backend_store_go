package routes

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/handlers"
	"github.com/WagaoCarvalho/backend_store_go/internal/repositories"
	"github.com/WagaoCarvalho/backend_store_go/internal/services"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	db := repositories.Connect()

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	r.HandleFunc("/user/id/{id}", userHandler.GetUserById).Methods("GET")
	r.HandleFunc("/user/email/{email}", userHandler.GetUserByEmail).Methods("GET")
	r.HandleFunc("/user", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/user/{id}", userHandler.UpdateUser).Methods("PUT")

	return r
}

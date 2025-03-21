package routes

import (
	"log"

	homeHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/home"
	loginHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/login"
	userHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user"
	"github.com/WagaoCarvalho/backend_store_go/internal/middlewares"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/db_postgres"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/user"
	loginServices "github.com/WagaoCarvalho/backend_store_go/internal/services/login"
	userServices "github.com/WagaoCarvalho/backend_store_go/internal/services/user"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// Cria uma instância de PgxPool (usando a implementação real)
	pgx := &repo.RealPgxPool{}

	// Conecta ao banco de dados e trata o erro
	db, err := repo.Connect(pgx)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	r.Use(middlewares.Logging)
	r.Use(middlewares.RecoverPanic)
	r.Use(middlewares.RateLimiter)
	r.Use(middlewares.CORS)

	userRepo := repositories.NewUserRepository(db)
	userService := userServices.NewUserService(userRepo)
	loginService := loginServices.NewLoginService(userRepo)
	userHandler := userHandlers.NewUserHandler(userService)
	loginHandler := loginHandlers.NewLoginHandler(loginService)

	r.HandleFunc("/", homeHandlers.GetHome).Methods("GET")
	r.HandleFunc("/user", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/login", loginHandler.Login).Methods("POST")

	protectedRoutes := r.PathPrefix("/").Subrouter()
	protectedRoutes.Use(middlewares.IsAuthByBearerToken)

	protectedRoutes.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	protectedRoutes.HandleFunc("/user/id/{id}", userHandler.GetUserById).Methods("GET")
	protectedRoutes.HandleFunc("/user/email/{email}", userHandler.GetUserByEmail).Methods("GET")
	protectedRoutes.HandleFunc("/user/{id}", userHandler.UpdateUser).Methods("PUT")
	protectedRoutes.HandleFunc("/user/{id}", userHandler.DeleteUserById).Methods("DELETE")

	return r
}

package routes

import (
	"log"
	"net/http"

	homeHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/home"
	loginHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/login"
	productHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/product"
	userHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user"
	userCategoryHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user_categories"
	"github.com/WagaoCarvalho/backend_store_go/internal/middlewares"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/db_postgres"
	productRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/products"
	userRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/user"
	userCategoryRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/user_categories"
	loginServices "github.com/WagaoCarvalho/backend_store_go/internal/services/login"
	productServices "github.com/WagaoCarvalho/backend_store_go/internal/services/products"
	userServices "github.com/WagaoCarvalho/backend_store_go/internal/services/user"
	userCategoryServices "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_categories"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	db, err := repo.Connect(&repo.RealPgxPool{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	r.Use(middlewares.Logging)
	r.Use(middlewares.RecoverPanic)
	r.Use(middlewares.RateLimiter)
	r.Use(middlewares.CORS)

	// User repository, service and handler setup
	userRepo := userRepositories.NewUserRepository(db)
	userService := userServices.NewUserService(userRepo)
	loginService := loginServices.NewLoginService(userRepo)
	userHandler := userHandlers.NewUserHandler(userService)
	loginHandler := loginHandlers.NewLoginHandler(loginService)

	// UserCategory repository, service and handler setup
	userCategoryRepo := userCategoryRepositories.NewUserCategoryRepository(db)
	userCategoryService := userCategoryServices.NewUserCategoryService(userCategoryRepo)
	userCategoryHandler := userCategoryHandlers.NewUserCategoryHandler(userCategoryService)

	// Home, login and user routes
	r.HandleFunc("/", homeHandlers.GetHome).Methods(http.MethodGet)
	r.HandleFunc("/user", userHandler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/login", loginHandler.Login).Methods(http.MethodPost)

	protectedRoutes := r.PathPrefix("/").Subrouter()
	protectedRoutes.Use(middlewares.IsAuthByBearerToken)

	// User routes
	protectedRoutes.HandleFunc("/users", userHandler.GetUsers).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/user/id/{id}", userHandler.GetUserById).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/user/email/{email}", userHandler.GetUserByEmail).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/user/{id}", userHandler.UpdateUser).Methods(http.MethodPut)
	protectedRoutes.HandleFunc("/user/{id}", userHandler.DeleteUserById).Methods(http.MethodDelete)

	// Product routes
	productRepo := productRepositories.NewProductRepository(db)
	productService := productServices.NewProductService(productRepo)
	productHandler := productHandlers.NewProductHandler(productService)

	protectedRoutes.HandleFunc("/products", productHandler.GetProducts).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/product/{id}", productHandler.GetProductById).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/product", productHandler.CreateProduct).Methods(http.MethodPost)
	protectedRoutes.HandleFunc("/product/{id}", productHandler.UpdateProduct).Methods(http.MethodPut)
	protectedRoutes.HandleFunc("/product/{id}", productHandler.DeleteProductById).Methods(http.MethodDelete)

	// UserCategory routes
	protectedRoutes.HandleFunc("/user-categories", userCategoryHandler.GetCategories).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/user-category/{id}", userCategoryHandler.GetCategoryById).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/user-category", userCategoryHandler.CreateCategory).Methods(http.MethodPost)
	protectedRoutes.HandleFunc("/user-category/{id}", userCategoryHandler.UpdateCategory).Methods(http.MethodPut)
	protectedRoutes.HandleFunc("/user-category/{id}", userCategoryHandler.DeleteCategoryById).Methods(http.MethodDelete)

	// Additional product search routes
	protectedRoutes.HandleFunc("/products/search", productHandler.GetProductsByName).Methods(http.MethodGet)
	//protectedRoutes.HandleFunc("/products/price", productHandler.GetProductsByPriceRange).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/products/low-stock", productHandler.GetProductsLowInStock).Methods(http.MethodGet)

	return r
}

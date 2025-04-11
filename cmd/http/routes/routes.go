package routes

import (
	"log"
	"net/http"

	addressHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/address"
	contactHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/contacts"
	homeHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/home"
	loginHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/login"
	productHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/product"
	supplierHandler "github.com/WagaoCarvalho/backend_store_go/internal/handlers/supplier"
	userHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user"
	userCategoryHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user/user_categories"
	"github.com/WagaoCarvalho/backend_store_go/internal/middlewares"
	addressRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	contactRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/db_postgres"
	productRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/products"
	supplierRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers"
	userRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	userCategoryRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_categories"
	addressServices "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	contactServices "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"
	loginServices "github.com/WagaoCarvalho/backend_store_go/internal/services/login"
	productServices "github.com/WagaoCarvalho/backend_store_go/internal/services/products"
	supplierServices "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers"
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

	// Address repository, service and handler setup
	addressRepo := addressRepositories.NewAddressRepository(db)
	addressService := addressServices.NewAddressService(addressRepo)
	addressHandler := addressHandlers.NewAddressHandler(addressService)

	// Contact repository, service and handler setup
	contactRepo := contactRepositories.NewContactRepository(db)
	contactService := contactServices.NewContactService(contactRepo)
	contactHandler := contactHandlers.NewContactHandler(contactService)

	// Supplier repository, service and handler setup
	supplierRepo := supplierRepositories.NewSupplierRepository(db)
	supplierService := supplierServices.NewSupplierService(supplierRepo)
	supplierHandler := supplierHandler.NewSupplierHandler(supplierService)

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
	protectedRoutes.HandleFunc("/products/search", productHandler.GetProductsByName).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/products/price", productHandler.GetProductsBySalePriceRange).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/products/low-stock", productHandler.GetProductsLowInStock).Methods(http.MethodGet)

	// UserCategory routes
	protectedRoutes.HandleFunc("/user-categories", userCategoryHandler.GetCategories).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/user-category/{id}", userCategoryHandler.GetCategoryById).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/user-category", userCategoryHandler.CreateCategory).Methods(http.MethodPost)
	protectedRoutes.HandleFunc("/user-category/{id}", userCategoryHandler.UpdateCategory).Methods(http.MethodPut)
	protectedRoutes.HandleFunc("/user-category/{id}", userCategoryHandler.DeleteCategoryById).Methods(http.MethodDelete)

	// Address routes
	protectedRoutes.HandleFunc("/addresses", addressHandler.Create).Methods(http.MethodPost)
	//protectedRoutes.HandleFunc("/addresses", addressHandler.GetAddresses).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/address/{id:[0-9]+}", addressHandler.GetByID).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/address/{id:[0-9]+}", addressHandler.Update).Methods(http.MethodPut)
	protectedRoutes.HandleFunc("/address/{id:[0-9]+}", addressHandler.Delete).Methods(http.MethodDelete)

	// Contact routes (padronizadas com regex e ordem consistente)
	protectedRoutes.HandleFunc("/contact", contactHandler.Create).Methods(http.MethodPost)
	protectedRoutes.HandleFunc("/contact/{id:[0-9]+}", contactHandler.GetByID).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/contact/user/{userID:[0-9]+}", contactHandler.GetByUser).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/contact/client/{clientID:[0-9]+}", contactHandler.GetByClient).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/contact/supplier/{supplierID:[0-9]+}", contactHandler.GetBySupplier).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/contact/{id:[0-9]+}", contactHandler.Update).Methods(http.MethodPut)
	protectedRoutes.HandleFunc("/contact/{id:[0-9]+}", contactHandler.Delete).Methods(http.MethodDelete)

	// Supplier routes (padronizadas com regex e ordem consistente)
	protectedRoutes.HandleFunc("/supplier", supplierHandler.CreateSupplier).Methods(http.MethodPost)
	protectedRoutes.HandleFunc("/supplier/{id:[0-9]+}", supplierHandler.GetSupplierByID).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/suppliers", supplierHandler.GetAllSuppliers).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/supplier/{id:[0-9]+}", supplierHandler.UpdateSupplier).Methods(http.MethodPut)
	protectedRoutes.HandleFunc("/supplier/{id:[0-9]+}", supplierHandler.DeleteSupplier).Methods(http.MethodDelete)

	return r
}

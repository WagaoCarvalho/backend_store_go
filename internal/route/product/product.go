package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/handler/product/filter"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/product/product"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwtMiddlewares "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repoFilter "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/filter"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/product"
	serviceFilter "github.com/WagaoCarvalho/backend_store_go/internal/service/product/filter"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product/product"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterProductRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwtMiddlewares.TokenBlacklist,
) {
	serverConfig := config.LoadServerConfig()
	baseURL := serverConfig.BaseURL
	idPath := serverConfig.IDPath

	// Repositórios
	newRepoProduct := repo.NewProduct(db)
	newRepoFilter := repoFilter.NewFilterProduct(db)

	// Serviços
	newServiceProduct := service.NewProductService(newRepoProduct)
	newServiceFilter := serviceFilter.NewProductFilterService(newRepoFilter)

	// Handlers
	newHandlerProduct := handler.NewProductHandler(newServiceProduct, log)
	newHandlerFilter := filter.NewProductFilterHandler(newServiceFilter, log)

	// Config JWT
	jwtCfg := config.LoadJwtConfig()
	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	// Rotas protegidas
	s := r.PathPrefix("/").Subrouter()
	s.Use(jwtMiddlewares.IsAuthByBearerToken(blacklist, log, jwtManager))

	// Constantes para caminhos
	const (
		product     = "/product"
		products    = "/products"
		filterPath  = "/filter"
		update      = "/update"
		delete      = "/delete"
		enable      = "/enable"
		disable     = "/disable"
		version     = "/version"
		stock       = "/stock"
		increase    = "/increase-stock"
		decrease    = "/decrease-stock"
		getStock    = "/get-stock"
		enableDisc  = "/enable-discount"
		disableDisc = "/disable-discount"
		applyDisc   = "/apply-discount"
	)

	// Rotas CRUD básicas
	s.HandleFunc(baseURL+product, newHandlerProduct.Create).Methods(http.MethodPost)
	s.HandleFunc(baseURL+product+idPath, newHandlerProduct.GetByID).Methods(http.MethodGet)
	s.HandleFunc(baseURL+product+idPath, newHandlerProduct.Update).Methods(http.MethodPut)
	s.HandleFunc(baseURL+product+idPath, newHandlerProduct.Delete).Methods(http.MethodDelete)

	// Rotas de status
	s.HandleFunc(baseURL+product+idPath+enable, newHandlerProduct.EnableProduct).Methods(http.MethodPatch)
	s.HandleFunc(baseURL+product+idPath+disable, newHandlerProduct.DisableProduct).Methods(http.MethodPatch)

	// Rotas de versão
	s.HandleFunc(baseURL+product+idPath+version, newHandlerProduct.GetVersionByID).Methods(http.MethodGet)

	// Rotas de estoque
	s.HandleFunc(baseURL+product+idPath+stock, newHandlerProduct.UpdateStock).Methods(http.MethodPatch)
	s.HandleFunc(baseURL+product+idPath+increase, newHandlerProduct.IncreaseStock).Methods(http.MethodPatch)
	s.HandleFunc(baseURL+product+idPath+decrease, newHandlerProduct.DecreaseStock).Methods(http.MethodPatch)
	s.HandleFunc(baseURL+product+idPath+getStock, newHandlerProduct.GetStock).Methods(http.MethodGet)

	// Rotas de desconto
	s.HandleFunc(baseURL+product+idPath+enableDisc, newHandlerProduct.EnableDiscount).Methods(http.MethodPatch)
	s.HandleFunc(baseURL+product+idPath+disableDisc, newHandlerProduct.DisableDiscount).Methods(http.MethodPatch)
	s.HandleFunc(baseURL+product+idPath+applyDisc, newHandlerProduct.ApplyDiscount).Methods(http.MethodPatch)

	// Rota de filtro
	s.HandleFunc(baseURL+products+filterPath, newHandlerFilter.Filter).Methods(http.MethodGet)
}

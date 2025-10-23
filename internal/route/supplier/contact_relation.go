package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/supplier/supplier_contact_relation"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_contact_relation"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier_contact_relation"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSupplierContactRelationRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	contactRepo := repo.NewSupplierContactRelation(db)
	contactService := service.NewSupplierContactRelation(contactRepo)
	contactHandler := handler.NewSupplierContactRelation(contactService, log)

	// Carregar config JWT
	jwtCfg := config.LoadJwtConfig()

	// Criar jwtManager que implementa JWTService
	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, jwtManager))

	// Rotas Supplier Contact Relations
	s.HandleFunc("/supplier-contact-relation", contactHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/supplier-contact-relations/{supplier_id:[0-9]+}", contactHandler.GetAllBySupplierID).Methods(http.MethodGet)
	s.HandleFunc("/supplier-contact-relation/{supplier_id:[0-9]+}/contact/{contact_id:[0-9]+}/exists", contactHandler.HasSupplierContactRelation).Methods(http.MethodGet)
	s.HandleFunc("/supplier-contact-relation/{supplier_id:[0-9]+}/contact/{contact_id:[0-9]+}", contactHandler.Delete).Methods(http.MethodDelete)
	s.HandleFunc("/supplier-contact-relation/{supplier_id:[0-9]+}", contactHandler.DeleteAll).Methods(http.MethodDelete)
}

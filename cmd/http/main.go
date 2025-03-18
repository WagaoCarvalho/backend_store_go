package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/cmd/http/routes"
	"github.com/WagaoCarvalho/backend_store_go/config"
	"github.com/WagaoCarvalho/backend_store_go/internal/repositories"
)

func main() {
	configs := config.LoadConfig()
	port := configs.Server.Port
	if port == "" {
		port = "5000"
	}

	repositories.TestConnection()
	fmt.Printf("API running on port %s\n", port)

	r := routes.NewRouter()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

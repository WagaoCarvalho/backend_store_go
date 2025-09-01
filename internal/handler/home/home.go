package handler

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func GetHome(w http.ResponseWriter, _ *http.Request) {
	utils.ToJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "Go RESTful Api backend_store",
	})
}

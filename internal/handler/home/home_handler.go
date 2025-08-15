package handler

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/pkg/utils"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	utils.ToJson(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "Go RESTful Api backend_store",
	})
}

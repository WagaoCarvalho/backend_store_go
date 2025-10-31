package validators

import (
	"strings"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

// --- Validação comum de paginação e ordenação ---
func ValidatePagination(limit, offset int) error {
	if limit <= 0 {
		return errMsg.ErrInvalidLimit
	}
	if offset < 0 {
		return errMsg.ErrInvalidOffset
	}
	return nil
}

func ValidateOrder(orderBy string, allowedFields map[string]bool, orderDir string) (string, error) {
	if !allowedFields[orderBy] {
		return "", errMsg.ErrInvalidOrderField
	}

	orderDir = strings.ToLower(orderDir)
	if orderDir != "asc" && orderDir != "desc" {
		return "", errMsg.ErrInvalidOrderDirection
	}

	return orderDir, nil
}

package repo

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

type supplierContactRelationRepo struct {
	db repo.DBExecutor
}

func NewSupplierContactRelation(db repo.DBExecutor) SupplierContactRelation {
	return &supplierContactRelationRepo{db: db}
}

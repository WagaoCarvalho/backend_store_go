package repo

import "github.com/WagaoCarvalho/backend_store_go/internal/iface"

type SaleRepo interface {
	iface.SaleReader
	iface.SaleWriter
	iface.SaleStatus
}

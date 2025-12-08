package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/sale"

type SaleRepo interface {
	iface.SaleReader
	iface.SaleWriter
	iface.SaleStatus
	iface.SaleVersion
}

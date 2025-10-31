package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/sale"

type SaleItemRepo interface {
	iface.SaleItemReader
	iface.SaleItemWriter
	iface.SaleItemChecker
}

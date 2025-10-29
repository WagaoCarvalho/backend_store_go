package services

import "github.com/WagaoCarvalho/backend_store_go/internal/iface"

type SaleService interface {
	iface.SaleReader
	iface.SaleWriter
	iface.SaleStatus
}

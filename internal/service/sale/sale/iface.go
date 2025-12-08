package services

import sale_iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/sale"

type SaleService interface {
	sale_iface.SaleReader
	sale_iface.SaleWriter
	sale_iface.SaleStatus
	sale_iface.SaleVersion
}

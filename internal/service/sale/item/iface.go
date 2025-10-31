package services

import item "github.com/WagaoCarvalho/backend_store_go/internal/iface/sale"

type SaleItemService interface {
	item.SaleItemReader
	item.SaleItemWriter
	item.SaleItemChecker
}

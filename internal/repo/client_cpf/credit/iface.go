package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/client_cpf"

type ClientCredit interface {
	iface.ClientCreditReader
	iface.ClientCreditWriter
	iface.ClientCreditStatus
}

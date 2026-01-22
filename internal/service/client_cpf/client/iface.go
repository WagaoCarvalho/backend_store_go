package services

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/client_cpf"

type Client interface {
	iface.ClientCpfReader
	iface.ClientCpfWriter
	iface.ClientCpfStatus
	iface.ClientCpfVersion
}

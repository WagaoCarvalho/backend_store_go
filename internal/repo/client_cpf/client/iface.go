package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/client_cpf"

type ClientCpf interface {
	iface.ClientCpfReader
	iface.ClientCpfWriter
	iface.ClientCpfStatus
	iface.ClientCpfVersion
}

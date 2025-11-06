package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/contact"

type Contact interface {
	iface.ContactReader
	iface.ContactWriter
}

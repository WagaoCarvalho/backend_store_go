package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/user"

type UserContactRelation interface {
	iface.UserContactRelationReader
	iface.UserContactRelationWriter
}

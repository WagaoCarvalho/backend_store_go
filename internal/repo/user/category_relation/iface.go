package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/user"

type UserCategoryRelation interface {
	iface.UserCategoryRelationReader
	iface.UserCategoryRelationWriter
}

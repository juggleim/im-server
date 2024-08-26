package mongodbs

import "im-server/commons/mongocommons"

func RegistCollections() {
	mongocommons.Register(&UserTagsDao{})
}

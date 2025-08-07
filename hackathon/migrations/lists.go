package migrations

import (
	"github.com/minhho2511/elotusteam-test/pkgs/db"
)

func MigrationLists() []db.MFile {
	return []db.MFile{
		CreateUserTable{},
		CreateFileTable{},
	}
}

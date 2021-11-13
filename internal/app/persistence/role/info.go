package role

import (
	"pbRole"
	store "src/_infrastructure/store/mysql"
)

type setName struct {
	pbRole.EventSetName
}

func (s *setName) Write(db *store.MySQLClient) error {
	panic("implement me")
}



